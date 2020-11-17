package cmd

import (
	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Fetch an existing robot from Robocorp Cloud.",
	Long:  "Fetch an existing robot from Robocorp Cloud.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Download lasted").Report()
		}
		account := operations.AccountByName(AccountName())
		if account == nil {
			pretty.Exit(1, "Could not find account by name: %v", AccountName())
		}
		client, err := cloud.NewClient(account.Endpoint)
		if err != nil {
			pretty.Exit(2, "Could not create client for endpoint: %v, reason: %v", account.Endpoint, err)
		}
		err = operations.DownloadCommand(client, account, workspaceId, robotId, zipfile, common.DebugFlag)
		if err != nil {
			pretty.Exit(3, "Error: %v", err)
		}
		pretty.Ok()
	},
}

func init() {
	cloudCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&zipfile, "zipfile", "z", "robot.zip", "The filename for the downloaded robot.")
	downloadCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "The workspace id to use as the download source.")
	downloadCmd.MarkFlagRequired("workspace")
	downloadCmd.Flags().StringVarP(&robotId, "robot", "r", "", "The robot id to use as the download source.")
	downloadCmd.MarkFlagRequired("robot")
}
