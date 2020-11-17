package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var lshCmd = &cobra.Command{
	Use:   "lsh",
	Short: "Locality-sensitive hash calculation",
	Long: `This lsh command calculates locality-sensitive hash from environment.yaml,
or requirements.txt, or similar text files.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		baseline := ""
		failure := false
		for _, arg := range args {
			out, err := conda.HashConfig(arg)
			if err != nil {
				common.Error("lsh", err)
				failure = true
				continue
			}
			if baseline == "" {
				baseline = out
			}
			distance, _ := conda.Distance(baseline, out)
			common.Log("%s: %s <%d>", out, arg, distance)
		}
		if failure {
			pretty.Exit(1, "Error!")
		}
	},
}

func init() {
	internalCmd.AddCommand(lshCmd)
}
