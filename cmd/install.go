package cmd

import (
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "Install miniconda into the managed location.",
	Long: `Install miniconda into the rcc managed location. Before executing this command,
you must successfully run the "download" command and verify that the miniconda SHA256
matches the one on the conda site.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !conda.DoInstall() {
			pretty.Exit(1, "Error: Install failed. See above.")
		}
	},
}

func init() {
	condaCmd.AddCommand(installCmd)
}
