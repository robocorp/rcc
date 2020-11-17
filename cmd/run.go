package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/xviper"

	"github.com/spf13/cobra"
)

var (
	rcHosts  = []string{"RC_API_SECRET_HOST", "RC_API_WORKITEM_HOST"}
	rcTokens = []string{"RC_API_SECRET_TOKEN", "RC_API_WORKITEM_TOKEN"}
)

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Short:   "Run task in place, to debug current setup.",
	Long: `Local task run, in place, to see how full run execution works
in your own machine.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Task run lasted").Report()
		}
		ok := conda.MustConda()
		if !ok {
			pretty.Exit(2, "Could not get miniconda installed.")
		}
		defer xviper.RunMinutes().Done()
		simple, config, todo, label := operations.LoadTaskWithEnvironment(robotFile, runTask, forceFlag)
		operations.BackgroundMetric("rcc", "rcc.cli.run", common.Version)
		operations.SelectExecutionModel(captureRunFlags(), simple, todo.Commandline(), config, todo, label, false, nil)
	},
}

func captureRunFlags() *operations.RunFlags {
	return &operations.RunFlags{
		AccountName:     AccountName(),
		WorkspaceId:     workspaceId,
		ValidityTime:    validityTime,
		EnvironmentFile: environmentFile,
		RobotYaml:       robotFile,
	}
}

func init() {
	taskCmd.AddCommand(runCmd)
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&environmentFile, "environment", "e", "", "Full path to the 'env.json' development environment data file.")
	runCmd.Flags().StringVarP(&robotFile, "robot", "r", "robot.yaml", "Full path to the 'robot.yaml' configuration file. (Backward compatibility with 'package.yaml')")
	runCmd.Flags().StringVarP(&runTask, "task", "t", "", "Task to run from the configuration file.")
	runCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "Optional workspace id to get authorization tokens for. OPTIONAL")
	runCmd.Flags().IntVarP(&validityTime, "minutes", "m", 0, "How many minutes the authorization should be valid for. OPTIONAL")
	runCmd.Flags().StringVarP(&accountName, "account", "", "", "Account used for workspace. OPTIONAL")
	runCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force conda cache update (only for new environments).")
}
