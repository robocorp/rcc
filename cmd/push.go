package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Wrap the local directory and push it into Robocorp Cloud as a specific robot.",
	Long:  "Wrap the local directory and push it into Robocorp Cloud as a specific robot.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Push lasted").Report()
		}
		account := operations.AccountByName(AccountName())
		if account == nil {
			pretty.Exit(1, "Could not find account by name: %v", AccountName())
		}
		client, err := cloud.NewClient(account.Endpoint)
		if err != nil {
			pretty.Exit(2, "Could not create client for endpoint: %v reason: %v", account.Endpoint, err)
		}

		zipfile := filepath.Join(os.TempDir(), fmt.Sprintf("push%x.zip", time.Now().Unix()))
		defer os.Remove(zipfile)
		common.Debug("Using temporary zipfile at %v", zipfile)

		err = operations.Zip(directory, zipfile, ignores)
		if err != nil {
			pretty.Exit(3, "Error: %v", err)
		}
		err = operations.UploadCommand(client, account, workspaceId, robotId, zipfile, common.DebugFlag)
		if err != nil {
			pretty.Exit(4, "Error: %v", err)
		}
		operations.BackgroundMetric("rcc", "rcc.cli.push", common.Version)
		pretty.Ok()
	},
}

func init() {
	cloudCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringVarP(&directory, "directory", "d", ".", "The root directory to create the robot from.")
	pushCmd.Flags().StringArrayVarP(&ignores, "ignore", "i", []string{}, "Files containing ignore patterns.")
	pushCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "The workspace id to use as the upload target.")
	pushCmd.MarkFlagRequired("workspace")
	pushCmd.Flags().StringVarP(&robotId, "robot", "r", "", "The robot id to use as the upload target.")
	pushCmd.MarkFlagRequired("robot")
}
