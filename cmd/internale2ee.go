package cmd

import (
	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var e2eeCmd = &cobra.Command{
	Use:   "encryption",
	Short: "Internal end-to-end encryption tester method",
	Long:  "Internal end-to-end encryption tester method",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Encryption lasted").Report()
		}
		account := operations.AccountByName(AccountName())
		if account == nil {
			pretty.Exit(1, "Could not find account by name: %v", AccountName())
		}
		client, err := cloud.NewClient(account.Endpoint)
		if err != nil {
			pretty.Exit(2, "Could not create client for endpoint: %v, reason: %v", account.Endpoint, err)
		}
		key, err := operations.GenerateEphemeralKey()
		if err != nil {
			pretty.Exit(3, "Problem with key generation, reason: %v", err)
		}
		request := client.NewRequest("/assistant-v1/test/encryption")
		request.Body, err = key.RequestBody(args[0])
		if err != nil {
			pretty.Exit(4, "Problem with body generation, reason: %v", err)
		}
		response := client.Post(request)
		if response.Status != 200 {
			pretty.Exit(5, "Problem with test request, status=%d, body=%s", response.Status, response.Body)
		}
		plaintext, err := key.Decode(response.Body)
		if err != nil {
			pretty.Exit(6, "Decode problem with body %s, reason: %v", response.Body, err)
		}
		common.Log("Response: %s", string(plaintext))
		pretty.Ok()
	},
}

func init() {
	internalCmd.AddCommand(e2eeCmd)
	e2eeCmd.Flags().StringVarP(&accountName, "account", "a", "", "Account used for Robocorp Cloud operations.")
}
