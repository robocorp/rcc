// +build windows !darwin !linux

package pretty

import (
	"syscall"

	"github.com/robocorp/rcc/common"
)

const (
	ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x4
)

func localSetup() {
	Disabled = true
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	if kernel32 == nil {
		common.Trace("Cannot use colors. Did not get kernel32.dll!")
		return
	}
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	if setConsoleMode == nil {
		common.Trace("Cannot use colors. Did not get SetConsoleMode!")
		return
	}
	_, _, err := setConsoleMode.Call(ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	Disabled = err != nil
	if err != nil {
		common.Trace("Cannot use colors. Got error '%v'!", err)
	}
}
