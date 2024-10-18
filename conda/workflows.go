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
	SkipNoLayers         SkipLayer = iota
	SkipMicromambaLayer  SkipLayer = iota
	SkipPipLayer         SkipLayer = iota
	SkipPostinstallLayer SkipLayer = iota
	SkipError            SkipLayer = iota
)

const (
	micromambaInstall  = `micromamba install`
	pipInstall         = `pip install`
	uvInstall          = `uv install`
	postInstallScripts = `post-install script execution`
)

type (
	pipTool func(string, string, string, fmt.Stringer, io.Writer) (bool, bool, bool, string)

	SkipLayer uint8
	Recorder  interface {
		Record([]byte) error
	}
	PlanWriter struct {
		filename string
		blob     []byte
	}
)

func NewPlanWriter(filename string) *PlanWriter {
	return &PlanWriter{
		filename: filename,
		blob:     make([]byte, 0, 50000),
	}
}

func (it *PlanWriter) AsText() string {
	return string(it.blob)
}

func (it *PlanWriter) Write(blob []byte) (int, error) {
	it.blob = append(it.blob, blob...)
	return len(blob), nil
}

func (it *PlanWriter) Save() error {
	return os.WriteFile(it.filename, it.blob, 0o644)
}

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
		cloud.InternalBackgroundMetric(common.ControllerIdentity(), "rcc.env.creation.failure", common.Version)
		renameRemove(targetFolder)
		location := filepath.Join(common.Product.Home(), "pkgs")
		common.Log("%sWARNING! Conda environment is unstable, see above error.%s", pretty.Red, pretty.Reset)
		common.Log("%sWARNING! To fix it, try to remove directory: %v%s", pretty.Red, location, pretty.Reset)
		return true
	}
	return false
}

func newLive(yaml, condaYaml, requirementsText, key string, force, freshInstall bool, skip SkipLayer, finalEnv *Environment, recorder Recorder) (bool, error) {
	if !MustMicromamba() {
		return false, fmt.Errorf("Could not get micromamba installed.")
	}
	targetFolder := common.StageFolder
	if skip == SkipNoLayers {
		common.Debug("===  pre cleanup phase ===")
		common.Timeline("pre cleanup phase.")
		err := renameRemove(targetFolder)
		if err != nil {
			return false, err
		}
	}
	common.Debug("===  first try phase ===")
	common.Timeline("first try.")
	success, fatal := newLiveInternal(yaml, condaYaml, requirementsText, key, force, freshInstall, skip, finalEnv, recorder)
	if !success && !force && !fatal && !common.NoRetryBuild {
		journal.CurrentBuildEvent().Rebuild()
		cloud.InternalBackgroundMetric(common.ControllerIdentity(), "rcc.env.creation.retry", common.Version)
		common.Debug("===  second try phase ===")
		common.Timeline("second try.")
		common.ForceDebug()
		common.Log("Retry! First try failed ... now retrying with debug and force options!")
		err := renameRemove(targetFolder)
		if err != nil {
			return false, err
		}
		success, _ = newLiveInternal(yaml, condaYaml, requirementsText, key, true, freshInstall, SkipNoLayers, finalEnv, recorder)
	}
	if success {
		journal.CurrentBuildEvent().Successful()
	}
	return success, nil
}

func assertStageFolder(location string) {
	base := filepath.Base(location)
	holotree := strings.HasPrefix(base, "h") && strings.HasSuffix(base, "t")
	virtual := strings.HasPrefix(base, "v") && strings.HasSuffix(base, "h")
	if !(holotree || virtual) {
		panic(fmt.Sprintf("FATAL: incorrect stage %q for environment building!", location))
	}
}

