// Package projectfs guards paths and bounded manifest reads within a project.
package projectfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func SafePath(root, rel string) (string, error) {
	if filepath.IsAbs(rel) || rel == "" {
		return "", fmt.Errorf("caminho relativo inválido: %q", rel)
	}
	clean := filepath.Clean(rel)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("caminho fora do projeto: %q", rel)
	}
	current := root
	for _, part := range strings.Split(clean, string(filepath.Separator)) {
		current = filepath.Join(current, part)
		st, err := os.Lstat(current)
		if err == nil && st.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("symlink não permitido: %s", current)
		}
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}
	return filepath.Join(root, clean), nil
}

func ReadManifest(root, dir, file string) ([]byte, error) {
	p, e := SafePath(root, filepath.Join(dir, file))
	if e != nil {
		return nil, e
	}
	st, e := os.Stat(p)
	if e != nil {
		return nil, e
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("manifest deve ser arquivo regular: %s", p)
	}
	const limit = 4 << 20
	if st.Size() > limit {
		return nil, fmt.Errorf("manifest grande demais: %s", p)
	}
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	st, e = f.Stat()
	if e != nil {
		return nil, e
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("manifest deve ser arquivo regular: %s", p)
	}
	b, e := io.ReadAll(io.LimitReader(f, limit+1))
	if e != nil {
		return nil, e
	}
	if len(b) > limit {
		return nil, fmt.Errorf("manifest grande demais: %s", p)
	}
	return b, nil
}
