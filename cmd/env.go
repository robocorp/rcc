package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/settings"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:     "env",
	Aliases: []string{"environment", "e"},
	Short:   "Group of commands related to `environment management`.",
	Long: `This "env" command set is for managing virtual environments,
used in task context locally.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		settings.CriticalEnvironmentSettingsCheck()
	},
}

func init() {

	rootCmd.AddCommand(envCmd)
	envCmd.PersistentFlags().StringVar(&common.StageFolder, "stage", "", "internal, DO NOT USE (unless you know what you are doing)")
}
