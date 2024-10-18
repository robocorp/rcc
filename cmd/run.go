package cmd

import (
	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/journal"
	"github.com/robocorp/rcc/operations"

	"github.com/spf13/cobra"
)

var (
	rcHosts         = []string{"RC_API_SECRET_HOST", "RC_API_WORKITEM_HOST"}
	rcTokens        = []string{"RC_API_SECRET_TOKEN", "RC_API_WORKITEM_TOKEN"}
	interactiveFlag bool
)

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Short:   "Run task in place, to debug current setup.",
	Long: `Local task run, in place, to see how full run execution works
in your own machine.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer conda.RemoveCurrentTemp()
		defer journal.BuildEventStats("robot")
		defer journal.StopRunJournal()
		if common.DebugFlag() {
			defer common.Stopwatch("Task run lasted").Report()
		}
		simple, config, todo, label := operations.LoadTaskWithEnvironment(robotFile, runTask, forceFlag)
		cloud.InternalBackgroundMetric(common.ControllerIdentity(), "rcc.cli.run", common.Version)
		commandline := todo.Commandline()
		commandline = append(commandline, args...)
		operations.SelectExecutionModel(captureRunFlags(false), simple, commandline, config, todo, label, interactiveFlag, nil)
	},
}

func captureRunFlags(assistant bool) *operations.RunFlags {
	return &operations.RunFlags{
		TokenPeriod: &operations.TokenPeriod{
			ValidityTime: validityTime,
			GracePeriod:  gracePeriod,
		},
		AccountName:     AccountName(),
		WorkspaceId:     workspaceId,
		EnvironmentFile: environmentFile,
		RobotYaml:       robotFile,
		Assistant:       assistant,
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
	taskCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&environmentFile, "environment", "e", "", "Full path to the 'env.json' development environment data file.")
	runCmd.Flags().StringVarP(&robotFile, "robot", "r", "robot.yaml", "Full path to the 'robot.yaml' configuration file.")
	runCmd.Flags().StringVarP(&runTask, "task", "t", "", "Task to run from the configuration file.")
	runCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "Optional workspace id to get authorization tokens for. OPTIONAL")
	runCmd.Flags().IntVarP(&validityTime, "minutes", "m", 15, "How many minutes the authorization should be valid for (minimum 15 minutes).")
	runCmd.Flags().IntVarP(&gracePeriod, "graceperiod", "", 5, "What is grace period buffer in minutes on top of validity minutes (minimum 5 minutes).")
	runCmd.Flags().StringVarP(&accountName, "account", "", "", "Account used for workspace. OPTIONAL")
	runCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force conda cache update (only for new environments).")
	runCmd.Flags().BoolVarP(&interactiveFlag, "interactive", "", false, "Allow robot to be interactive in terminal/command prompt. For development only, not for production!")
	runCmd.Flags().StringVarP(&common.HolotreeSpace, "space", "s", "user", "Client specific name to identify this environment.")
	runCmd.Flags().BoolVarP(&common.NoOutputCapture, "no-outputs", "", false, "Do not capture stderr/stdout into files.")
	runCmd.Flags().BoolVarP(&common.DeveloperFlag, "dev", "", false, "Use devTasks instead of normal tasks. For development work only. Strategy selection.")
}
