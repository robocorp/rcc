package cmd

import (
	"github.com/robocorp/rcc/common"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show rcc version number.",
	Long:    `Show current version number of installed rcc.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		common.Out("%s", common.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
