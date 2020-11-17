package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var unwrapCmd = &cobra.Command{
	Use:   "unwrap",
	Short: "Unpack a robot back into a directory structure.",
	Long: `Unpack a robot back into a directory structure. This command expects to get
robot filename, and target directory. And using --force option, files will
be overwritten.`,
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Unwrap lasted").Report()
		}
		err := operations.Unzip(directory, zipfile, forceFlag, false)
		if err != nil {
			pretty.Exit(1, "Error: %v", err)
		}
		pretty.Ok()
	},
}

func init() {
	robotCmd.AddCommand(unwrapCmd)
	unwrapCmd.Flags().StringVarP(&zipfile, "zipfile", "z", "robot.zip", "The filename for the robot to extract.")
	unwrapCmd.Flags().StringVarP(&directory, "directory", "d", "", "The root directory to extract robot into.")
	unwrapCmd.MarkFlagRequired("directory")
	unwrapCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Remove the safety nets around the unwrapping of the robot.")
}
