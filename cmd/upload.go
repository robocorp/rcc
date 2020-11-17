package cmd

import (
	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Push an existing robot to Robocorp Cloud.",
	Long:  "Push an existing robot to Robocorp Cloud.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Upload lasted").Report()
		}
		account := operations.AccountByName(AccountName())
		if account == nil {
			pretty.Exit(1, "Could not find account by name: %v", AccountName())
		}
		client, err := cloud.NewClient(account.Endpoint)
		if err != nil {
			pretty.Exit(2, "Could not create client for endpoint: %v, reason: %v", account.Endpoint, err)
		}
		err = operations.UploadCommand(client, account, workspaceId, robotId, zipfile, common.DebugFlag)
		if err != nil {
			pretty.Exit(3, "Error: %v", err)
		}
		pretty.Ok()
	},
}

func init() {
	cloudCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&zipfile, "zipfile", "z", "robot.zip", "The filename for the robot.")
	uploadCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "The workspace id to use as the upload target.")
	uploadCmd.MarkFlagRequired("workspace")
	uploadCmd.Flags().StringVarP(&robotId, "robot", "r", "", "The robot id to use as the upload target.")
	uploadCmd.MarkFlagRequired("robot")
}
