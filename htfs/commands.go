package htfs

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/robocorp/rcc/anywork"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/fail"
	"github.com/robocorp/rcc/journal"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/robot"
	"github.com/robocorp/rcc/settings"
	"github.com/robocorp/rcc/xviper"
)

type CatalogPuller func(string, string, bool) error

func NewEnvironment(condafile, holozip string, restore, force bool, puller CatalogPuller) (label string, scorecard common.Scorecard, err error) {
	defer fail.Around(&err)

	who, _ := user.Current()
	host, _ := os.Hostname()
	pretty.Progress(0, "Context: %q <%v@%v> [%v/%v].", who.Name, who.Username, host, common.Platform(), settings.OperatingSystem())

	if common.WarrantyVoided() {
		tree, err := New()
		fail.On(err != nil, "%s", err)

		path := tree.WarrantyVoidedDir([]byte(common.ControllerIdentity()), []byte(common.HolotreeSpace))
		return path, common.NewScorecard(), err
	}

	journal.CurrentBuildEvent().StartNow(force)

	if settings.Global.NoBuild() {
		pretty.Note("'no-build' setting is active. Only cached, prebuild, or imported environments are allowed!")
	}

	haszip := len(holozip) > 0
	if haszip {
		common.Debug("New zipped environment from %q!", holozip)
	}

	path, externally := "", ""
	defer func() {
		if err != nil {
			pretty.Regression(15, "Holotree restoration failure, see above [with %d workers on %d CPUs].", anywork.Scale(), runtime.NumCPU())
		} else {
			pretty.Progress(15, "Fresh %sholotree done [with %d workers on %d CPUs].", externally, anywork.Scale(), runtime.NumCPU())
			if len(path) > 0 {
				usefile := fmt.Sprintf("%s.use", path)
				pathlib.AppendFile(usefile, []byte{'.'})
			}
		}
		if haszip {
			pretty.Note("There is hololib.zip present at: %q", holozip)
		}
		if len(path) > 0 {
			dependencies := conda.LoadWantedDependencies(conda.GoldenMasterFilename(path))
			dependencies.WarnVulnerability(
				"https://robocorp.com/docs/faq/openssl-cve-2022-11-01",
				"HIGH",
				"openssl",
				"3.0.0", "3.0.1", "3.0.2", "3.0.3", "3.0.4", "3.0.5", "3.0.6")
		}
	}()
	if common.SharedHolotree {
		pretty.Progress(1, "Fresh [shared mode] holotree environment %v. (parent/pid: %d/%d)", xviper.TrackingIdentity(), os.Getppid(), os.Getpid())
	} else {
		pretty.Progress(1, "Fresh [private mode] holotree environment %v. (parent/pid: %d/%d)", xviper.TrackingIdentity(), os.Getppid(), os.Getpid())
	}

	lockfile := common.HolotreeLock()
	completed := pathlib.LockWaitMessage(lockfile, "Serialized environment creation [holotree lock]")
	locker, err := pathlib.Locker(lockfile, 30000, common.SharedHolotree)
	completed()
	fail.On(err != nil, "Could not get lock for holotree. Quiting.")
	defer locker.Release()

	_, holotreeBlueprint, err := ComposeFinalBlueprint([]string{condafile}, "", false)
	fail.Fast(err)

	common.EnvironmentHash, common.FreshlyBuildEnvironment = common.BlueprintHash(holotreeBlueprint), false
	pretty.Progress(2, "Holotree blueprint is %q [%s with %d workers on %d CPUs from %q].", common.EnvironmentHash, common.Platform(), anywork.Scale(), runtime.NumCPU(), filepath.Base(condafile))
	journal.CurrentBuildEvent().Blueprint(common.EnvironmentHash)

	tree, err := New()
	fail.Fast(err)

	if !haszip && !tree.HasBlueprint(holotreeBlueprint) && common.Liveonly {
		tree = Virtual()
		common.Timeline("downgraded to virtual holotree library")
	}
	if common.UnmanagedSpace {
		tree = Unmanaged(tree)
	}
	fail.Fast(tree.ValidateBlueprint(holotreeBlueprint))
	scorecard = common.NewScorecard()
	var library Library
	if haszip {
		library, err = ZipLibrary(holozip)
		fail.On(err != nil, "Failed to load %q -> %s", holozip, err)
		common.Timeline("downgraded to holotree zip library")
	} else {
		scorecard.Start()
		fail.Fast(RecordEnvironment(tree, holotreeBlueprint, force, scorecard, puller))
		library = tree
	}

	if restore {
		pretty.Progress(14, "Restore space from library [with %d workers on %d CPUs; with compression: %v].", anywork.Scale(), runtime.NumCPU(), Compress())
		path, err = library.Restore(holotreeBlueprint, []byte(common.ControllerIdentity()), []byte(common.HolotreeSpace))
		fail.On(err != nil, "Failed to restore blueprint %q, reason: %v", string(holotreeBlueprint), err)
		journal.CurrentBuildEvent().RestoreComplete()
	} else {
		pretty.Progress(14, "Restoring space skipped.")
	}

	externally, err = conda.ApplyExternallyManaged(path)
	fail.Fast(err)
	fail.Fast(CleanupHolotreeStage(tree))

	return path, scorecard, nil
}

func CleanupHolotreeStage(tree MutableLibrary) error {
	common.TimelineBegin("holotree stage removal")
	defer common.TimelineEnd()
	location := tree.Stage()
	common.Timeline("- removing %q", location)
	return pathlib.TryRemoveAll("stage", location)
}

