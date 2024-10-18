package cmd

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/htfs"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/set"
	"github.com/spf13/cobra"
)

var (
	metafileFlag bool
	forceBuild   bool
	exportFile   string
)

func conditionalExpand(filename string) string {
	if !pathlib.IsFile(filename) {
		return filename
	}
	fullpath, err := filepath.Abs(filename)
	if err != nil {
		return filename
	}
	return fullpath
}

func resolveMetafileURL(link string) ([]string, error) {
	origin, err := url.Parse(link)
	refok := err == nil
	raw, err := cloud.ReadFile(link)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, line := range strings.SplitAfter(string(raw), "\n") {
		flat := strings.TrimSpace(line)
		if strings.HasPrefix(flat, "#") || len(flat) == 0 {
			continue
		}
		here, err := url.Parse(flat)
		if refok && err == nil {
			relative := origin.ResolveReference(here)
			result = append(result, relative.String())
		} else {
			result = append(result, flat)
		}
	}
	return result, nil
}

func resolveMetafile(link string) ([]string, error) {
	if !pathlib.IsFile(link) {
		return resolveMetafileURL(link)
	}
	fullpath, err := filepath.Abs(link)
	if err != nil {
		return nil, err
	}
	basedir := filepath.Dir(fullpath)
	raw, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, line := range strings.SplitAfter(string(raw), "\n") {
		flat := strings.TrimSpace(line)
		if strings.HasPrefix(flat, "#") || len(flat) == 0 {
			continue
		}
		result = append(result, filepath.Join(basedir, flat))
	}
	return result, nil
}

func metafileExpansion(links []string, expand bool) []string {
	if !expand {
		return links
	}
	result := []string{}
	for _, metalink := range links {
		links, err := resolveMetafile(conditionalExpand(metalink))
		if err != nil {
			pretty.Warning("Failed to resolve %q metafile, reason: %v", metalink, err)
			continue
		}
		result = append(result, links...)
	}
	return result
}

var holotreePrebuildCmd = &cobra.Command{
	Use:   "prebuild",
	Short: "Prebuild hololib from given set of environment descriptors.",
	Long:  "Prebuild hololib from given set of environment descriptors. Requires shared holotree to be enabled and active.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag() {
			defer common.Stopwatch("Holotree prebuild lasted").Report()
		}

		pretty.Guard(common.SharedHolotree, 1, "Shared holotree must be enabled and in use for prebuild environments to work correctly.")

		configurations := metafileExpansion(args, metafileFlag)
		total, failed := len(configurations), 0
		success := make([]string, 0, total)
		exporting := len(exportFile) > 0
		for at, configfile := range configurations {
			environment, err := conda.ReadPackageCondaYaml(configfile, false)
			if err != nil {
				pretty.Warning("%d/%d: Failed to load %q, reason: %v (ignored)", at+1, total, configfile, err)
				continue
			}
			pretty.Note("%d/%d: Now building config %q", at+1, total, configfile)
			_, _, err = htfs.NewEnvironment(configfile, "", false, forceBuild, operations.PullCatalog)
			if err != nil {
				failed += 1
				pretty.Warning("%d/%d: Holotree recording error: %v", at+1, total, err)
			} else {
				for _, hash := range environment.FingerprintLayers() {
					key := htfs.CatalogName(hash)
					if exporting && !set.Member(success, key) {
						success = append(success, key)
						pretty.Note("Added catalog %q to be exported.", key)
					}
				}
			}
		}
		if exporting && len(success) > 0 {
			holotreeExport(selectCatalogs(success), nil, exportFile)
		}
		pretty.Guard(failed == 0, 2, "%d out of %d environment builds failed! See output above for details.", failed, total)
		pretty.Ok()
	},
}

func init() {
	holotreeCmd.AddCommand(holotreePrebuildCmd)
	holotreePrebuildCmd.Flags().BoolVarP(&metafileFlag, "metafile", "m", false, "Input arguments are actually files containing links/filenames of environment descriptors.")
	holotreePrebuildCmd.Flags().BoolVarP(&forceBuild, "force", "f", false, "Force environment builds, even when blueprint is already present.")
	holotreePrebuildCmd.Flags().StringVarP(&exportFile, "export", "e", "", "Optional filename to export new, successfully build catalogs.")
}
