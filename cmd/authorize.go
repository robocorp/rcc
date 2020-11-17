package cmd

import (
	"encoding/json"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Convert an API key to a valid authorization JWT token.",
	Long:  "Convert an API key to a valid authorization JWT token.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Authorize query lasted").Report()
		}
		var claims *operations.Claims
		if granularity == "user" {
			claims = operations.WorkspaceTreeClaims(validityTime * 60)
		} else {
			claims = operations.RunClaims(validityTime*60, workspaceId)
		}
		data, err := operations.AuthorizeClaims(AccountName(), claims)
		if err != nil {
			pretty.Exit(3, "Error: %v", err)
		}
		nice, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			pretty.Exit(4, "Error: Could not format reply: %v", err)
		}
		common.Out("%s", nice)
	},
}

func init() {
	cloudCmd.AddCommand(authorizeCmd)
	authorizeCmd.Flags().IntVarP(&validityTime, "minutes", "m", 0, "How many minutes the authorization should be valid for.")
	authorizeCmd.Flags().StringVarP(&granularity, "granularity", "g", "", "Authorization granularity (user/workspace) used in.")
	authorizeCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "Workspace id to use with this command.")
}
