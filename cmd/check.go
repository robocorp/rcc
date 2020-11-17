package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"c"},
	Short:   "Check if conda is installed in managed location.",
	Long: `Check if conda is installed. And optionally also force download and install
conda using "rcc conda download" and "rcc conda install" commands.  `,
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Conda check took").Report()
		}
		if conda.HasConda() {
			pretty.Exit(0, "OK.")
		}
		common.Debug("Conda is missing ...")
		if !autoInstall {
			pretty.Exit(1, "Error: No conda.")
		}
		common.Debug("Starting conda download ...")
		if !(conda.DoDownload() || conda.DoDownload() || conda.DoDownload()) {
			pretty.Exit(2, "Error: Conda download failed.")
		}
		common.Debug("Starting conda install ...")
		if !conda.DoInstall() {
			pretty.Exit(3, "Error: Conda install failed.")
		}
		common.Debug("Conda install completed ...")
		pretty.Ok()
	},
}

func init() {
	condaCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolVarP(&autoInstall, "install", "i", false, "If conda is missing, download and install it automatically.")
}
