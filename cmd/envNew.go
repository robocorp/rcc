package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var newEnvCmd = &cobra.Command{
	Use:   "new <conda.yaml+>",
	Short: "Creates a new managed virtual environment.",
	Long: `The new command can be used to create a new managed virtual environment.
When given multiple conda.yaml files, they will be merged together and the
end result will be a composite environment.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("New environment creation lasted").Report()
		}
		ok := conda.MustConda()
		if !ok {
			pretty.Exit(2, "Could not get miniconda installed.")
		}
		label, err := conda.NewEnvironment(forceFlag, args...)
		if err != nil {
			pretty.Exit(1, "Environment creation failed: %v", err)
		} else {
			common.Log("Environment for %v as %v created.", args, label)
		}
		if common.Silent {
			common.Out("%s", label)
		}
	},
}

func init() {
	envCmd.AddCommand(newEnvCmd)
	newEnvCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force conda cache update. (only for new environments)")
}
