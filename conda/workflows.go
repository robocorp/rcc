package conda

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/journal"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/settings"
	"github.com/robocorp/rcc/shell"
)

const (
	venvTemplate = `home = %s
include-system-site-packages = true
version = %s
`
)

func metafile(folder string) string {
	return common.ExpandPath(folder + ".meta")
}

func livePrepare(liveFolder string, command ...string) (*shell.Task, error) {
	commandName := command[0]
	task, ok := HolotreePath(liveFolder).Which(commandName, FileExtensions)
	if !ok {
		return nil, fmt.Errorf("Cannot find command: %v", commandName)
	}
	common.Debug("Using %v as command %v.", task, commandName)
	command[0] = task
	environment := CondaExecutionEnvironment(liveFolder, nil, true)
	return shell.New(environment, ".", command...), nil
}

func LiveCapture(liveFolder string, command ...string) (string, int, error) {
	task, err := livePrepare(liveFolder, command...)
	if err != nil {
		return "", 9999, err
	}
	return task.CaptureOutput()
}

func LiveExecution(sink io.Writer, liveFolder string, command ...string) (int, error) {
	fmt.Fprintf(sink, "Command %q at %q:\n", command, liveFolder)
	task, err := livePrepare(liveFolder, command...)
	if err != nil {
		return 0, err
	}
	return task.Tracked(sink, false)
}

type InstallObserver map[string]bool

func (it InstallObserver) Write(content []byte) (int, error) {
	text := strings.ToLower(string(content))
	if strings.Contains(text, "safetyerror:") {
		it["safetyerror"] = true
	}
	if strings.Contains(text, "pkgs") {
		it["pkgs"] = true
	}
	if strings.Contains(text, "appears to be corrupted") {
		it["corrupted"] = true
	}
	return len(content), nil
}

func (it InstallObserver) HasFailures(targetFolder string) bool {
	if it["safetyerror"] && it["corrupted"] && len(it) > 2 {
		cloud.BackgroundMetric(common.ControllerIdentity(), "rcc.env.creation.failure", common.Version)
		renameRemove(targetFolder)
		location := filepath.Join(common.RobocorpHome(), "pkgs")
		common.Log("%sWARNING! Conda environment is unstable, see above error.%s", pretty.Red, pretty.Reset)
		common.Log("%sWARNING! To fix it, try to remove directory: %v%s", pretty.Red, location, pretty.Reset)
		return true
	}
	return false
}

func newLive(yaml, condaYaml, requirementsText, key string, force, freshInstall bool, postInstall []string) (bool, error) {
	if !MustMicromamba() {
		return false, fmt.Errorf("Could not get micromamba installed.")
	}
	targetFolder := common.StageFolder
	common.Debug("===  pre cleanup phase ===")
	common.Timeline("pre cleanup phase.")
	err := renameRemove(targetFolder)
	if err != nil {
		return false, err
	}
	common.Debug("===  first try phase ===")
	common.Timeline("first try.")
	success, fatal := newLiveInternal(yaml, condaYaml, requirementsText, key, force, freshInstall, postInstall)
	if !success && !force && !fatal {
		journal.CurrentBuildEvent().Rebuild()
		cloud.BackgroundMetric(common.ControllerIdentity(), "rcc.env.creation.retry", common.Version)
		common.Debug("===  second try phase ===")
		common.Timeline("second try.")
		common.ForceDebug()
		common.Log("Retry! First try failed ... now retrying with debug and force options!")
		err = renameRemove(targetFolder)
		if err != nil {
			return false, err
		}
		success, _ = newLiveInternal(yaml, condaYaml, requirementsText, key, true, freshInstall, postInstall)
	}
	if success {
		journal.CurrentBuildEvent().Successful()
	}
	return success, nil
}

