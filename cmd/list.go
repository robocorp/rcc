package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Listing currently managed virtual environments.",
	Long: `List shows listing of currently managed virtual environments
in human readable form.`,
	Run: func(cmd *cobra.Command, args []string) {
		templates := conda.TemplateList()
		if len(templates) == 0 {
			pretty.Exit(1, "No environments available.")
		}
		lines := make([]string, 0, len(templates))
		common.Log("%-25s  %-25s  %s", "Last used", "Last cloned", "Environment (TLSH)")
		for _, template := range templates {
			cloned := "N/A"
			used := cloned
			when, err := conda.LastUsed(conda.TemplateFrom(template))
			if err == nil {
				cloned = when.Format(time.RFC3339)
			}
			when, err = conda.LastUsed(conda.LiveFrom(template))
			if err == nil {
				used = when.Format(time.RFC3339)
			}
			lines = append(lines, fmt.Sprintf("%-25s  %-25s  %s", used, cloned, template))
		}
		sort.Strings(lines)
		for _, line := range lines {
			common.Log("%s", line)
		}
	},
}

func init() {
	envCmd.AddCommand(listCmd)
}
