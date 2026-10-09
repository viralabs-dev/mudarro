//go:build windows

package winapi

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procThread32First            = kernel32.NewProc("Thread32First")
	procThread32Next             = kernel32.NewProc("Thread32Next")
	procOpenThread               = kernel32.NewProc("OpenThread")
	procResumeThread             = kernel32.NewProc("ResumeThread")
)

const (
	// CreateSuspended lets a caller assign the process to a Job Object
	// before it runs, so no descendant can escape the job.
	CreateSuspended     = 0x00000004
	th32csSnapThread    = 0x00000004
	threadSuspendResume = 0x0002
)

type threadEntry32 struct {
	Size           uint32
	Usage          uint32
	ThreadID       uint32
	OwnerProcessID uint32
	BasePri        int32
	DeltaPri       int32
	Flags          uint32
}

// ResumeProcess resumes every thread of a process created with
// CreateSuspended. exec.Cmd does not expose the primary thread handle, so the
// threads are found through a Toolhelp snapshot.
func ResumeProcess(pid int) error {
	r, _, err := procCreateToolhelp32Snapshot.Call(th32csSnapThread, 0)
	if syscall.Handle(r) == syscall.InvalidHandle {
		return callErr(0, err)
	}
	snapshot := syscall.Handle(r)
	defer syscall.CloseHandle(snapshot)
	entry := threadEntry32{Size: uint32(unsafe.Sizeof(threadEntry32{}))}
	resumed := 0
	r, _, err = procThread32First.Call(uintptr(snapshot), uintptr(unsafe.Pointer(&entry)))
	for r != 0 {
		if entry.OwnerProcessID == uint32(pid) {
			h, _, openErr := procOpenThread.Call(threadSuspendResume, 0, uintptr(entry.ThreadID))
			if h == 0 {
				return callErr(0, openErr)
			}
			count, _, resumeErr := procResumeThread.Call(h)
			syscall.CloseHandle(syscall.Handle(h))
			if count == ^uintptr(0) || uint32(count) == ^uint32(0) {
				return callErr(0, resumeErr)
			}
			resumed++
		}
		entry.Size = uint32(unsafe.Sizeof(threadEntry32{}))
		r, _, err = procThread32Next.Call(uintptr(snapshot), uintptr(unsafe.Pointer(&entry)))
	}
	if resumed == 0 {
		return fmt.Errorf("no thread to resume for process %d: %v", pid, err)
	}
	return nil
}