func newLiveInternal(yaml, condaYaml, requirementsText, key string, force, freshInstall bool, postInstall []string) (bool, bool) {
	targetFolder := common.StageFolder
	planfile := fmt.Sprintf("%s.plan", targetFolder)
	planSink, err := os.OpenFile(planfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return false, false
	}
	defer func() {
		planSink.Close()
		content, err := os.ReadFile(planfile)
		if err == nil {
			common.Log("%s", string(content))
		}
		os.Remove(planfile)
	}()

	planalyzer := NewPlanAnalyzer(true)
	defer planalyzer.Close()

	planWriter := io.MultiWriter(planSink, planalyzer)
	fmt.Fprintf(planWriter, "---  installation plan %q %s [force: %v, fresh: %v| rcc %s]  ---\n\n", key, time.Now().Format(time.RFC3339), force, freshInstall, common.Version)
	stopwatch := common.Stopwatch("installation plan")
	fmt.Fprintf(planWriter, "---  plan blueprint @%ss  ---\n\n", stopwatch)
	fmt.Fprintf(planWriter, "%s\n", yaml)

	common.Debug("Setting up new conda environment using %v to folder %v", condaYaml, targetFolder)
	ttl := "57600"
	if force {
		ttl = "0"
	}
	common.Progress(5, "Running micromamba phase. (micromamba v%s)", MicromambaVersion())
	mambaCommand := common.NewCommander(BinMicromamba(), "create", "--always-copy", "--no-env", "--safety-checks", "enabled", "--extra-safety-checks", "--retry-clean-cache", "--strict-channel-priority", "--repodata-ttl", ttl, "-y", "-f", condaYaml, "-p", targetFolder)
	mambaCommand.Option("--channel-alias", settings.Global.CondaURL())
	mambaCommand.ConditionalFlag(common.VerboseEnvironmentBuilding(), "--verbose")
	mambaCommand.ConditionalFlag(!settings.Global.HasMicroMambaRc(), "--no-rc")
	mambaCommand.ConditionalFlag(settings.Global.HasMicroMambaRc(), "--rc-file", common.MicroMambaRcFile())
	observer := make(InstallObserver)
	common.Debug("===  micromamba create phase ===")
	fmt.Fprintf(planWriter, "\n---  micromamba plan @%ss  ---\n\n", stopwatch)
	tee := io.MultiWriter(observer, planWriter)
	code, err := shell.New(CondaEnvironment(), ".", mambaCommand.CLI()...).Tracked(tee, false)
	if err != nil || code != 0 {
		cloud.BackgroundMetric(common.ControllerIdentity(), "rcc.env.fatal.micromamba", fmt.Sprintf("%d_%x", code, code))
		common.Timeline("micromamba fail.")
		common.Fatal(fmt.Sprintf("Micromamba [%d/%x]", code, code), err)
		return false, false
	}
	journal.CurrentBuildEvent().MicromambaComplete()
	common.Timeline("micromamba done.")
	if observer.HasFailures(targetFolder) {
		return false, true
	}
	fmt.Fprintf(planWriter, "\n---  pip plan @%ss  ---\n\n", stopwatch)
	python, pyok := FindPython(targetFolder)
	if !pyok {
		fmt.Fprintf(planWriter, "Note: no python in target folder: %s\n", targetFolder)
	}
	pipUsed, pipCache, wheelCache := false, common.PipCache(), common.WheelCache()
	size, ok := pathlib.Size(requirementsText)
	if !ok || size == 0 {
		common.Progress(6, "Skipping pip install phase -- no pip dependencies.")
	} else {
		if !pyok {
			cloud.BackgroundMetric(common.ControllerIdentity(), "rcc.env.fatal.pip", fmt.Sprintf("%d_%x", 9999, 9999))
			common.Timeline("pip fail. no python found.")
			common.Fatal("pip fail. no python found.", errors.New("No python found, but required!"))
			return false, false
		}
		common.Progress(6, "Running pip install phase. (pip v%s)", PipVersion(python))
		common.Debug("Updating new environment at %v with pip requirements from %v (size: %v)", targetFolder, requirementsText, size)
		pipCommand := common.NewCommander(python, "-m", "pip", "install", "--isolated", "--no-color", "--disable-pip-version-check", "--prefer-binary", "--cache-dir", pipCache, "--find-links", wheelCache, "--requirement", requirementsText)
		pipCommand.Option("--index-url", settings.Global.PypiURL())
		pipCommand.Option("--trusted-host", settings.Global.PypiTrustedHost())
		pipCommand.ConditionalFlag(common.VerboseEnvironmentBuilding(), "--verbose")
		common.Debug("===  pip install phase ===")
		code, err = LiveExecution(planWriter, targetFolder, pipCommand.CLI()...)
		planSink.Sync()
		if err != nil || code != 0 {
			cloud.BackgroundMetric(common.ControllerIdentity(), "rcc.env.fatal.pip", fmt.Sprintf("%d_%x", code, code))
			common.Timeline("pip fail.")
			common.Fatal(fmt.Sprintf("Pip [%d/%x]", code, code), err)
			return false, false
		}
		journal.CurrentBuildEvent().PipComplete()
		common.Timeline("pip done.")
		pipUsed = true
	}
	fmt.Fprintf(planWriter, "\n---  post install plan @%ss  ---\n\n", stopwatch)
	if postInstall != nil && len(postInstall) > 0 {
		common.Progress(7, "Post install scripts phase started.")
		common.Debug("===  post install phase ===")
		for _, script := range postInstall {
			scriptCommand, err := shell.Split(script)
			if err != nil {
				common.Fatal("post-install", err)
				common.Log("%sScript '%s' parsing failure: %v%s", pretty.Red, script, err, pretty.Reset)
				return false, false
			}
			common.Debug("Running post install script '%s' ...", script)
			_, err = LiveExecution(planWriter, targetFolder, scriptCommand...)
			planSink.Sync()
			if err != nil {
				common.Fatal("post-install", err)
				common.Log("%sScript '%s' failure: %v%s", pretty.Red, script, err, pretty.Reset)
				return false, false
			}
		}
		journal.CurrentBuildEvent().PostInstallComplete()
	} else {
		common.Progress(7, "Post install scripts phase skipped -- no scripts.")
	}
	common.Progress(8, "Activate environment started phase.")
	common.Debug("===  activate phase ===")
	fmt.Fprintf(planWriter, "\n---  activation plan @%ss  ---\n\n", stopwatch)
	err = Activate(planWriter, targetFolder)
	if err != nil {
		common.Log("%sActivation failure: %v%s", pretty.Yellow, err, pretty.Reset)
	}
	for _, line := range LoadActivationEnvironment(targetFolder) {
		fmt.Fprintf(planWriter, "%s\n", line)
	}
	err = goldenMaster(targetFolder, pipUsed)
	if err != nil {
		common.Log("%sGolden EE failure: %v%s", pretty.Yellow, err, pretty.Reset)
	}
	fmt.Fprintf(planWriter, "\n---  pip check plan @%ss  ---\n\n", stopwatch)
	if common.StrictFlag && pipUsed {
		common.Progress(9, "Running pip check phase.")
		pipCommand := common.NewCommander(python, "-m", "pip", "check", "--no-color")
		pipCommand.ConditionalFlag(common.VerboseEnvironmentBuilding(), "--verbose")
		common.Debug("===  pip check phase ===")
		code, err = LiveExecution(planWriter, targetFolder, pipCommand.CLI()...)
		planSink.Sync()
		if err != nil || code != 0 {
			cloud.BackgroundMetric(common.ControllerIdentity(), "rcc.env.fatal.pipcheck", fmt.Sprintf("%d_%x", code, code))
			common.Timeline("pip check fail.")
			common.Fatal(fmt.Sprintf("Pip check [%d/%x]", code, code), err)
			return false, false
		}
		common.Timeline("pip check done.")
	} else {
		common.Progress(9, "Pip check skipped.")
	}
	fmt.Fprintf(planWriter, "\n---  installation plan complete @%ss  ---\n\n", stopwatch)
	planSink.Sync()
	planSink.Close()
	common.Progress(10, "Update installation plan.")
	finalplan := filepath.Join(targetFolder, "rcc_plan.log")
	os.Rename(planfile, finalplan)
	common.Debug("===  finalize phase ===")

	markerFile := filepath.Join(targetFolder, "identity.yaml")
	err = os.WriteFile(markerFile, []byte(yaml), 0o644)
	if err != nil {
		return false, false
	}

	_, ok = HolotreePath(targetFolder).Which("python", FileExtensions)
	if ok {
		venvContent := fmt.Sprintf(venvTemplate, targetFolder, pythonVersionAt(targetFolder))
		venvFile := filepath.Join(targetFolder, "pyvenv.cfg")
		err = os.WriteFile(venvFile, []byte(venvContent), 0o644)
		if err != nil {
			return false, false
		}
	}

	return true, false
}

