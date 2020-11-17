package operations

import (
	"fmt"
	"path/filepath"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/robot"
	"github.com/robocorp/rcc/shell"
)

var (
	rcHosts  = []string{"RC_API_SECRET_HOST", "RC_API_WORKITEM_HOST"}
	rcTokens = []string{"RC_API_SECRET_TOKEN", "RC_API_WORKITEM_TOKEN"}
)

type RunFlags struct {
	AccountName     string
	WorkspaceId     string
	ValidityTime    int
	EnvironmentFile string
	RobotYaml       string
}

func PipFreeze(searchPath pathlib.PathParts, directory, outputDir string, environment []string) bool {
	pip, ok := searchPath.Which("pip", conda.FileExtensions)
	if !ok {
		return false
	}
	fullPip, err := filepath.EvalSymlinks(pip)
	if err != nil {
		return false
	}
	common.Log("Installed pip packages:")
	_, err = shell.New(environment, directory, fullPip, "freeze", "--all").Tee(outputDir, false)
	if err != nil {
		return false
	}
	common.Log("--")
	return true
}

func LoadTaskWithEnvironment(packfile, theTask string, force bool) (bool, robot.Robot, robot.Task, string) {
	FixRobot(packfile)
	config, err := robot.LoadYamlConfiguration(packfile)
	if err != nil {
		pretty.Exit(1, "Error: %v", err)
	}

	ok, err := config.Validate()
	if !ok {
		pretty.Exit(2, "Error: %v", err)
	}

	todo := config.TaskByName(theTask)
	if todo == nil {
		pretty.Exit(3, "Error: Could not resolve task to run. Available tasks are: %v", config.AvailableTasks())
	}

	if !config.UsesConda() {
		return true, config, todo, ""
	}

	label, err := conda.NewEnvironment(force, config.CondaConfigFile())
	if err != nil {
		pretty.Exit(4, "Error: %v", err)
	}
	return false, config, todo, label
}

func SelectExecutionModel(runFlags *RunFlags, simple bool, template []string, config robot.Robot, todo robot.Task, label string, interactive bool, extraEnv map[string]string) {
	if simple {
		ExecuteSimpleTask(runFlags, template, config, todo, interactive, extraEnv)
	} else {
		ExecuteTask(runFlags, template, config, todo, label, interactive, extraEnv)
	}
}

func ExecuteSimpleTask(flags *RunFlags, template []string, config robot.Robot, todo robot.Task, interactive bool, extraEnv map[string]string) {
	common.Debug("Command line is: %v", template)
	task := make([]string, len(template))
	copy(task, template)
	searchPath := pathlib.TargetPath()
	searchPath = searchPath.Prepend(todo.Paths(config)...)
	found, ok := searchPath.Which(task[0], conda.FileExtensions)
	if !ok {
		pretty.Exit(6, "Error: Cannot find command: %v", task[0])
	}
	fullpath, err := filepath.EvalSymlinks(found)
	if err != nil {
		pretty.Exit(7, "Error: %v", err)
	}
	var data Token
	if len(flags.WorkspaceId) > 0 {
		claims := RunClaims(flags.ValidityTime*60, flags.WorkspaceId)
		data, err = AuthorizeClaims(flags.AccountName, claims)
	}
	if err != nil {
		pretty.Exit(8, "Error: %v", err)
	}
	task[0] = fullpath
	directory := todo.WorkingDirectory(config)
	environment := robot.PlainEnvironment([]string{searchPath.AsEnvironmental("PATH")}, true)
	if len(data) > 0 {
		endpoint := data["endpoint"]
		for _, key := range rcHosts {
			environment = append(environment, fmt.Sprintf("%s=%s", key, endpoint))
		}
		token := data["token"]
		for _, key := range rcTokens {
			environment = append(environment, fmt.Sprintf("%s=%s", key, token))
		}
		environment = append(environment, fmt.Sprintf("RC_WORKSPACE_ID=%s", flags.WorkspaceId))
	}
	if extraEnv != nil {
		for key, value := range extraEnv {
			environment = append(environment, fmt.Sprintf("%s=%s", key, value))
		}
	}
	outputDir := todo.ArtifactDirectory(config)
	common.Debug("DEBUG: about to run command - %v", task)
	_, err = shell.New(environment, directory, task...).Tee(outputDir, interactive)
	if err != nil {
		pretty.Exit(9, "Error: %v", err)
	}
	pretty.Ok()
}

func ExecuteTask(flags *RunFlags, template []string, config robot.Robot, todo robot.Task, label string, interactive bool, extraEnv map[string]string) {
	common.Debug("Command line is: %v", template)
	developmentEnvironment, err := robot.LoadEnvironmentSetup(flags.EnvironmentFile)
	if err != nil {
		pretty.Exit(5, "Error: %v", err)
	}
	task := make([]string, len(template))
	copy(task, template)
	searchPath := todo.SearchPath(config, label)
	found, ok := searchPath.Which(task[0], conda.FileExtensions)
	if !ok {
		pretty.Exit(6, "Error: Cannot find command: %v", task[0])
	}
	fullpath, err := filepath.EvalSymlinks(found)
	if err != nil {
		pretty.Exit(7, "Error: %v", err)
	}
	var data Token
	if len(flags.WorkspaceId) > 0 {
		claims := RunClaims(flags.ValidityTime*60, flags.WorkspaceId)
		data, err = AuthorizeClaims(flags.AccountName, claims)
	}
	if err != nil {
		pretty.Exit(8, "Error: %v", err)
	}
	task[0] = fullpath
	directory := todo.WorkingDirectory(config)
	environment := todo.ExecutionEnvironment(config, label, developmentEnvironment.AsEnvironment(), true)
	if len(data) > 0 {
		endpoint := data["endpoint"]
		for _, key := range rcHosts {
			environment = append(environment, fmt.Sprintf("%s=%s", key, endpoint))
		}
		token := data["token"]
		for _, key := range rcTokens {
			environment = append(environment, fmt.Sprintf("%s=%s", key, token))
		}
		environment = append(environment, fmt.Sprintf("RC_WORKSPACE_ID=%s", flags.WorkspaceId))
	}
	if extraEnv != nil {
		for key, value := range extraEnv {
			environment = append(environment, fmt.Sprintf("%s=%s", key, value))
		}
	}
	outputDir := todo.ArtifactDirectory(config)
	if !common.Silent && !interactive {
		PipFreeze(searchPath, directory, outputDir, environment)
	}
	common.Debug("DEBUG: about to run command - %v", task)
	_, err = shell.New(environment, directory, task...).Tee(outputDir, interactive)
	if err != nil {
		pretty.Exit(9, "Error: %v", err)
	}
	pretty.Ok()
}
