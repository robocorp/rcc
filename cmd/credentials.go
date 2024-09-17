package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/settings"

	"github.com/spf13/cobra"
)

var (
	deleteCredentialsFlag bool
)

var credentialsCmd = &cobra.Command{
	Use:   "credentials [credentials]",
	Short: fmt.Sprintf("Manage %s Control Room API credentials.", common.Product.Name()),
	Long:  fmt.Sprintf("Manage %s Control Room API credentials for later use.", common.Product.Name()),
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag() {
			defer common.Stopwatch("Credentials query lasted").Report()
		}
		var account, credentials, endpoint string
		if len(args) == 1 {
			credentials = strings.TrimSpace(args[0])
		}
		show := len(credentials) == 0
		if show && verifiedFlag {
			operations.VerifyAccounts(forceFlag)
		}
		if show && !deleteCredentialsFlag {
			operations.ListAccounts(jsonFlag, secretsFlag)
			return
		}
		account = strings.TrimSpace(AccountName())
		if len(account) == 0 {
			account = "Default account"
		}
		if deleteCredentialsFlag {
			localDelete(account)
		}
		endpoint = endpointUrl
		if len(endpoint) == 0 {
			endpoint = settings.Global.DefaultEndpoint()
		}
		https, err := cloud.EnsureHttps(endpoint)
		if err != nil {
			pretty.Exit(1, "Error: %v", err)
		}
		parts := strings.Split(credentials, ":")
		if len(parts) != 2 {
			pretty.Exit(1, "Error: No valid credentials detected. Copy them from %s Control Room.", common.Product.Name())
		}
		common.Log("Adding credentials: %v", parts)
		operations.UpdateCredentials(account, https, parts[0], parts[1])
		if defaultFlag {
			operations.SetDefaultAccount(account)
		}
		pretty.Ok()
	},
}

func localDelete(accountName string) {
	account := operations.AccountByName(accountName)
	if account == nil {
		pretty.Exit(1, "Could not find account by name: %q", accountName)
	}
	err := account.Delete(10 * time.Second)
	if err != nil {
		pretty.Exit(3, "Error: %v", err)
	}
	pretty.Exit(0, "OK.")
}

func init() {
	configureCmd.AddCommand(credentialsCmd)

	credentialsCmd.Flags().BoolVarP(&deleteCredentialsFlag, "delete", "", false, "Delete this account and corresponding Control Room credentials! DANGER!")
	credentialsCmd.Flags().BoolVarP(&defaultFlag, "default", "d", false, "Set this as the default account.")
	credentialsCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output in JSON format.")
	credentialsCmd.Flags().BoolVarP(&verifiedFlag, "verified", "v", false, "Updates the verified timestamp, if the credentials are still active.")
	credentialsCmd.Flags().StringVarP(&endpointUrl, "endpoint", "e", "", fmt.Sprintf("%s Control Room endpoint used with the given account (or default).", common.Product.Name))
	credentialsCmd.Flags().BoolVarP(&secretsFlag, "showSecrets", "", false, "Caution: When used with --json this exposes the credential secret for clients that need them.")
}
