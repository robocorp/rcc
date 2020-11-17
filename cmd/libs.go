package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var (
	channelFlag bool
	pipFlag     bool
	dryFlag     bool

	condaOption string
	nameOption  string
	addMany     []string
	removeMany  []string
)

var libsCmd = &cobra.Command{
	Use:     "libs",
	Aliases: []string{"library", "libraries"},
	Short:   "Manage library dependencies in an action oriented way.",
	Long:    "Manage library dependencies in an action oriented way.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Robot libs lasted").Report()
		}
		changes := &conda.Changes{
			Name:    nameOption,
			Pip:     pipFlag,
			Dryrun:  dryFlag,
			Channel: channelFlag,
			Add:     addMany,
			Remove:  removeMany,
		}
		output, err := conda.UpdateEnvironment(condaOption, changes)
		if err != nil {
			pretty.Exit(1, "Error: %v", err)
		}
		common.Out("%s", output)
		pretty.Ok()
	},
}

func init() {
	robotCmd.AddCommand(libsCmd)
	libsCmd.Flags().StringVarP(&nameOption, "name", "n", "", "Change the name of the configuration.")
	libsCmd.Flags().StringVarP(&condaOption, "conda", "", "", "Full path to the conda environment configuration file (conda.yaml).")
	libsCmd.MarkFlagRequired("conda")
	libsCmd.Flags().StringArrayVarP(&addMany, "add", "a", []string{}, "Add new libraries as requirements.")
	libsCmd.Flags().StringArrayVarP(&removeMany, "remove", "r", []string{}, "Remove existing libraries from requirements.")
	libsCmd.Flags().BoolVarP(&channelFlag, "channel", "c", false, "Operate on channels (default is packages).")
	libsCmd.Flags().BoolVarP(&pipFlag, "pip", "p", false, "Operate on pip packages (the default is to operate on conda packages).")
	libsCmd.Flags().BoolVarP(&dryFlag, "dryrun", "d", false, "Do not save the end result, just show what would happen.")
}
