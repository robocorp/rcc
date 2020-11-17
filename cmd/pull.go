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

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull a robot from Robocorp Cloud and unwrap it into local directory.",
	Long:  "Pull a robot from Robocorp Cloud and unwrap it into local directory.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Pull lasted").Report()
		}

		account := operations.AccountByName(AccountName())
		if account == nil {
			pretty.Exit(1, "Could not find account by name: %v", AccountName())
		}

		client, err := cloud.NewClient(account.Endpoint)
		if err != nil {
			pretty.Exit(2, "Could not create client for endpoint: %v reason %v", account.Endpoint, err)
		}

		zipfile := filepath.Join(os.TempDir(), fmt.Sprintf("pull%x.zip", time.Now().Unix()))
		defer os.Remove(zipfile)
		common.Debug("Using temporary zipfile at %v", zipfile)

		err = operations.DownloadCommand(client, account, workspaceId, robotId, zipfile, common.DebugFlag)
		if err != nil {
			pretty.Exit(3, "Error: %v", err)
		}

		err = operations.Unzip(directory, zipfile, forceFlag, false)
		if err != nil {
			pretty.Exit(4, "Error: %v", err)
		}

		pretty.Ok()
	},
}

func init() {
	cloudCmd.AddCommand(pullCmd)
	pullCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "The workspace id to use as the download source.")
	pullCmd.MarkFlagRequired("workspace")
	pullCmd.Flags().StringVarP(&robotId, "robot", "r", "", "The robot id to use as the download source.")
	pullCmd.MarkFlagRequired("robot")
	pullCmd.Flags().StringVarP(&directory, "directory", "d", "", "The root directory to extract the robot into.")
	pullCmd.MarkFlagRequired("directory")
	pullCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Remove safety nets around the unwrapping of the robot.")
}
