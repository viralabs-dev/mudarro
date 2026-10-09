//go:build windows

package mudarro

import (
	"fmt"
	"os"
)

// openPreviewSource refuses symlinks and other reparse points before opening.
// Windows has no O_NOFOLLOW; callers re-check the opened handle is a regular file.
func openPreviewSource(path string) (*os.File, error) {
	st, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("preview source is not a regular file: %s", path)
	}
	return os.Open(path)
}
