package pretty

import (
	"os"

	"github.com/mattn/go-isatty"
	"github.com/robocorp/rcc/common"
)

var (
	Disabled    bool
	Interactive bool
	Red         string
	Green       string
	Cyan        string
	Reset       string
)

func Setup() {
	stdin := isatty.IsTerminal(os.Stdin.Fd())
	stdout := isatty.IsTerminal(os.Stdout.Fd())
	stderr := isatty.IsTerminal(os.Stderr.Fd())
	Interactive = stdin && stdout && stderr

	localSetup()

	common.Trace("Interactive mode enabled: %v; colors enabled: %v", Interactive, !Disabled)
	if Interactive && !Disabled {
		Red = csi("1;31m")
		Green = csi("1;32m")
		Cyan = csi("1;36m")
		Reset = csi("0m")
	}
}
