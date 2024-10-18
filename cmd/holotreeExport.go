package cmd

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/htfs"
	"github.com/robocorp/rcc/pretty"
	"github.com/spf13/cobra"
)

var (
	holozip     string
	exportRobot string
)

func holotreeExport(catalogs, known []string, archive string) {
	common.Debug("Ignoring content from catalogs:")
	for _, catalog := range known {
		common.Debug("- %s", catalog)
	}

	common.Debug("Exporting catalogs:")
	for _, catalog := range catalogs {
		common.Debug("- %s", catalog)
	}

	tree, err := htfs.New()
	pretty.Guard(err == nil, 2, "%s", err)

	err = tree.Export(catalogs, known, archive)
	pretty.Guard(err == nil, 3, "%s", err)
}

func listCatalogs(jsonForm bool) {
	if jsonForm {
		nice, err := json.MarshalIndent(htfs.CatalogNames(), "", "  ")
		pretty.Guard(err == nil, 2, "%s", err)
		common.Stdout("%s\n", nice)
	} else {
		common.Log("Selectable catalogs (you can use substrings):")
		for _, catalog := range htfs.CatalogNames() {
			common.Log("- %s", catalog)
		}
	}
}

func selectCatalogs(filters []string) []string {
	result := make([]string, 0, len(filters))
	for _, catalog := range htfs.CatalogNames() {
		for _, filter := range filters {
			if strings.Contains(catalog, filter) {
				result = append(result, catalog)
				break
			}
		}
	}
	sort.Strings(result)
	return result
}

var holotreeExportCmd = &cobra.Command{
	Use:   "export catalog+",
	Short: "Export existing holotree catalog and library parts.",
	Long:  "Export existing holotree catalog and library parts.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag() {
			defer common.Stopwatch("Holotree export command lasted").Report()
		}
		if len(exportRobot) > 0 {
			devDependencies := false
			_, holotreeBlueprint, err := htfs.ComposeFinalBlueprint(nil, exportRobot, devDependencies)
			pretty.Guard(err == nil, 1, "Blueprint calculation failed: %v", err)
			hash := common.BlueprintHash(holotreeBlueprint)
			args = append(args, htfs.CatalogName(hash))
		}
		if len(args) == 0 {
			listCatalogs(jsonFlag)
		} else {
			holotreeExport(selectCatalogs(args), nil, holozip)
		}
		pretty.Ok()
	},
}

func init() {
	holotreeCmd.AddCommand(holotreeExportCmd)
	holotreeExportCmd.Flags().StringVarP(&holozip, "zipfile", "z", "hololib.zip", "Name of zipfile to export.")
	holotreeExportCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output in JSON format")
	holotreeExportCmd.Flags().StringVarP(&exportRobot, "robot", "r", "", "Full path to 'robot.yaml' configuration file to export as catalog. <optional>")
}
