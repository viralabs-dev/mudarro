//go:build linux || darwin

package terminal

import (
	"context"
	"os"
	"sync"
	"syscall"
	"time"
)

// ExecutionInput owns its pipe and input pump; stop joins before restoring termios.
// The child receives a pipe, never the user's terminal or an uncancellable Reader.
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
	writeFD := int(w.Fd())
	if err = syscall.SetNonblock(writeFD, true); err != nil {
		r.Close()
		w.Close()
		restore()
		return nil, nil, err
	}
	pumpCtx, stop := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer w.Close()
		buf := make([]byte, 1)
		for {
			if pumpCtx.Err() != nil {
				return
			}
			ready, e := interactiveReady(int(v.frame.in.Fd()))
			if e == syscall.EINTR {
				continue
			}
			if e != nil {
				cancel()
				return
			}
			if !ready {
				continue
			}
			n, e := syscall.Read(int(v.frame.in.Fd()), buf)
			if e == syscall.EINTR || e == syscall.EAGAIN {
				continue
			}
			if e != nil || n == 0 {
				cancel()
				return
			}
			if buf[0] == 3 {
				cancel()
				return
			}
			if buf[0] == 4 {
				return
			}
			if buf[0] == 13 {
				buf[0] = 10
			}
			for {
				if pumpCtx.Err() != nil {
					return
				}
				_, e = syscall.Write(writeFD, buf[:1])
				if e == nil {
					break
				}
				if e != syscall.EAGAIN && e != syscall.EINTR {
					return
				}
				select {
				case <-pumpCtx.Done():
					return
				case <-time.After(10 * time.Millisecond):
				}
			}
		}
	}()
	var once sync.Once
	cleanup := func() { once.Do(func() { stop(); <-done; r.Close(); restore() }) }
	return r, cleanup, nil
}
