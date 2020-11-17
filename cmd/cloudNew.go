package cmd

import (
	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var newCloudCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new robot into Robocorp Cloud.",
	Long:  "Create a new robot into Robocorp Cloud.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("New robot creation lasted").Report()
		}
		account := operations.AccountByName(AccountName())
		if account == nil {
			pretty.Exit(1, "Could not find account by name: %v", AccountName())
		}
		client, err := cloud.NewClient(account.Endpoint)
		if err != nil {
			pretty.Exit(2, "Could not create client for endpoint: %v, reason: %v", account.Endpoint, err)
		}
		reply, err := operations.NewRobotCommand(client, account, workspaceId, robotName)
		if err != nil {
			pretty.Exit(3, "Error: %v", err)
		}
		common.Log("Created new robot named '%s' with identity %s.", reply["name"], reply["id"])
	},
}

func init() {
	cloudCmd.AddCommand(newCloudCmd)
	newCloudCmd.Flags().StringVarP(&robotName, "robot", "r", "", "Name for new robot to create.")
	newCloudCmd.MarkFlagRequired("robot")
	newCloudCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "Workspace id to use as creation target.")
	newCloudCmd.MarkFlagRequired("workspace")
}
