package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/xviper"

	"github.com/spf13/cobra"
)

var (
	metricType  string
	metricName  string
	metricValue string
)

var metricCmd = &cobra.Command{
	Use:   "metric",
	Short: "Send some metric to Robocorp Cloud.",
	Long:  "Send some metric to Robocorp Cloud.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Feedback metric lasted").Report()
		}
		if !xviper.CanTrack() {
			pretty.Exit(1, "Tracking is disabled. Quitting.")
		}
		operations.SendMetric(metricType, metricName, metricValue)
		pretty.Exit(0, "OK")
	},
}

func init() {
	feedbackCmd.AddCommand(metricCmd)
	metricCmd.Flags().StringVarP(&metricType, "type", "t", "", "Type for metric source to use.")
	metricCmd.MarkFlagRequired("type")
	metricCmd.Flags().StringVarP(&metricName, "name", "n", "", "Name for metric to report.")
	metricCmd.MarkFlagRequired("name")
	metricCmd.Flags().StringVarP(&metricValue, "value", "v", "", "Value for metric to report.")
	metricCmd.MarkFlagRequired("value")
}
