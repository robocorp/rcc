package pathlib

import (
	"os"

	"github.com/robocorp/rcc/common"
)

type Releaser interface {
	Release() error
}

type Locked struct {
	*os.File
}

type fake bool

func (it fake) Release() error {
	return common.Trace("LOCKER: lockless mode release.")
}

func Fake() Releaser {
	common.Trace("LOCKER: lockless mode.")
	return fake(true)
}
