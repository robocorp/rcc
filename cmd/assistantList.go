package cmd

import (
	"encoding/json"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var assistantListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Robot Assistant listing",
	Long:    "Robot Assistant listing.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Robot Assistant list query lasted").Report()
		}
		account := operations.AccountByName(AccountName())
		if account == nil {
			pretty.Exit(1, "Could not find account by name: %v", AccountName())
		}
		client, err := cloud.NewClient(account.Endpoint)
		if err != nil {
			pretty.Exit(2, "Could not create client for endpoint: %v, reason: %v", account.Endpoint, err)
		}
		assistants, err := operations.ListAssistantsCommand(client, account, workspaceId)
		if err != nil {
			pretty.Exit(3, "Could not get list of assistants for workspace %v, reason: %v", workspaceId, err)
		}
		nice, err := json.MarshalIndent(assistants, "", "  ")
		if err != nil {
			pretty.Exit(4, "Could not format reply: %v", err)
		}
		common.Out("%s", nice)
	},
}

func init() {
	assistantCmd.AddCommand(assistantListCmd)
	assistantListCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "Workspace id to get assistant information.")
	assistantListCmd.MarkFlagRequired("workspace")
}