func RecordEnvironment(tree MutableLibrary, blueprint []byte, force bool, scorecard common.Scorecard, puller CatalogPuller) (err error) {
	defer fail.Around(&err)

	// following must be setup here
	common.StageFolder = tree.Stage()
	backup := common.Liveonly
	common.Liveonly = true
	defer func() {
		common.Liveonly = backup
	}()

	common.Debug("Holotree stage is %q.", tree.Stage())
	exists := tree.HasBlueprint(blueprint)
	common.Debug("Has blueprint environment: %v", exists)

	conda.LogUnifiedEnvironment(blueprint)

	if force || !exists {
		common.FreshlyBuildEnvironment = true
		remoteOrigin := common.RccRemoteOrigin()
		if len(remoteOrigin) > 0 {
			pretty.Progress(3, "Fill hololib from RCC_REMOTE_ORIGIN.")
			hash := common.BlueprintHash(blueprint)
			catalog := CatalogName(hash)
			err = puller(remoteOrigin, catalog, false)
			if err != nil {
				pretty.Warning("Failed to pull %q from %q, reason: %v", catalog, remoteOrigin, err)
			} else {
				return nil
			}
			exists = tree.HasBlueprint(blueprint)
		} else {
			pretty.Progress(3, "Fill hololib from RCC_REMOTE_ORIGIN skipped. RCC_REMOTE_ORIGIN was not defined.")
		}
		pretty.Progress(4, "Cleanup holotree stage for fresh install.")
		fail.On(settings.Global.NoBuild(), "Building new holotree environment is blocked by settings, and could not be found from hololib cache!")
		err = CleanupHolotreeStage(tree)
		fail.On(err != nil, "Failed to clean stage, reason %v.", err)
		journal.CurrentBuildEvent().PrepareComplete()

		err = os.MkdirAll(tree.Stage(), 0o755)
		fail.On(err != nil, "Failed to create stage, reason %v.", err)

		pretty.Progress(5, "Build environment into holotree stage %q.", tree.Stage())
		identityfile := filepath.Join(tree.Stage(), "identity.yaml")
		err = os.WriteFile(identityfile, blueprint, 0o644)
		fail.On(err != nil, "Failed to save %q, reason %w.", identityfile, err)

		skip := conda.SkipNoLayers
		if !force {
			pretty.Progress(6, "Restore partial environment into holotree stage %q.", tree.Stage())
			skip = RestoreLayersTo(tree, identityfile, tree.Stage())
		} else {
			pretty.Progress(6, "Restore partial environment skipped. Force used.")
		}

		err = os.WriteFile(identityfile, blueprint, 0o644)
		fail.On(err != nil, "Failed to save %q, reason %w.", identityfile, err)

		err = conda.LegacyEnvironment(tree, force, skip, identityfile)
		fail.On(err != nil, "Failed to create environment, reason %w.", err)

		scorecard.Midpoint()

		pretty.Progress(13, "Record holotree stage to hololib [with %d workers on %d CPUs].", anywork.Scale(), runtime.NumCPU())
		err = tree.Record(blueprint)
		fail.On(err != nil, "Failed to record blueprint %q, reason: %w", string(blueprint), err)
		journal.CurrentBuildEvent().RecordComplete()
	}

	return nil
}

func RestoreLayersTo(tree MutableLibrary, identityfile string, targetDir string) conda.SkipLayer {
	config, err := conda.ReadPackageCondaYaml(identityfile, false)
	if err != nil {
		return conda.SkipNoLayers
	}

	layers := config.AsLayers()
	mambaLayer := []byte(layers[0])
	pipLayer := []byte(layers[1])
	base := filepath.Base(targetDir)
	if tree.HasBlueprint(pipLayer) {
		_, err = tree.RestoreTo(pipLayer, base, common.ControllerIdentity(), common.HolotreeSpace, true)
		if err == nil {
			return conda.SkipPipLayer
		}
	}
	if tree.HasBlueprint(mambaLayer) {
		_, err = tree.RestoreTo(mambaLayer, base, common.ControllerIdentity(), common.HolotreeSpace, true)
		if err == nil {
			return conda.SkipMicromambaLayer
		}
	}
	return conda.SkipNoLayers
}

func RobotBlueprints(userBlueprints []string, packfile string) (robot.Robot, []string) {
	var err error
	var config robot.Robot

	blueprints := make([]string, 0, len(userBlueprints)+2)

	if len(packfile) > 0 {
		config, err = robot.LoadRobotYaml(packfile, false)
		if err == nil {
			blueprints = append(blueprints, config.CondaConfigFile())
		}
	}

	return config, append(blueprints, userBlueprints...)
}

func ComposeFinalBlueprint(userFiles []string, packfile string, devDependencies bool) (config robot.Robot, blueprint []byte, err error) {
	defer fail.Around(&err)

	var left, right *conda.Environment

	config, filenames := RobotBlueprints(userFiles, packfile)

	for _, filename := range filenames {
		left = right
		right, err = conda.ReadPackageCondaYaml(filename, devDependencies)
		fail.On(err != nil, "Failure: %v", err)
		if left == nil {
			continue
		}
		right, err = left.Merge(right)
		fail.On(err != nil, "Failure: %v", err)
	}
	fail.On(right == nil, "Missing environment specification(s).")
	content, err := right.AsYaml()
	fail.On(err != nil, "YAML error: %v", err)
	blueprint = []byte(strings.TrimSpace(content))
	if !right.IsCacheable() {
		fingerprint := common.BlueprintHash(blueprint)
		pretty.Warning("Holotree blueprint %q is not publicly cacheable. Use `rcc robot diagnostics` to find out more.", fingerprint)
	}
	return config, blueprint, nil
}
