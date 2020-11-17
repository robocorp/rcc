package cmd

import (
	"encoding/json"
	"os"
	"time"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var (
	taskDirectory string
)

const (
	timeFormat = `02.01.2006 15:04`
)

func updateRobotDirectory(directory string) {
	err := operations.UpdateRobot(directory)
	if err != nil {
		pretty.Exit(1, "Error: %v", err)
	}
}

func jsonRobots() {
	robots, err := operations.ListRobots()
	if err != nil {
		pretty.Exit(1, "Error: %v", err)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(robots)
	if err != nil {
		pretty.Exit(2, "Error: %v", err)
	}
}

func listRobots() {
	if jsonFlag {
		jsonRobots()
		return
	}
	robots, err := operations.ListRobots()
	if err != nil {
		pretty.Exit(1, "Error: %v", err)
	}
	if len(robots) == 0 {
		pretty.Exit(2, "Error: No robots found!")
	}
	common.Log("Updated at       | Created at       | Directory")
	for _, robot := range robots {
		updated := time.Unix(robot.Updated, 0)
		created := time.Unix(robot.Created, 0)
		status := ""
		if robot.Deleted > 0 {
			status = "<deleted>"
		}
		common.Log("%v | %v | %v %v", updated.Format(timeFormat), created.Format(timeFormat), robot.Path, status)
	}
}

var robotlistCmd = &cobra.Command{
	Use:   "list",
	Short: "List or update tracked robot directories.",
	Long:  "List or update tracked robot directories.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Robot list lasted").Report()
		}
		if len(taskDirectory) > 0 {
			updateRobotDirectory(taskDirectory)
		} else {
			listRobots()
		}
	},
}

func init() {
	robotCmd.AddCommand(robotlistCmd)
	robotlistCmd.Flags().StringVarP(&taskDirectory, "add", "a", "", "The root directory to add as robot.")
	robotlistCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output in JSON format.")
}
