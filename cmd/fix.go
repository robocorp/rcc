package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Automatically fix known issues inside robots.",
	Long: `Automatically fix known issues inside robots. Current fixes are:
- make files in PATH folder executable
- convert .sh newlines to unix form`,
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Fix run lasted").Report()
		}
		err := operations.FixRobot(robotFile)
		if err != nil {
			pretty.Exit(1, "Error: %v", err)
		}
		pretty.Ok()
	},
}

func init() {
	robotCmd.AddCommand(fixCmd)
	fixCmd.Flags().StringVarP(&robotFile, "robot", "r", "robot.yaml", "Full path to 'robot.yaml' configuration file.")
}
