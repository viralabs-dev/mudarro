//go:build windows

package winapi

import (
	"errors"
	"os"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// kernel32 is a KnownDLL, so loading it by name cannot be redirected by the
// current directory or PATH.
var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procCreateJobObjectW           = kernel32.NewProc("CreateJobObjectW")
	procSetInformationJobObject    = kernel32.NewProc("SetInformationJobObject")
	procAssignProcessToJobObject   = kernel32.NewProc("AssignProcessToJobObject")
	procTerminateJobObject         = kernel32.NewProc("TerminateJobObject")
	procGenerateConsoleCtrlEvent   = kernel32.NewProc("GenerateConsoleCtrlEvent")
	procLockFileEx                 = kernel32.NewProc("LockFileEx")
	procUnlockFileEx               = kernel32.NewProc("UnlockFileEx")
	procSetConsoleMode             = kernel32.NewProc("SetConsoleMode")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procPeekConsoleInputW          = kernel32.NewProc("PeekConsoleInputW")
	procReadConsoleInputW          = kernel32.NewProc("ReadConsoleInputW")
	procCreateMutexW               = kernel32.NewProc("CreateMutexW")
	procOpenMutexW                 = kernel32.NewProc("OpenMutexW")
	procQueryFullProcessImageNameW = kernel32.NewProc("QueryFullProcessImageNameW")
)

const (
	// Process access rights.
	processTerminate               = 0x0001
	processSetQuota                = 0x0100
	processQueryLimitedInformation = 0x1000
	synchronize                    = 0x00100000
	stillActive                    = 259

	// Job Objects.
	jobObjectExtendedLimitInformation = 9
	jobObjectLimitKillOnJobClose      = 0x2000

	// Console.
	ctrlBreakEvent        = 1
	EnableProcessedInput  = 0x0001
	EnableLineInput       = 0x0002
	EnableEchoInput       = 0x0004
	EnableVirtualTerminal = 0x0200 // ENABLE_VIRTUAL_TERMINAL_INPUT
	EnableVTProcessing    = 0x0004 // ENABLE_VIRTUAL_TERMINAL_PROCESSING (output)
	keyEvent              = 0x0001
	lockfileExclusiveLock = 0x0002
	errorAlreadyExists    = syscall.Errno(183)
	waitTimeout           = 0x00000102
	lockAllBytes          = ^uint32(0)
)

// CREATE_NEW_PROCESS_GROUP and DETACHED_PROCESS for SysProcAttr.CreationFlags.
const (
	CreateNewProcessGroup = syscall.CREATE_NEW_PROCESS_GROUP
	DetachedProcess       = 0x00000008
)

