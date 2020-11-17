package main

import (
	"os"

	"github.com/robocorp/rcc/cmd"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
)

func ExitProtection() {
	status := recover()
	if status != nil {
		exit, ok := status.(common.ExitCode)
		if ok {
			exit.ShowMessage()
			os.Exit(exit.Code)
		}
		operations.SendMetric("rcc", "rcc.panic.origin", cmd.Origin())
		panic(status)
	}
}

func main() {
	defer os.Stderr.Sync()
	defer os.Stdout.Sync()
	defer ExitProtection()
	cmd.Execute()
}