func micromambaLayer(fingerprint, condaYaml, targetFolder string, stopwatch fmt.Stringer, planWriter io.Writer, force bool) (bool, bool) {
	assertStageFolder(targetFolder)
	common.TimelineBegin("Layer: micromamba [%s]", fingerprint)
	defer common.TimelineEnd()

	common.Debug("Setting up new conda environment using %v to folder %v", condaYaml, targetFolder)
	ttl := "57600"
	if force {
		ttl = "0"
	}
	pretty.Progress(7, "Running micromamba phase. (micromamba v%s) [layer: %s]", MicromambaVersion(), fingerprint)
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
		cloud.InternalBackgroundMetric(common.ControllerIdentity(), "rcc.env.fatal.micromamba", fmt.Sprintf("%d_%x", code, code))
		common.Timeline("micromamba fail.")
		common.Fatal(fmt.Sprintf("Micromamba [%d/%x]", code, code), err)
		pretty.RccPointOfView(micromambaInstall, err)
		return false, false
	}
	journal.CurrentBuildEvent().MicromambaComplete()
	common.Timeline("micromamba done.")
	if observer.HasFailures(targetFolder) {
		return false, true
	}
	return true, false
}

func uvLayer(fingerprint, requirementsText, targetFolder string, stopwatch fmt.Stringer, planWriter io.Writer) (bool, bool, bool, string) {
	assertStageFolder(targetFolder)
	common.TimelineBegin("Layer: uv [%s]", fingerprint)
	defer common.TimelineEnd()

	pipUsed := false
	fmt.Fprintf(planWriter, "\n---  uv plan @%ss  ---\n\n", stopwatch)
	uv, uvok := FindUv(targetFolder)
	if !uvok {
		fmt.Fprintf(planWriter, "Note: no uv in target folder: %s\n", targetFolder)
		return false, false, pipUsed, ""
	}
	python, pyok := FindPython(targetFolder)
	if !pyok {
		fmt.Fprintf(planWriter, "Note: no python in target folder: %s\n", targetFolder)
	}
	uvCache, wheelCache := common.UvCache(), common.WheelCache()
	size, ok := pathlib.Size(requirementsText)
	if !ok || size == 0 {
		pretty.Progress(8, "Skipping pip install phase -- no pip dependencies.")
	} else {
		pretty.Progress(8, "Running uv install phase. (uv v%s) [layer: %s]", UvVersion(uv), fingerprint)
		common.Debug("Updating new environment at %v with uv requirements from %v (size: %v)", targetFolder, requirementsText, size)
		uvCommand := common.NewCommander(uv, "pip", "install", "--link-mode", "copy", "--color", "never", "--cache-dir", uvCache, "--find-links", wheelCache, "--requirement", requirementsText)
		uvCommand.Option("--index-url", settings.Global.PypiURL())
		// no "--trusted-host" on uv pip install
		// uvCommand.Option("--trusted-host", settings.Global.PypiTrustedHost())
		uvCommand.ConditionalFlag(common.VerboseEnvironmentBuilding(), "--verbose")
		common.Debug("===  uv install phase ===")
		code, err := LiveExecution(planWriter, targetFolder, uvCommand.CLI()...)
		if err != nil || code != 0 {
			cloud.InternalBackgroundMetric(common.ControllerIdentity(), "rcc.env.fatal.uv", fmt.Sprintf("%d_%x", code, code))
			common.Timeline("uv fail.")
			common.Fatal(fmt.Sprintf("uv [%d/%x]", code, code), err)
			pretty.RccPointOfView(uvInstall, err)
			return false, false, pipUsed, ""
		}
		journal.CurrentBuildEvent().PipComplete()
		common.Timeline("uv done.")
		pipUsed = true
	}
	return true, false, pipUsed, python
}