type ioCounters struct {
	ReadOperationCount, WriteOperationCount, OtherOperationCount uint64
	ReadTransferCount, WriteTransferCount, OtherTransferCount    uint64
}
type jobBasicLimit struct {
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit     int64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
}
type jobExtendedLimit struct {
	BasicLimitInformation jobBasicLimit
	IoInfo                ioCounters
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

func callErr(r uintptr, err error) error {
	if r != 0 {
		return nil
	}
	if errno, ok := err.(syscall.Errno); ok && errno != 0 {
		return errno
	}
	return syscall.EINVAL
}

// Job owns a Windows Job Object; terminating it ends every process assigned
// to it, including descendants created after the assignment.
type Job struct{ handle syscall.Handle }

// NewJob creates an anonymous Job Object. killOnClose ends the whole tree when
// the last handle closes, including when the owning process dies.
func NewJob(killOnClose bool) (*Job, error) {
	r, _, err := procCreateJobObjectW.Call(0, 0)
	if r == 0 {
		return nil, callErr(r, err)
	}
	j := &Job{handle: syscall.Handle(r)}
	if killOnClose {
		var info jobExtendedLimit
		info.BasicLimitInformation.LimitFlags = jobObjectLimitKillOnJobClose
		r, _, err = procSetInformationJobObject.Call(uintptr(j.handle), jobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
		if e := callErr(r, err); e != nil {
			j.Close()
			return nil, e
		}
	}
	return j, nil
}

// Assign places a running process (by PID) in the job.
func (j *Job) Assign(pid int) error {
	h, err := syscall.OpenProcess(processSetQuota|processTerminate, false, uint32(pid))
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(h)
	r, _, err := procAssignProcessToJobObject.Call(uintptr(j.handle), uintptr(h))
	return callErr(r, err)
}

// Terminate ends every process in the job with exit code 1.
func (j *Job) Terminate() error {
	r, _, err := procTerminateJobObject.Call(uintptr(j.handle), 1)
	return callErr(r, err)
}

// Close releases the handle; with killOnClose it also ends the tree.
func (j *Job) Close() error {
	if j == nil || j.handle == 0 {
		return nil
	}
	err := syscall.CloseHandle(j.handle)
	j.handle = 0
	return err
}

// CtrlBreak sends CTRL_BREAK_EVENT to a process group created with
// CREATE_NEW_PROCESS_GROUP that shares this process's console.
func CtrlBreak(pid int) error {
	r, _, err := procGenerateConsoleCtrlEvent.Call(ctrlBreakEvent, uintptr(pid))
	return callErr(r, err)
}

// LockFile blocks until an exclusive lock on the whole file is held.
func LockFile(f *os.File) error {
	var ol syscall.Overlapped
	r, _, err := procLockFileEx.Call(f.Fd(), lockfileExclusiveLock, 0, uintptr(lockAllBytes), uintptr(lockAllBytes), uintptr(unsafe.Pointer(&ol)))
	return callErr(r, err)
}

// UnlockFile releases LockFile.
func UnlockFile(f *os.File) error {
	var ol syscall.Overlapped
	r, _, err := procUnlockFileEx.Call(f.Fd(), 0, uintptr(lockAllBytes), uintptr(lockAllBytes), uintptr(unsafe.Pointer(&ol)))
	return callErr(r, err)
}

// ProcessAlive reports whether pid is a running process and returns its image path.
func ProcessAlive(pid int) (string, bool) {
	h, err := syscall.OpenProcess(processQueryLimitedInformation|synchronize, false, uint32(pid))
	if err != nil {
		return "", false
	}
	defer syscall.CloseHandle(h)
	var code uint32
	if syscall.GetExitCodeProcess(h, &code) != nil || code != stillActive {
		return "", false
	}
	buf := make([]uint16, 32768)
	size := uint32(len(buf))
	r, _, _ := procQueryFullProcessImageNameW.Call(uintptr(h), 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if r == 0 {
		return "", false
	}
	return string(utf16.Decode(buf[:size])), true
}

// SameExecutable compares two image paths the way NTFS does by default.
func SameExecutable(a, b string) bool {
	if strings.EqualFold(a, b) {
		return true
	}
	sa, ea := os.Stat(a)
	sb, eb := os.Stat(b)
	return ea == nil && eb == nil && os.SameFile(sa, sb)
}

// HoldMutex creates a named mutex that exists while this process lives. It
// refuses names that already exist, so a forged process cannot claim one.
func HoldMutex(name string) (func(), error) {
	p, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	r, _, callErrno := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(p)))
	if r == 0 {
		return nil, callErr(r, callErrno)
	}
	if errors.Is(callErrno, errorAlreadyExists) {
		syscall.CloseHandle(syscall.Handle(r))
		return nil, errorAlreadyExists
	}
	return func() { syscall.CloseHandle(syscall.Handle(r)) }, nil
}

// MutexExists reports whether a named mutex is currently held open by some process.
func MutexExists(name string) bool {
	p, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return false
	}
	r, _, _ := procOpenMutexW.Call(synchronize, 0, uintptr(unsafe.Pointer(p)))
	if r == 0 {
		return false
	}
	syscall.CloseHandle(syscall.Handle(r))
	return true
}