func pythonVersionAt(targetFolder string) string {
	versionText, code, _ := LiveCapture(targetFolder, "python", "--version")
	if code != 0 {
		return "?.?.?"
	}
	return strings.Replace(versionText, "Python ", "", 1)
}

func temporaryConfig(condaYaml, requirementsText string, save bool, filenames ...string) (string, string, *Environment, error) {
	var left, right *Environment
	var err error

	for _, filename := range filenames {
		left = right
		right, err = ReadCondaYaml(filename)
		if err != nil {
			return "", "", nil, err
		}
		if left == nil {
			continue
		}
		right, err = left.Merge(right)
		if err != nil {
			return "", "", nil, err
		}
	}
	yaml, err := right.AsYaml()
	if err != nil {
		return "", "", nil, err
	}
	hash := common.ShortDigest(yaml)
	if !save {
		return hash, yaml, right, nil
	}
	common.Log("FINAL union conda environment descriptor:\n---\n%v---", yaml)
	err = right.SaveAsRequirements(requirementsText)
	if err != nil {
		return "", "", nil, err
	}
	pure := right.AsPureConda()
	err = pure.SaveAs(condaYaml)
	return hash, yaml, right, err
}

func LegacyEnvironment(force bool, configurations ...string) error {
	cloud.BackgroundMetric(common.ControllerIdentity(), "rcc.env.create.start", common.Version)

	lockfile := common.RobocorpLock()
	completed := pathlib.LockWaitMessage("Serialized environment creation [robocorp lock]")
	locker, err := pathlib.Locker(lockfile, 30000)
	completed()
	if err != nil {
		common.Log("Could not get lock on live environment. Quitting!")
		return err
	}
	defer locker.Release()

	freshInstall := true

	condaYaml := filepath.Join(os.TempDir(), fmt.Sprintf("conda_%x.yaml", common.When))
	requirementsText := filepath.Join(os.TempDir(), fmt.Sprintf("require_%x.txt", common.When))
	common.Debug("Using temporary conda.yaml file: %v and requirement.txt file: %v", condaYaml, requirementsText)
	key, yaml, finalEnv, err := temporaryConfig(condaYaml, requirementsText, true, configurations...)
	if err != nil {
		return err
	}
	defer os.Remove(condaYaml)
	defer os.Remove(requirementsText)

	success, err := newLive(yaml, condaYaml, requirementsText, key, force, freshInstall, finalEnv.PostInstall)
	if err != nil {
		return err
	}
	if success {
		return nil
	}

	return errors.New("Could not create environment.")
}

func renameRemove(location string) error {
	if !pathlib.IsDir(location) {
		common.Trace("Location %q is not directory, not removed.", location)
		return nil
	}
	randomLocation := fmt.Sprintf("%s.%08X", location, rand.Uint32())
	common.Debug("Rename/remove %q using %q as random name.", location, randomLocation)
	err := os.Rename(location, randomLocation)
	if err != nil {
		common.Log("Rename %q -> %q failed as: %v!", location, randomLocation, err)
		return err
	}
	common.Trace("Rename %q -> %q was successful!", location, randomLocation)
	err = os.RemoveAll(randomLocation)
	if err != nil {
		common.Log("Removal of %q failed as: %v!", randomLocation, err)
		return err
	}
	common.Trace("Removal of %q was successful!", randomLocation)
	meta := metafile(location)
	if pathlib.IsFile(meta) {
		err = os.Remove(meta)
		common.Trace("Removal of %q result was %v.", meta, err)
		return err
	}
	common.Trace("Metafile %q was not file.", meta)
	return nil
}
