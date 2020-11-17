package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/robot"

	"github.com/spf13/cobra"
)

func Has(value string) bool {
	return len(value) > 0
}

func asSimpleMap(line string) map[string]string {
	parts := strings.SplitN(strings.TrimSpace(line), "=", 2)
	if len(parts) != 2 {
		return nil
	}
	result := make(map[string]string)
	result["key"] = parts[0]
	result["value"] = parts[1]
	return result
}

func asJson(items []string) error {
	result := make([]map[string]string, 0, len(items))
	for _, line := range items {
		entry := asSimpleMap(line)
		if entry != nil {
			result = append(result, entry)
		}
	}
	content, err := operations.NiceJsonOutput(result)
	if err != nil {
		return err
	}
	common.Out("%s", content)
	return nil
}

func asText(items []string) {
	for _, line := range items {
		common.Out(line)
	}
}

func exportEnvironment(condaYaml []string, packfile, taskName, environment, workspace string, validity int, jsonform bool) error {
	var err error
	var config robot.Robot
	var task robot.Task
	var extra []string
	var data operations.Token

	if Has(packfile) {
		config, err = robot.LoadYamlConfiguration(packfile)
		if err == nil {
			condaYaml = append(condaYaml, config.CondaConfigFile())
			task = config.TaskByName(taskName)
		}
	}

	if Has(environment) {
		developmentEnvironment, err := robot.LoadEnvironmentSetup(environmentFile)
		if err == nil {
			extra = developmentEnvironment.AsEnvironment()
		}
	}

	if len(condaYaml) < 1 {
		return errors.New("No robot.yaml, package.yaml or conda.yaml files given. Cannot continue.")
	}

	label, err := conda.NewEnvironment(forceFlag, condaYaml...)
	if err != nil {
		return err
	}

	env := conda.EnvironmentExtensionFor(label)
	if task != nil {
		env = task.ExecutionEnvironment(config, label, extra, false)
	}

	if Has(workspace) {
		claims := operations.RunClaims(validity*60, workspace)
		data, err = operations.AuthorizeClaims(AccountName(), claims)
	}

	if err != nil {
		return err
	}

	if len(data) > 0 {
		endpoint := data["endpoint"]
		for _, key := range rcHosts {
			env = append(env, fmt.Sprintf("%s=%s", key, endpoint))
		}
		token := data["token"]
		for _, key := range rcTokens {
			env = append(env, fmt.Sprintf("%s=%s", key, token))
		}
		env = append(env, fmt.Sprintf("RC_WORKSPACE_ID=%s", workspaceId))
	}

	if jsonform {
		return asJson(env)
	}

	asText(env)
	return nil
}

var variablesCmd = &cobra.Command{
	Use:     "variables <conda.yaml*>",
	Aliases: []string{"vars"},
	Short:   "Export environment specific variables as a JSON structure.",
	Long:    "Export environment specific variables as a JSON structure.",
	Run: func(cmd *cobra.Command, args []string) {
		silent := common.Silent
		common.Silent = true

		defer func() {
			common.Silent = silent
		}()

		ok := conda.MustConda()
		if !ok {
			pretty.Exit(2, "Could not get miniconda installed.")
		}
		err := exportEnvironment(args, robotFile, runTask, environmentFile, workspaceId, validityTime, jsonFlag)
		if err != nil {
			pretty.Exit(1, "Error: Variable exporting failed because: %v", err)
		}
	},
}

func init() {
	envCmd.AddCommand(variablesCmd)

	variablesCmd.Flags().StringVarP(&environmentFile, "environment", "e", "", "Full path to 'env.json' development environment data file. <optional>")
	variablesCmd.Flags().StringVarP(&robotFile, "robot", "r", "robot.yaml", "Full path to 'robot.yaml' configuration file. (Backward compatibility with 'package.yaml')  <optional>")
	variablesCmd.Flags().StringVarP(&runTask, "task", "t", "", "Task to run from configuration file. <optional>")
	variablesCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "Optional workspace id to get authorization tokens for. <optional>")
	variablesCmd.Flags().IntVarP(&validityTime, "minutes", "m", 0, "How many minutes the authorization should be valid for. <optional>")
	variablesCmd.Flags().StringVarP(&accountName, "account", "", "", "Account used for workspace. <optional>")

	variablesCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output in JSON format")
	variablesCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force conda cache update. (only for new environments)")
}
