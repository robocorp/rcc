package cmd

import (
	"github.com/robocorp/rcc/blobs"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var tutorialCmd = &cobra.Command{
	Use:     "tutorial",
	Short:   "Show the rcc tutorial.",
	Long:    "Show the rcc tutorial.",
	Aliases: []string{"tut"},
	Run: func(cmd *cobra.Command, args []string) {
		content, err := blobs.Asset("assets/man/tutorial.txt")
		if err != nil {
			pretty.Exit(1, "Cannot show tutorial text, reason: %v", err)
		}
		common.Out("%s", content)
	},
}

func init() {
	manCmd.AddCommand(tutorialCmd)
	rootCmd.AddCommand(tutorialCmd)
}