func pipLayer(fingerprint, requirementsText, targetFolder string, stopwatch fmt.Stringer, planWriter io.Writer) (bool, bool, bool, string) {
	assertStageFolder(targetFolder)
	common.TimelineBegin("Layer: pip [%s]", fingerprint)
	defer common.TimelineEnd()

	pipUsed := false
	fmt.Fprintf(planWriter, "\n---  pip plan @%ss  ---\n\n", stopwatch)
	python, pyok := FindPython(targetFolder)
	if !pyok {
		fmt.Fprintf(planWriter, "Note: no python in target folder: %s\n", targetFolder)
	}
	pipCache, wheelCache := common.PipCache(), common.WheelCache()
	size, ok := pathlib.Size(requirementsText)
	if !ok || size == 0 {
		pretty.Progress(8, "Skipping pip install phase -- no pip dependencies.")
	} else {
		if !pyok {
			cloud.InternalBackgroundMetric(common.ControllerIdentity(), "rcc.env.fatal.pip", fmt.Sprintf("%d_%x", 9999, 9999))
			common.Timeline("pip fail. no python found.")
			common.Fatal("pip fail. no python found.", errors.New("No python found, but required!"))
			return false, false, pipUsed, ""
		}
		pretty.Progress(8, "Running pip install phase. (pip v%s) [layer: %s]", PipVersion(python), fingerprint)
		common.Debug("Updating new environment at %v with pip requirements from %v (size: %v)", targetFolder, requirementsText, size)
		pipCommand := common.NewCommander(python, "-m", "pip", "install", "--isolated", "--no-color", "--disable-pip-version-check", "--prefer-binary", "--cache-dir", pipCache, "--find-links", wheelCache, "--requirement", requirementsText)
		pipCommand.Option("--index-url", settings.Global.PypiURL())
		pipCommand.Option("--trusted-host", settings.Global.PypiTrustedHost())
		pipCommand.ConditionalFlag(common.VerboseEnvironmentBuilding(), "--verbose")
		common.Debug("===  pip install phase ===")
		code, err := LiveExecution(planWriter, targetFolder, pipCommand.CLI()...)
		if err != nil || code != 0 {
			cloud.InternalBackgroundMetric(common.ControllerIdentity(), "rcc.env.fatal.pip", fmt.Sprintf("%d_%x", code, code))
			common.Timeline("pip fail.")
			common.Fatal(fmt.Sprintf("Pip [%d/%x]", code, code), err)
			pretty.RccPointOfView(pipInstall, err)
			return false, false, pipUsed, ""
		}
		journal.CurrentBuildEvent().PipComplete()
		common.Timeline("pip done.")
		pipUsed = true
	}
	return true, false, pipUsed, python
}

func postInstallLayer(fingerprint string, postInstall []string, targetFolder string, stopwatch fmt.Stringer, planWriter io.Writer) (bool, bool) {
	assertStageFolder(targetFolder)
	common.TimelineBegin("Layer: post install scripts [%s]", fingerprint)
	defer common.TimelineEnd()

	fmt.Fprintf(planWriter, "\n---  post install plan @%ss  ---\n\n", stopwatch)
	if postInstall != nil && len(postInstall) > 0 {
		pretty.Progress(9, "Post install scripts phase started. [layer: %s]", fingerprint)
		common.Debug("===  post install phase ===")
		for _, script := range postInstall {
			scriptCommand, err := shell.Split(script)
			if err != nil {
				common.Fatal("post-install", err)
				common.Log("%sScript '%s' parsing failure: %v%s", pretty.Red, script, err, pretty.Reset)
				pretty.RccPointOfView(postInstallScripts, err)
				return false, false
			}
			common.Debug("Running post install script '%s' ...", script)
			_, err = LiveExecution(planWriter, targetFolder, scriptCommand...)
			if err != nil {
				common.Fatal("post-install", err)
				common.Log("%sScript '%s' failure: %v%s", pretty.Red, script, err, pretty.Reset)
				pretty.RccPointOfView(postInstallScripts, err)
				return false, false
			}
		}
		journal.CurrentBuildEvent().PostInstallComplete()
	} else {
		pretty.Progress(9, "Post install scripts phase skipped -- no scripts.")
	}
	return true, false
}

