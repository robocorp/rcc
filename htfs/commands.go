package htfs

import (
	"os"
	"path/filepath"
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

func NewEnvironment(condafile, holozip string, restore, force bool) (label string, scorecard common.Scorecard, err error) {
	defer fail.Around(&err)

	journal.CurrentBuildEvent().StartNow(force)

	if settings.Global.NoBuild() {
		pretty.Note("'no-build' setting is active. Only cached, prebuild, or imported environments are allowed!")
	}

	haszip := len(holozip) > 0
	if haszip {
		common.Debug("New zipped environment from %q!", holozip)
	}

	path := ""
	defer func() {
		common.Progress(13, "Fresh holotree done [with %d workers].", anywork.Scale())
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
		common.Progress(1, "Fresh [shared mode] holotree environment %v.", xviper.TrackingIdentity())
	} else {
		common.Progress(1, "Fresh [private mode] holotree environment %v.", xviper.TrackingIdentity())
	}

	completed := pathlib.LockWaitMessage("Serialized environment creation [holotree lock]")
	locker, err := pathlib.Locker(common.HolotreeLock(), 30000)
	completed()
	fail.On(err != nil, "Could not get lock for holotree. Quiting.")
	defer locker.Release()

	_, holotreeBlueprint, err := ComposeFinalBlueprint([]string{condafile}, "")
	fail.On(err != nil, "%s", err)
	common.EnvironmentHash = BlueprintHash(holotreeBlueprint)
	common.Progress(2, "Holotree blueprint is %q [%s].", common.EnvironmentHash, common.Platform())
	journal.CurrentBuildEvent().Blueprint(common.EnvironmentHash)

	tree, err := New()
	fail.On(err != nil, "%s", err)

	if !haszip && !tree.HasBlueprint(holotreeBlueprint) && common.Liveonly {
		tree = Virtual()
		common.Timeline("downgraded to virtual holotree library")
	}
	if common.UnmanagedSpace {
		tree = Unmanaged(tree)
	}
	err = tree.ValidateBlueprint(holotreeBlueprint)
	fail.On(err != nil, "%s", err)
	scorecard = common.NewScorecard()
	var library Library
	if haszip {
		library, err = ZipLibrary(holozip)
		fail.On(err != nil, "Failed to load %q -> %s", holozip, err)
		common.Timeline("downgraded to holotree zip library")
	} else {
		scorecard.Start()
		err = RecordEnvironment(tree, holotreeBlueprint, force, scorecard)
		fail.On(err != nil, "%s", err)
		library = tree
	}

	if restore {
		common.Progress(12, "Restore space from library [with %d workers].", anywork.Scale())
		path, err = library.Restore(holotreeBlueprint, []byte(common.ControllerIdentity()), []byte(common.HolotreeSpace))
		fail.On(err != nil, "Failed to restore blueprint %q, reason: %v", string(holotreeBlueprint), err)
		journal.CurrentBuildEvent().RestoreComplete()
	} else {
		common.Progress(12, "Restoring space skipped.")
	}

	return path, scorecard, nil
}

func CleanupHolotreeStage(tree MutableLibrary) error {
	common.Timeline("holotree stage removal start")
	defer common.Timeline("holotree stage removal done")
	return TryRemoveAll("stage", tree.Stage())
}

func RecordEnvironment(tree MutableLibrary, blueprint []byte, force bool, scorecard common.Scorecard) (err error) {
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

	if force || !exists {
		common.Progress(3, "Cleanup holotree stage for fresh install.")
		fail.On(settings.Global.NoBuild(), "Building new holotree environment is blocked by settings, and could not be found from hololib cache!")
		err = CleanupHolotreeStage(tree)
		fail.On(err != nil, "Failed to clean stage, reason %v.", err)
		journal.CurrentBuildEvent().PrepareComplete()

		err = os.MkdirAll(tree.Stage(), 0o755)
		fail.On(err != nil, "Failed to create stage, reason %v.", err)

		common.Progress(4, "Build environment into holotree stage.")
		identityfile := filepath.Join(tree.Stage(), "identity.yaml")
		err = os.WriteFile(identityfile, blueprint, 0o644)
		fail.On(err != nil, "Failed to save %q, reason %w.", identityfile, err)
		err = conda.LegacyEnvironment(force, identityfile)
		fail.On(err != nil, "Failed to create environment, reason %w.", err)

		scorecard.Midpoint()

		common.Progress(11, "Record holotree stage to hololib [with %d workers].", anywork.Scale())
		err = tree.Record(blueprint)
		fail.On(err != nil, "Failed to record blueprint %q, reason: %w", string(blueprint), err)
		journal.CurrentBuildEvent().RecordComplete()
	}

	return nil
}

func FindEnvironment(fragment string) []string {
	result := make([]string, 0, 10)
	for directory, _ := range Spacemap() {
		name := filepath.Base(directory)
		if strings.Contains(name, fragment) {
			result = append(result, name)
		}
	}
	return result
}

func InstallationPlan(hash string) (string, bool) {
	finalplan := filepath.Join(common.HolotreeLocation(), hash, "rcc_plan.log")
	return finalplan, pathlib.IsFile(finalplan)
}

func RemoveHolotreeSpace(label string) (err error) {
	defer fail.Around(&err)

	for directory, metafile := range Spacemap() {
		name := filepath.Base(directory)
		if name != label {
			continue
		}
		TryRemove("metafile", metafile)
		TryRemove("lockfile", directory+".lck")
		err = TryRemoveAll("space", directory)
		fail.On(err != nil, "Problem removing %q, reason: %s.", directory, err)
	}
	return nil
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

func ComposeFinalBlueprint(userFiles []string, packfile string) (config robot.Robot, blueprint []byte, err error) {
	defer fail.Around(&err)

	var left, right *conda.Environment

	config, filenames := RobotBlueprints(userFiles, packfile)

	for _, filename := range filenames {
		left = right
		right, err = conda.ReadCondaYaml(filename)
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
	return config, []byte(strings.TrimSpace(content)), nil
}
