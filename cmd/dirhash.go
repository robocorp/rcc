package cmd

import (
	"os"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"

	"github.com/spf13/cobra"
)

var dirhashCmd = &cobra.Command{
	Use:   "dirhash",
	Short: "Calculate hash for directory content.",
	Long:  `Calculate SHA256 of directory tree structure.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer common.Stopwatch("rcc dirhash lasted").Report()
		for _, directory := range args {
			stat, err := os.Stat(directory)
			if err != nil {
				common.Error("dirhash", err)
				continue
			}
			if !stat.IsDir() {
				continue
			}
			digest, err := conda.DigestFor(directory)
			if err != nil {
				common.Error("dirhash", err)
				continue
			}
			result := conda.Hexdigest(digest)
			common.Log("+ %v %v", result, directory)
		}
	},
}

func init() {
	internalCmd.AddCommand(dirhashCmd)
}