func holotreeLayers(condaYaml, requirementsText string, finalEnv *Environment, targetFolder string, stopwatch fmt.Stringer, planWriter io.Writer, theplan *PlanWriter, force bool, skip SkipLayer, recorder Recorder) (bool, bool, bool, string) {
	assertStageFolder(targetFolder)
	common.TimelineBegin("Holotree layers at %q", targetFolder)
	defer common.TimelineEnd()

	pipNeeded := len(requirementsText) > 0
	postInstall := len(finalEnv.PostInstall) > 0

	var pypiSelector pipTool = pipLayer

	hasUv := finalEnv.HasCondaDependency("uv")
	if hasUv {
		pypiSelector = uvLayer
	}

	layers := finalEnv.AsLayers()
	fingerprints := finalEnv.FingerprintLayers()

	var success, fatal, pipUsed bool
	var python string

	if skip < SkipMicromambaLayer {
		success, fatal = micromambaLayer(fingerprints[0], condaYaml, targetFolder, stopwatch, planWriter, force)
		if !success {
			return success, fatal, false, ""
		}
		if pipNeeded || postInstall {
			fmt.Fprintf(theplan, "\n---  micromamba layer complete [on layered holotree]  ---\n\n")
			common.Error("saving rcc_plan.log", theplan.Save())
			common.Error("saving golden master", goldenMaster(targetFolder, false))
			recorder.Record([]byte(layers[0]))
		}
	} else {
		pretty.Progress(7, "Skipping micromamba phase, layer exists.")
		fmt.Fprintf(planWriter, "\n---  micromamba plan skipped, layer exists ---\n\n")
	}
	if skip < SkipPipLayer {
		success, fatal, pipUsed, python = pypiSelector(fingerprints[1], requirementsText, targetFolder, stopwatch, planWriter)
		if !success {
			return success, fatal, pipUsed, python
		}
		if pipUsed && postInstall {
			fmt.Fprintf(theplan, "\n---  pip layer complete [on layered holotree]  ---\n\n")
			common.Error("saving rcc_plan.log", theplan.Save())
			common.Error("saving golden master", goldenMaster(targetFolder, true))
			recorder.Record([]byte(layers[1]))
		}
	} else {
		pretty.Progress(8, "Skipping pip phase, layer exists.")
		fmt.Fprintf(planWriter, "\n---  pip plan skiped, layer exists  ---\n\n")
	}
	if skip < SkipPostinstallLayer {
		success, fatal = postInstallLayer(fingerprints[2], finalEnv.PostInstall, targetFolder, stopwatch, planWriter)
		if !success {
			return success, fatal, pipUsed, python
		}
	} else {
		pretty.Progress(9, "Skipping post install scripts phase, layer exists.")
		fmt.Fprintf(planWriter, "\n---  post install plan skipped, layer exists  ---\n\n")
	}
	return true, false, pipUsed, python
}

