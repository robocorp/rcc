package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/htfs"
	"github.com/robocorp/rcc/journal"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/robot"
	"github.com/spf13/cobra"
)

const (
	newEnvironment = `environment creation`
)

var (
	holotreeBlueprint []byte
	holotreeForce     bool
	holotreeJson      bool
)

func asSimpleMap(line string) map[string]string {
	parts := strings.SplitN(strings.TrimSpace(line), "=", 2)
	if len(parts) != 2 {
		return nil
	}
	if len(parts[0]) == 0 {
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
	common.Stdout("%s\n", content)
	return nil
}

func asExportedText(items []string) {
	prefix := "export"
	if conda.IsWindows() {
		prefix = "SET"
	}
	for _, line := range items {
		common.Stdout("%s %s\n", prefix, line)
	}
}

func holotreeExpandEnvironment(userFiles []string, packfile, environment, workspace string, validity int, force bool, devDependencies bool) []string {
	var extra []string
	var data operations.Token
	common.TimelineBegin("environment expansion start")
	defer common.TimelineEnd()

	config, holotreeBlueprint, err := htfs.ComposeFinalBlueprint(userFiles, packfile, devDependencies)
	pretty.Guard(err == nil, 5, "%s", err)

	condafile := filepath.Join(common.ProductTemp(), common.BlueprintHash(holotreeBlueprint))
	err = pathlib.WriteFile(condafile, holotreeBlueprint, 0o644)
	pretty.Guard(err == nil, 6, "%s", err)

	holozip := ""
	if config != nil {
		holozip = config.Holozip()
	}

	// i.e.: the conda file is now already created in the temp folder, so, there's no need to use the devDependencies flag
	// anymore.
	path, _, err := htfs.NewEnvironment(condafile, holozip, true, force, operations.PullCatalog)
	if !common.WarrantyVoided() {
		pretty.RccPointOfView(newEnvironment, err)
	}
	pretty.Guard(err == nil, 6, "%s", err)

	if Has(environment) {
		common.Timeline("load robot environment")
		developmentEnvironment, err := robot.LoadEnvironmentSetup(environment)
		if err == nil {
			extra = developmentEnvironment.AsEnvironment()
		}
	}

	common.Timeline("load robot environment")
	var env []string
	if config != nil {
		env = config.RobotExecutionEnvironment(path, extra, false)
	} else {
		env = conda.CondaExecutionEnvironment(path, extra, false)
	}

	if Has(workspace) {
		common.Timeline("get run robot claims")
		period := &operations.TokenPeriod{
			ValidityTime: validityTime,
			GracePeriod:  gracePeriod,
		}
		period.EnforceGracePeriod()
		claims := operations.RunRobotClaims(period.RequestSeconds(), workspace)
		data, err = operations.AuthorizeClaims(AccountName(), claims, period)
		pretty.Guard(err == nil, 9, "Failed to get cloud data, reason: %v", err)
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

	return removeUnwanted(env)
}

func removeUnwanted(variables []string) []string {
	result := make([]string, 0, len(variables))
	for _, line := range variables {
		switch {
		case strings.HasPrefix(line, "PS1="):
			continue
		default:
			result = append(result, line)
		}
	}

	return result
}

var holotreeVariablesCmd = &cobra.Command{
	Use:     "variables conda.yaml+",
	Aliases: []string{"vars"},
	Short:   "Do holotree operations.",
	Long:    "Do holotree operations.",
	Run: func(cmd *cobra.Command, args []string) {
		defer journal.BuildEventStats("variables")
		if common.DebugFlag() {
			defer common.Stopwatch("Holotree variables command lasted").Report()
		}

		env := holotreeExpandEnvironment(args, robotFile, environmentFile, workspaceId, validityTime, holotreeForce, common.DevDependencies)
		if holotreeJson {
			asJson(env)
		} else {
			asExportedText(env)
		}
	},
}

func init() {
	holotreeCmd.AddCommand(holotreeVariablesCmd)
	holotreeVariablesCmd.Flags().StringVarP(&environmentFile, "environment", "e", "", "Full path to 'env.json' development environment data file. <optional>")
	holotreeVariablesCmd.Flags().StringVarP(&robotFile, "robot", "r", "robot.yaml", "Full path to 'robot.yaml' configuration file. <optional>")
	holotreeVariablesCmd.Flags().StringVarP(&workspaceId, "workspace", "w", "", "Optional workspace id to get authorization tokens for. <optional>")
	holotreeVariablesCmd.Flags().IntVarP(&validityTime, "minutes", "m", 15, "How many minutes the authorization should be valid for (minimum 15 minutes).")
	holotreeVariablesCmd.Flags().IntVarP(&gracePeriod, "graceperiod", "", 5, "What is grace period buffer in minutes on top of validity minutes (minimum 5 minutes).")
	holotreeVariablesCmd.Flags().StringVarP(&accountName, "account", "a", "", "Account used for workspace. <optional>")

	holotreeVariablesCmd.Flags().StringVarP(&common.HolotreeSpace, "space", "s", "user", "Client specific name to identify this environment.")
	holotreeVariablesCmd.Flags().BoolVarP(&holotreeForce, "force", "f", false, "Force environment creation with refresh.")
	holotreeVariablesCmd.Flags().BoolVarP(&holotreeJson, "json", "j", false, "Show environment as JSON.")
	holotreeVariablesCmd.Flags().BoolVarP(&common.DevDependencies, "devdeps", "", false, "Include dev-dependencies from the `package.yaml` file in the environment (only valid when dealing with a `package.yaml` file).")
}
