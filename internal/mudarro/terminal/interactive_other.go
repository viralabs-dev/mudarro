//go:build !linux && !darwin

package terminal

import "os"

func (v *Menu) ChooseInteractive(in *os.File, title string, labels, previews []string) (int, bool, error) {
	return 0, false, nil
}
