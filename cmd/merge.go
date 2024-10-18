package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

func dumpYaml(title string, environment *conda.Environment) {
	common.Log("%s", title)
	content, err := environment.AsYaml()
	if err != nil {
		pretty.Exit(3, err.Error())
	}
	common.Stdout("%s\n", content)
}

var mergeCmd = &cobra.Command{
	Use:   "merge conda.yaml+",
	Short: "Tool for testing conda.yaml merging.",
	Long:  "Tool for testing conda.yaml merging.",
	Run: func(cmd *cobra.Command, args []string) {
		var left, right *conda.Environment
		var err error

		for _, filename := range args {
			left = right
			right, err = conda.ReadPackageCondaYaml(filename, false)
			if err != nil {
				pretty.Exit(1, err.Error())
			}
			if left == nil {
				continue
			}
			right, err = left.Merge(right)
			if err != nil {
				pretty.Exit(2, err.Error())
			}
		}
		dumpYaml("Full merge:", right)
		dumpYaml("Pure conda:", right.AsPureConda())
		common.Log("Requirements text:")
		common.Stdout("%s\n", right.AsRequirementsText())
	},
}

func init() {
	internalCmd.AddCommand(mergeCmd)
}
