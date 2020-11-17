// +build windows

package pathlib

import (
	"os"
	"syscall"
	"time"

	"github.com/robocorp/rcc/common"
)

const (
	LOCKFILE_EXCLUSIVE_LOCK = 2
)

// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-lockfile
// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-unlockfile

var (
	kernel32, _   = syscall.LoadLibrary("kernel32.dll")
	lockFile, _   = syscall.GetProcAddress(kernel32, "LockFile")
	unlockFile, _ = syscall.GetProcAddress(kernel32, "UnlockFile")
)

type filehandle interface {
	Fd() uintptr
}

func Locker(filename string, trycount int) (Releaser, error) {
	if Lockless {
		return Fake(), nil
	}
	var file *os.File
	var err error
	if common.TraceFlag {
		defer func() {
			common.Stopwatch("LOCKER: Leaving lock on %v with %v retries left in", filename, trycount).Report()
		}()
	}
	common.Trace("LOCKER: Want lock on: %v", filename)
	_, err = EnsureParentDirectory(filename)
	if err != nil {
		return nil, err
	}
	for {
		trycount -= 1
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil && trycount < 0 {
			return nil, err
		}
		if err != nil {
			time.Sleep(40 * time.Millisecond)
			continue
		}
		break
	}
	for {
		trycount -= 1
		success, err := trylock(lockFile, file)
		if err != nil && trycount < 0 {
			return nil, err
		}
		if success {
			return &Locked{file}, nil
		}
		time.Sleep(40 * time.Millisecond)
	}
}

func (it Locked) Release() error {
	success, err := trylock(unlockFile, it)
	common.Trace("LOCKER: release success: %v with err: %v", success, err)
	return err
}

func trylock(tool uintptr, identity filehandle) (bool, error) {
	handle := syscall.Handle(identity.Fd())
	primary, _, err := syscall.Syscall6(
		tool,
		5,
		uintptr(handle),
		uintptr(0),
		uintptr(0),
		uintptr(1),
		uintptr(0),
		uintptr(0))
	if primary == 0 {
		return false, err
	}
	return true, nil
}
