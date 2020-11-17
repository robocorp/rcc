package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"sh", "s"},
	Short:   "Run the given command inside the given environment",
	Long: `Shell command executes the given command inside a managed virtual environment.
It can be used to get inside a managed environment and execute your own
command within that environment.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("rcc shell lasted").Report()
		}
		ok := conda.MustConda()
		if !ok {
			pretty.Exit(2, "Could not get miniconda installed.")
		}
		simple, config, todo, label := operations.LoadTaskWithEnvironment(robotFile, runTask, forceFlag)
		if simple {
			pretty.Exit(1, "Cannot do shell for simple execution model.")
		}
		operations.ExecuteTask(captureRunFlags(), conda.Shell, config, todo, label, true, nil)
	},
}

func init() {
	taskCmd.AddCommand(shellCmd)

	shellCmd.Flags().StringVarP(&environmentFile, "environment", "e", "", "Full path to the 'env.json' development environment data file.")
	shellCmd.Flags().StringVarP(&robotFile, "robot", "r", "robot.yaml", "Full path to the 'robot.yaml' configuration file. (With backward compatibility with 'package.yaml')")
	shellCmd.Flags().StringVarP(&runTask, "task", "t", "", "Task to configure shell from configuration file.")
	shellCmd.MarkFlagRequired("config")
}
