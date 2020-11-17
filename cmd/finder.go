package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var finderCmd = &cobra.Command{
	Use:     "finder",
	Aliases: []string{"f"},
	Short:   "Internal command to breadth-first finding of named file.",
	Long: `Internal tool for finding some named file. It uses breadth first stragety
and accepts only one finding on found level. If there is more than one result,
then that is an error.

Example:
    rcc internal finder -d /starting/path/somewhere robot.yaml`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Finder run lasted").Report()
		}
		found, err := pathlib.FindNamedPath(shellDirectory, args[0])
		if err != nil {
			pretty.Exit(1, "Error: %v", err)
		} else {
			common.Out("%s", found)
		}
	},
}

func init() {
	internalCmd.AddCommand(finderCmd)

	finderCmd.Flags().StringVarP(&shellDirectory, "directory", "d", ".", "The working directory for the debug command.")
}
