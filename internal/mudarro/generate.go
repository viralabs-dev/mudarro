package mudarro

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type GeneratedFile struct {
	Data []byte
	Mode os.FileMode
}
type fileManifest struct {
	Version int               `json:"version"`
	Files   map[string]string `json:"files"`
}

func digest(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func Generate(root string, c Config, dry bool, w io.Writer) error {
	if e := c.Validate(root); e != nil {
		return e
	}
	files := map[string]GeneratedFile{}
	put := func(path, data string, mode os.FileMode) error {
		path = filepath.Clean(path)
		if _, ok := files[path]; ok {
			return fmt.Errorf("dois geradores reivindicam %s; configure um único proprietário", path)
		}
		files[path] = GeneratedFile{[]byte(data), mode}
		return nil
	}
	launcher := "#!/usr/bin/env bash\nset -euo pipefail\nROOT=\"$(cd -- \"$(dirname -- \"${BASH_SOURCE[0]}\")\" && pwd)\"\nexec mudarro menu --root \"$ROOT\" \"$@\"\n"
	if e := put("menu.sh", launcher, 0755); e != nil {
		return e
	}
	for _, s := range c.Services {
		for _, a := range Actions(s) {
			p := filepath.Join(".mudarro", "scripts", s.ID, a.Name+".sh")
			body := "#!/usr/bin/env bash\nset -euo pipefail\nROOT=\"$(cd -- \"$(dirname -- \"${BASH_SOURCE[0]}\")/../../..\" && pwd)\"\nexec mudarro run " + shellQuote(s.ID+":"+a.Name) + " --root \"$ROOT\" \"$@\"\n"
			if e := put(p, body, 0755); e != nil {
				return e
			}
		}
		scaffold, e := scaffoldFiles(s)
		if e != nil {
			return fmt.Errorf("%s: %w", s.ID, e)
		}
		for p, f := range scaffold {
			if e := put(filepath.Join(s.Dir, p), string(f.Data), f.Mode); e != nil {
				return e
			}
		}
	}
	manifestPath, e := safePath(root, ".mudarro/generated.json")
	if e != nil {
		return e
	}
	old := fileManifest{Version: 1, Files: map[string]string{}}
	if b, e := os.ReadFile(manifestPath); e == nil {
		if e = json.Unmarshal(b, &old); e != nil {
			return fmt.Errorf("manifest inválido: %w", e)
		}
	} else if !os.IsNotExist(e) {
		return e
	}
	next := fileManifest{Version: 1, Files: map[string]string{}}
	for p, h := range old.Files {
		next.Files[p] = h
	}
	paths := []string{}
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	writes := []string{}
	// Preflight every destination before writing anything. Existing user files are never adopted.
	for _, p := range paths {
		full, e := safePath(root, p)
		if e != nil {
			return e
		}
		f := files[p]
		b, e := os.ReadFile(full)
		if e == nil {
			if h, owned := old.Files[p]; owned {
				if digest(b) != h {
					return fmt.Errorf("arquivo gerado alterado manualmente: %s; preserve sua versão ou restaure antes de gerar", p)
				}
				next.Files[p] = digest(f.Data)
				if digest(b) == digest(f.Data) {
					continue
				}
			} else {
				fmt.Fprintln(w, "preservado:", p)
				continue
			}
		} else if !os.IsNotExist(e) {
			return e
		} else {
			next.Files[p] = digest(f.Data)
		}
		writes = append(writes, p)
	}
	for _, p := range writes {
		fmt.Fprintln(w, "gerar:", p)
		if dry {
			continue
		}
		full, _ := safePath(root, p)
		if e := atomicWrite(full, files[p].Data, files[p].Mode); e != nil {
			return e
		}
	}
	if dry {
		return nil
	}
	b, _ := json.MarshalIndent(next, "", "  ")
	return atomicWrite(manifestPath, append(b, '\n'), 0644)
}
func atomicWrite(path string, b []byte, mode os.FileMode) error {
	if e := os.MkdirAll(filepath.Dir(path), 0755); e != nil {
		return e
	}
	f, e := os.CreateTemp(filepath.Dir(path), ".mudarro-tmp-*")
	if e != nil {
		return e
	}
	defer os.Remove(f.Name())
	if e = f.Chmod(mode); e != nil {
		f.Close()
		return e
	}
	if _, e = f.Write(b); e != nil {
		f.Close()
		return e
	}
	if e = f.Close(); e != nil {
		return e
	}
	return os.Rename(f.Name(), path)
}
