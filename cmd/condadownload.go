package cmd

import (
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var condaDownloadCmd = &cobra.Command{
	Use:     "download",
	Aliases: []string{"dl", "d"},
	Short:   "Download the miniconda3 installer.",
	Long:    `Downloads the miniconda3 installer for this platform.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if !(conda.DoDownload() || conda.DoDownload() || conda.DoDownload()) {
			pretty.Exit(1, "Download failed.")
		}
	},
}

func init() {
	condaCmd.AddCommand(condaDownloadCmd)
}
