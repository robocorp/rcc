package cmd

import (
	"encoding/json"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "List the available workspaces and their tasks (with --workspace option).",
	Long:  "List the available workspaces and their tasks (with --workspace option).",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Workspace query lasted").Report()
		}
		account := operations.AccountByName(AccountName())
		if account == nil {
			pretty.Exit(1, "Could not find account by name: %v", AccountName())
		}
		client, err := cloud.NewClient(account.Endpoint)
		if err != nil {
			pretty.Exit(2, "Could not create client for endpoint: %v, reason: %v", account.Endpoint, err)
		}

		var data interface{}
		if len(workspaceId) > 0 {
			data, err = operations.WorkspaceTreeCommand(client, account, workspaceId)
		} else {
			data, err = operations.WorkspacesCommand(client, account)
		}

		if err != nil {
			pretty.Exit(3, "Could not receive workspace data: %v", err)
		}
		nice, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			pretty.Exit(3, "Could not format reply: %v", err)
		}
		common.Out("%s", nice)
	},
}

func init() {
	cloudCmd.AddCommand(workspaceCmd)
	workspaceCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "The id of the workspace to use with this command.")
	workspaceCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output in JSON format")
}
