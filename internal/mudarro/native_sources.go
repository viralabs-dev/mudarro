package mudarro

import (
	"fmt"
	"github.com/viralabs-dev/mudarro/internal/mudarro/adapters"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// scanNativeSources inventories native source files without interpreting source,
// invoking tools, or deciding compiler flags, application commands or startup.
func scanNativeSources(root string, paths []string, dirs map[string]map[string]bool, r *Report) {
	makes := map[string]bool{}
	for _, dir := range paths {
		if !dirs[dir]["Makefile"] {
			continue
		}
		if _, err := read(root, dir, "Makefile"); err == nil {
			makes[dir] = true
		}
	}
	groups := map[string]bool{}
	for _, dir := range paths {
		files := make([]string, 0, len(dirs[dir]))
		for file := range dirs[dir] {
			files = append(files, file)
		}
		sort.Strings(files)
		for _, file := range files {
			ext := filepath.Ext(file)
			kind := ""
			source := false
			switch ext {
			case ".c":
				kind = "c-source"
				source = true
			case ".cpp", ".cc", ".cxx":
				kind = "cpp-source"
				source = true
			case ".h":
				kind = "c-header"
			case ".hpp":
				kind = "cpp-header"
			default:
				continue
			}
			rel := filepath.Join(dir, file)
			st, err := os.Lstat(filepath.Join(root, rel))
			if err != nil || !st.Mode().IsRegular() {
				r.Warnings = append(r.Warnings, fmt.Sprintf("%s: native source/header must be a regular file", filepath.ToSlash(rel)))
				continue
			}
			r.Evidence = append(r.Evidence, Evidence{Path: filepath.ToSlash(rel), Kind: kind})
			if !source {
				continue
			}
			makeRoot := ""
			makeDepth := -1
			for candidate := range makes {
				if nativeContains(candidate, dir) && len(candidate) > makeDepth {
					makeRoot = candidate
					makeDepth = len(candidate)
				}
			}
			existing := -1
			existingDepth := -1
			for i, s := range r.Config.Services {
				if nativeContains(s.Dir, dir) && len(s.Dir) > existingDepth {
					existing = i
					existingDepth = len(s.Dir)
				}
			}
			if existing >= 0 && (makeRoot == "" || makeDepth <= existingDepth) {
				nativePending(&r.Config.Services[existing], makeRoot != "")
				continue
			}
			anchor := dir
			if makeRoot != "" {
				anchor = makeRoot
			}
			groups[anchor] = makes[anchor]
		}
	}
	anchors := make([]string, 0, len(groups))
	for anchor := range groups {
		anchors = append(anchors, anchor)
	}
	sort.Strings(anchors)
	for _, anchor := range anchors {
		existing := -1
		for i, s := range r.Config.Services {
			if s.Dir == anchor {
				existing = i
				break
			}
		}
		if existing >= 0 {
			nativePending(&r.Config.Services[existing], groups[anchor])
			continue
		}
		s := adapters.BaseService(anchor, "custom")
		s.Infrastructure.Kind = "local"
		nativePending(s, groups[anchor])
		collision := false
		for _, other := range r.Config.Services {
			if s.ID == other.ID {
				collision = true
				break
			}
		}
		if collision {
			r.Warnings = append(r.Warnings, fmt.Sprintf("%s: native service ID %s is ambiguous; configure services explicitly", filepath.ToSlash(anchor), s.ID))
			continue
		}
		r.Config.Services = append(r.Config.Services, *s)
	}
}
func nativeContains(parent, dir string) bool {
	rel, err := filepath.Rel(parent, dir)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
func nativePending(s *Service, hasMake bool) {
	message := "Native sources detected: declare commands explicitly; compiler flags and startup are not inferred"
	if hasMake {
		message = "Native sources with Makefile: select explicit Make targets or declare commands; compiler flags and startup are not inferred"
	}
	for _, p := range s.Pending {
		if p == message {
			return
		}
	}
	s.Pending = append(s.Pending, message)
}
