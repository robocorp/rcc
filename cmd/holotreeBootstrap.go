package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/htfs"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/robot"
	"github.com/spf13/cobra"
)

var (
	holotreeQuick bool
)

func updateEnvironments(robots []string) {
	tree, err := htfs.New()
	pretty.Guard(err == nil, 2, "Holotree creation error: %v", err)
	for at, template := range robots {
		workarea := filepath.Join(pathlib.TempDir(), fmt.Sprintf("workarea%x%x", common.When, at))
		defer os.RemoveAll(workarea)
		common.Debug("Using temporary workarea: %v", workarea)
		err = operations.InitializeWorkarea(workarea, template, false, forceFlag)
		pretty.Guard(err == nil, 2, "Could not create robot %q, reason: %v", template, err)
		targetRobot := robot.DetectConfigurationName(workarea)
		_, blueprint, err := htfs.ComposeFinalBlueprint([]string{}, targetRobot, false)
		if tree.HasBlueprint(blueprint) {
			continue
		}
		config, err := robot.LoadRobotYaml(targetRobot, false)
		pretty.Guard(err == nil, 2, "Could not load robot config %q, reason: %w", targetRobot, err)
		if !config.UsesConda() {
			continue
		}
		_, _, err = htfs.NewEnvironment(config.CondaConfigFile(), "", false, false, operations.PullCatalog)
		pretty.Guard(err == nil, 2, "Holotree recording error: %v", err)
	}
}

var holotreeBootstrapCmd = &cobra.Command{
	Use:     "bootstrap",
	Aliases: []string{"boot"},
	Short:   "Bootstrap holotree from set of templates.",
	Long:    "Bootstrap holotree from set of templates.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag() {
			defer common.Stopwatch("Holotree bootstrap lasted").Report()
		}

		robots := operations.ListTemplates(false)

		if !holotreeQuick {
			updateEnvironments(robots)
		}

		pretty.Ok()
	},
}

func init() {
	holotreeCmd.AddCommand(holotreeBootstrapCmd)
	holotreeBootstrapCmd.Flags().BoolVar(&holotreeQuick, "quick", false, "Do not create environments, just download templates.")
}
