package cmd

import (
	"fmt"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/htfs"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"
	"github.com/spf13/cobra"
)

func holotreeExpandBlueprint(userFiles []string, packfile string) map[string]interface{} {
	result := make(map[string]interface{})

    devDependencies := false
	_, holotreeBlueprint, err := htfs.ComposeFinalBlueprint(userFiles, packfile, devDependencies)
	pretty.Guard(err == nil, 5, "%s", err)

	common.Debug("FINAL blueprint:\n%s", string(holotreeBlueprint))

	tree, err := htfs.New()
	pretty.Guard(err == nil, 6, "%s", err)

	result["hash"] = common.BlueprintHash(holotreeBlueprint)
	result["exist"] = tree.HasBlueprint(holotreeBlueprint)

	return result
}

var holotreeBlueprintCmd = &cobra.Command{
	Use:     "blueprint conda.yaml+",
	Short:   "Verify that resulting blueprint is in hololibrary.",
	Long:    "Verify that resulting blueprint is in hololibrary.",
	Aliases: []string{"bp"},
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag() {
			defer common.Stopwatch("Holotree blueprints command lasted").Report()
		}

		status := holotreeExpandBlueprint(args, robotFile)
		if holotreeJson {
			out, err := operations.NiceJsonOutput(status)
			pretty.Guard(err == nil, 6, "%s", err)
			fmt.Println(out)
		} else {
			common.Log("Blueprint %q is available: %v", status["hash"], status["exist"])
		}
	},
}

func init() {
	holotreeCmd.AddCommand(holotreeBlueprintCmd)
	holotreeBlueprintCmd.Flags().StringVarP(&robotFile, "robot", "r", "robot.yaml", "Full path to 'robot.yaml' configuration file. <optional>")
	holotreeBlueprintCmd.Flags().BoolVarP(&holotreeJson, "json", "j", false, "Show environment as JSON.")
}