func newLiveInternal(yaml, condaYaml, requirementsText, key string, force, freshInstall bool, skip SkipLayer, finalEnv *Environment, recorder Recorder) (bool, bool) {
	targetFolder := common.StageFolder
	theplan := NewPlanWriter(filepath.Join(targetFolder, "rcc_plan.log"))
	failure := true
	defer func() {
		if failure {
			common.Log("%s", theplan.AsText())
		}
	}()

	planalyzer := NewPlanAnalyzer(true)
	defer planalyzer.Close()

	planWriter := io.MultiWriter(theplan, planalyzer)
	fmt.Fprintf(planWriter, "---  installation plan %q %s [force: %v, fresh: %v| rcc %s]  ---\n\n", key, time.Now().Format(time.RFC3339), force, freshInstall, common.Version)
	stopwatch := common.Stopwatch("installation plan")
	fmt.Fprintf(planWriter, "---  plan blueprint @%ss  ---\n\n", stopwatch)
	fmt.Fprintf(planWriter, "%s\n", yaml)

	success, fatal, pipUsed, python := holotreeLayers(condaYaml, requirementsText, finalEnv, targetFolder, stopwatch, planWriter, theplan, force, skip, recorder)
	if !success {
		return success, fatal
	}

	pretty.Progress(10, "Activate environment started phase.")
	common.Debug("===  activate phase ===")
	fmt.Fprintf(planWriter, "\n---  activation plan @%ss  ---\n\n", stopwatch)
	err := Activate(planWriter, targetFolder)
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
		pretty.Progress(11, "Running pip check phase.")
		pipCommand := common.NewCommander(python, "-m", "pip", "check", "--no-color")
		pipCommand.ConditionalFlag(common.VerboseEnvironmentBuilding(), "--verbose")
		common.Debug("===  pip check phase ===")
		code, err := LiveExecution(planWriter, targetFolder, pipCommand.CLI()...)
		if err != nil || code != 0 {
			cloud.InternalBackgroundMetric(common.ControllerIdentity(), "rcc.env.fatal.pipcheck", fmt.Sprintf("%d_%x", code, code))
			common.Timeline("pip check fail.")
			common.Fatal(fmt.Sprintf("Pip check [%d/%x]", code, code), err)
			return false, false
		}
		common.Timeline("pip check done.")
	} else {
		pretty.Progress(11, "Pip check skipped.")
	}
	fmt.Fprintf(planWriter, "\n---  installation plan complete @%ss  ---\n\n", stopwatch)
	pretty.Progress(12, "Update installation plan.")
	common.Error("saving rcc_plan.log", theplan.Save())
	common.Debug("===  finalize phase ===")

	failure = false

	return true, false
}

func LogUnifiedEnvironment(content []byte) {
	environment, err := CondaYamlFrom(content)
	if err != nil {
		return
	}
	yaml, err := environment.AsYaml()
	if err != nil {
		return
	}
	common.Log("FINAL unified conda environment descriptor:\n---\n%v---", yaml)
}

func finalUnifiedEnvironment(filename string) (string, *Environment, error) {
	right, err := ReadPackageCondaYaml(filename, false)
	if err != nil {
		return "", nil, err
	}
	yaml, err := right.AsYaml()
	if err != nil {
		return "", nil, err
	}
	return yaml, right, nil
}

func temporaryConfig(condaYaml, requirementsText, filename string) (string, string, *Environment, error) {
	yaml, right, err := finalUnifiedEnvironment(filename)
	if err != nil {
		return "", "", nil, err
	}
	hash := common.ShortDigest(yaml)
	err = right.SaveAsRequirements(requirementsText)
	if err != nil {
		return "", "", nil, err
	}
	pure := right.AsPureConda()
	err = pure.SaveAs(condaYaml)
	return hash, yaml, right, err
}

func LegacyEnvironment(recorder Recorder, force bool, skip SkipLayer, configuration string) error {
	cloud.InternalBackgroundMetric(common.ControllerIdentity(), "rcc.env.create.start", common.Version)

	lockfile := common.ProductLock()
	completed := pathlib.LockWaitMessage(lockfile, "Serialized environment creation [robocorp lock]")
	locker, err := pathlib.Locker(lockfile, 30000, false)
	completed()
	if err != nil {
		common.Log("Could not get lock on live environment. Quitting!")
		return err
	}
	defer locker.Release()

	freshInstall := true

	condaYaml := filepath.Join(pathlib.TempDir(), fmt.Sprintf("conda_%x.yaml", common.When))
	requirementsText := filepath.Join(pathlib.TempDir(), fmt.Sprintf("require_%x.txt", common.When))
	common.Debug("Using temporary conda.yaml file: %v and requirement.txt file: %v", condaYaml, requirementsText)
	key, yaml, finalEnv, err := temporaryConfig(condaYaml, requirementsText, configuration)
	if err != nil {
		return err
	}
	defer os.Remove(condaYaml)
	defer os.Remove(requirementsText)

	success, err := newLive(yaml, condaYaml, requirementsText, key, force, freshInstall, skip, finalEnv, recorder)
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
