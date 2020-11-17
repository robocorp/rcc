package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/robot"
	"github.com/robocorp/rcc/xviper"

	"github.com/spf13/cobra"
)

var testrunCmd = &cobra.Command{
	Use:     "testrun",
	Aliases: []string{"test", "t"},
	Short:   "Run a task in a clean environment and clean directory.",
	Long:    "Run a task in a clean environment and clean directory.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Task testrun lasted").Report()
		}
		ok := conda.MustConda()
		if !ok {
			pretty.Exit(4, "Could not get miniconda installed.")
		}
		defer xviper.RunMinutes().Done()
		now := time.Now()
		marker := now.Unix()
		zipfile := filepath.Join(os.TempDir(), fmt.Sprintf("testrun%x.zip", marker))
		defer os.Remove(zipfile)
		common.Debug("Using temporary zip file: %v", zipfile)
		sourceDir := filepath.Dir(robotFile)
		testrunDir := filepath.Join(sourceDir, "testrun", now.Format("2006-01-02_15_04_05"))
		err := os.MkdirAll(testrunDir, 0o755)
		if err != nil {
			pretty.Exit(1, "Error: %v", err)
		}
		err = operations.Zip(sourceDir, zipfile, ignores)
		if err != nil {
			pretty.Exit(2, "Error: %v", err)
		}
		sentinelTime := time.Now()
		workarea := filepath.Join(os.TempDir(), fmt.Sprintf("workarea%x", marker))
		defer os.RemoveAll(workarea)
		common.Debug("Using temporary workarea: %v", workarea)
		err = operations.Unzip(workarea, zipfile, false, true)
		if err != nil {
			pretty.Exit(3, "Error: %v", err)
		}
		defer pathlib.Walk(workarea, pathlib.IgnoreOlder(sentinelTime).Ignore, TargetDir(testrunDir).CopyBack)
		targetRobot := robot.DetectConfigurationName(workarea)
		simple, config, todo, label := operations.LoadTaskWithEnvironment(targetRobot, runTask, forceFlag)
		defer common.Log("Moving outputs to %v directory.", testrunDir)
		operations.BackgroundMetric("rcc", "rcc.cli.testrun", common.Version)
		operations.SelectExecutionModel(captureRunFlags(), simple, todo.Commandline(), config, todo, label, false, nil)
	},
}

type TargetDir string

func (it TargetDir) CopyBack(fullpath, relativepath string, details os.FileInfo) {
	targetFile := filepath.Join(string(it), relativepath)
	err := pathlib.CopyFile(fullpath, targetFile, false)
	if err != nil {
		common.Log("Warning %v while copying %v", err, targetFile)
	}
}

func (it TargetDir) OverwriteBack(fullpath, relativepath string, details os.FileInfo) {
	targetFile := filepath.Join(string(it), relativepath)
	err := pathlib.CopyFile(fullpath, targetFile, true)
	if err != nil {
		common.Log("Warning %v while copying %v", err, targetFile)
	}
}

func init() {
	taskCmd.AddCommand(testrunCmd)

	testrunCmd.Flags().StringArrayVarP(&ignores, "ignore", "i", []string{}, "File with ignore patterns.")
	testrunCmd.Flags().StringVarP(&environmentFile, "environment", "e", "", "Full path to the 'env.json' development environment data file.")
	testrunCmd.Flags().StringVarP(&robotFile, "robot", "r", "robot.yaml", "Full path to the 'robot.yaml' configuration file. (Backward compatibility with 'package.yaml')")
	testrunCmd.Flags().StringVarP(&runTask, "task", "t", "", "Task to run from configuration file.")
	testrunCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "Optional workspace id to get authorization tokens for. OPTIONAL")
	testrunCmd.Flags().IntVarP(&validityTime, "minutes", "m", 0, "How many minutes the authorization should be valid for. OPTIONAL")
	testrunCmd.Flags().StringVarP(&accountName, "account", "", "", "Account used for workspace. OPTIONAL")
	testrunCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force conda cache update. (only for new environments)")
}
