package cmd

import (
	"github.com/robocorp/rcc/blobs"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Show the rcc License.",
	Long:  "Show the rcc License.",
	Run: func(cmd *cobra.Command, args []string) {
		content, err := blobs.Asset("assets/man/LICENSE.txt")
		if err != nil {
			pretty.Exit(1, "Cannot show LICENSE, reason: %v", err)
		}
		common.Out("%s", content)
	},
}

func init() {
	manCmd.AddCommand(licenseCmd)
}
