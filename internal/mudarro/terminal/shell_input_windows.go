//go:build windows

package terminal

import (
	"context"
	"os"
	"sync"
)

// ExecutionInput owns its pipe and input pump; stop joins before restoring the
// console mode. The child receives a pipe, never the console handle.
func (v *Menu) ExecutionInput(ctx context.Context, cancel context.CancelFunc) (*os.File, func(), error) {
	restore, err := interactiveRaw(v.frame.in)
	if err != nil {
		return nil, nil, err
	}
	r, w, err := os.Pipe()
	if err != nil {
		restore()
		return nil, nil, err
	}
	pumpCtx, stop := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer w.Close()
		for {
			if pumpCtx.Err() != nil {
				return
			}
			chunk, e := interactiveNext(v.frame.in)
			if e != nil {
				cancel()
				return
			}
			for _, b := range chunk {
				switch b {
				case 3:
					cancel()
					return
				case 4, 26: // Ctrl+D, Ctrl+Z: end of input
					return
				case 13:
					b = 10
				}
				if pumpCtx.Err() != nil {
					return
				}
				// Anonymous pipes cannot be non-blocking on Windows; closing
				// the read end in cleanup unblocks a full pipe.
				if _, e = w.Write([]byte{b}); e != nil {
					return
				}
			}
		}
	}()
	var once sync.Once
	cleanup := func() { once.Do(func() { stop(); r.Close(); <-done; restore() }) }
	return r, cleanup, nil
}
