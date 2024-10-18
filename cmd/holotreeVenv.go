package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/robocorp/rcc/blobs"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/fail"
	"github.com/robocorp/rcc/htfs"
	"github.com/robocorp/rcc/journal"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/shell"

	"github.com/spf13/cobra"
)

func deleteByExactIdentity(exact string) {
	_, roots := htfs.LoadCatalogs()
	for _, label := range roots.FindEnvironments([]string{exact}) {
		common.Log("Removing %v", label)
		err := roots.RemoveHolotreeSpace(label)
		pretty.Guard(err == nil, 4, "Error: %v", err)
	}
}

var holotreeVenvCmd = &cobra.Command{
	Use:   "venv conda.yaml+",
	Short: "Create user managed virtual python environment inside automation folder.",
	Long:  "Create user managed virtual python environment inside automation folder.",
	Args:  cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		defer journal.BuildEventStats("venv")
		if common.DebugFlag() {
			defer common.Stopwatch("Holotree venv command lasted").Report()
		}

		// following settings are forced in venv environments
		common.UnmanagedSpace = true
		common.ExternallyManaged = true
		common.ControllerType = "venv"

		where, err := os.Getwd()
		pretty.Guard(err == nil, 1, "Error: %v", err)
		location := filepath.Join(where, "venv")

		previous := pathlib.IsDir(location)
		if holotreeForce && previous {
			pretty.Note("Trying to remove existing venv at %q ...", location)
			err := pathlib.TryRemoveAll("venv", location)
			pretty.Guard(err == nil, 2, "Error: %v", err)
		}

		pretty.Guard(!pathlib.Exists(location), 3, "Name %q aready exists! Remove it, or use force.", location)

		if holotreeForce {
			identity := htfs.ControllerSpaceName([]byte(common.ControllerIdentity()), []byte(common.HolotreeSpace))
			deleteByExactIdentity(identity)
		}

		env := holotreeExpandEnvironment(args, "", "", "", 0, holotreeForce, common.DeveloperFlag)
		envPath := pathlib.EnvironmentPath(env)
		python, ok := envPath.Which("python", conda.FileExtensions)
		if !ok {
			python, ok = envPath.Which("python3", conda.FileExtensions)
		}
		pretty.Guard(ok, 5, "For some reason, could not find python executable in environment paths. Report a bug. PATH: %q", envPath)
		pretty.Note("Trying to make new venv at %q using %q ...", location, python)
		task := shell.New(env, ".", python, "-m", "venv", "--system-site-packages", location)
		code, err := task.Execute(false)
		pretty.Guard(err == nil, 6, "Error: %v", err)
		pretty.Guard(code == 0, 7, "Exit code %d from venv creation.", code)

		target := listActivationScripts(location)
		if len(target) > 0 {
			blob, err := blobs.Asset("assets/depxtraction.py")
			fail.Fast(err)
			location := filepath.Join(target, "depxtraction.py")
			fail.Fast(os.WriteFile(location, blob, 0o755))
			fmt.Printf("Experimental dependency extraction tool is available at %q.\nTry it after pip installing things into your venv.\n", location)
		}

		pretty.Ok()
	},
}

func listActivationScripts(root string) string {
	pretty.Highlight("New venv is located at %s. Following scripts seem to be available:", root)
	base := filepath.Dir(root)
	pathcandidate := ""
	filepath.Walk(root, func(path string, entry fs.FileInfo, err error) error {
		if entry.Mode().IsRegular() && strings.HasPrefix(strings.ToLower(entry.Name()), "activ") {
			short, err := filepath.Rel(base, path)
			if err == nil {
				pretty.Highlight(" - %s", short)
			}
			pathcandidate = filepath.Dir(short)
		}
		return nil
	})
	return pathcandidate
}

func init() {
	rootCmd.AddCommand(holotreeVenvCmd)
	holotreeCmd.AddCommand(holotreeVenvCmd)

	holotreeVenvCmd.Flags().StringVarP(&common.HolotreeSpace, "space", "s", "user", "Client specific name to identify this environment.")
	holotreeVenvCmd.Flags().BoolVarP(&holotreeForce, "force", "f", false, "Force environment creation by deleting unmanaged space. Dangerous, do not use unless you understand what it means.")
	holotreeVenvCmd.Flags().BoolVarP(&common.DevDependencies, "devdeps", "", false, "Include dev-dependencies from the `package.yaml` file in the environment (only valid when dealing with a `package.yaml` file).")
}
