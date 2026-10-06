package mudarro

import (
	"fmt"
	"os"
	"path/filepath"
)

func (r Runner) sqliteInit(root string, s Service) error {
	path := s.Database.Path
	if path == "" {
		path = "app.db"
	}
	full, err := safePath(root, filepath.Join(s.Dir, path))
	if err != nil {
		return err
	}
	if r.Dry {
		fmt.Fprintln(r.Out, "criar SQLite se ausente:", full)
		return nil
	}
	if err = os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(full, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if os.IsExist(err) {
		fmt.Fprintln(r.Out, "SQLite existente preservado")
		return nil
	}
	if err != nil {
		return err
	}
	return f.Close()
}
