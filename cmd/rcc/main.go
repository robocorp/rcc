package main

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/cmd"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pathlib"
)

const (
	timezonekey = `rcc.cli.tz`
	oskey       = `rcc.cli.os`
	daily       = 60 * 60 * 24
)

var (
	markedAlready = false
)

func TimezoneMetric() error {
	cache, err := operations.SummonCache()
	if err != nil {
		return err
	}
	deadline, ok := cache.Stamps[timezonekey]
	if ok && deadline > common.When {
		return nil
	}
	cache.Stamps[timezonekey] = common.When + daily
	zone := time.Now().Format("MST-0700")
	cloud.BackgroundMetric(common.ControllerIdentity(), timezonekey, zone)
	cloud.BackgroundMetric(common.ControllerIdentity(), oskey, common.Platform())
	return cache.Save()
}

func ExitProtection() {
	status := recover()
	if status != nil {
		markTempForRecycling()
		exit, ok := status.(common.ExitCode)
		if ok {
			exit.ShowMessage()
			cloud.WaitTelemetry()
			common.WaitLogs()
			os.Exit(exit.Code)
		}
		cloud.BackgroundMetric(common.ControllerIdentity(), "rcc.panic.origin", cmd.Origin())
		cloud.WaitTelemetry()
		common.WaitLogs()
		panic(status)
	}
	cloud.WaitTelemetry()
	common.WaitLogs()
}

func startTempRecycling() {
	defer common.Timeline("temp recycling done")
	pattern := filepath.Join(common.RobocorpTempRoot(), "*", "recycle.now")
	found, err := filepath.Glob(pattern)
	if err != nil {
		common.Debug("Recycling failed, reason: %v", err)
		return
	}
	for _, filename := range found {
		folder := filepath.Dir(filename)
		changed, err := pathlib.Modtime(folder)
		if err == nil && time.Since(changed) > 48*time.Hour {
			go os.RemoveAll(folder)
		}
	}
	runtime.Gosched()
}

func markTempForRecycling() {
	if markedAlready {
		return
	}
	target := common.RobocorpTempName()
	if pathlib.Exists(target) {
		filename := filepath.Join(target, "recycle.now")
		os.WriteFile(filename, []byte("True"), 0o644)
		common.Debug("Marked %q for recycling.", target)
		markedAlready = true
	}
}

func main() {
	defer ExitProtection()

	if common.SharedHolotree {
		common.TimelineBegin("Start [shared mode].")
	} else {
		common.TimelineBegin("Start [private mode].")
	}
	defer common.EndOfTimeline()
	go startTempRecycling()
	defer markTempForRecycling()
	defer os.Stderr.Sync()
	defer os.Stdout.Sync()
	cmd.Execute()
	common.Timeline("Command execution done.")
	TimezoneMetric()
}
