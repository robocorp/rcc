// +build darwin linux !windows

package pathlib

import (
	"os"
	"syscall"

	"github.com/robocorp/rcc/common"
)

func Locker(filename string, trycount int) (Releaser, error) {
	if Lockless {
		return Fake(), nil
	}
	if common.TraceFlag {
		defer common.Stopwatch("LOCKER: Got lock on %v in", filename).Report()
	}
	common.Trace("LOCKER: Want lock on: %v", filename)
	_, err := EnsureParentDirectory(filename)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, err
	}
	err = syscall.Flock(int(file.Fd()), int(syscall.LOCK_EX))
	if err != nil {
		return nil, err
	}
	return &Locked{file}, nil
}

func (it Locked) Release() error {
	defer it.Close()
	err := syscall.Flock(int(it.Fd()), int(syscall.LOCK_UN))
	common.Trace("LOCKER: release with err: %v", err)
	return err
}
