package robot

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/shell"

	"gopkg.in/yaml.v2"
)

var (
	GoosPattern   = regexp.MustCompile("(?i:(windows|darwin|linux))")
	GoarchPattern = regexp.MustCompile("(?i:(amd64|arm64))")
)

type Robot interface {
	IgnoreFiles() []string
	AvailableTasks() []string
	DefaultTask() Task
	TaskByName(string) Task
	UsesConda() bool
	CondaConfigFile() string
	PreRunScripts() []string
	RootDirectory() string
	HasHolozip() bool
	Holozip() string
	Validate() (bool, error)
	Diagnostics(*common.DiagnosticStatus, bool)
	DependenciesFile() (string, bool)

	WorkingDirectory() string
	ArtifactDirectory() string
	FreezeFilename() string
	Paths() pathlib.PathParts
	PythonPaths() pathlib.PathParts
	SearchPath(location string) pathlib.PathParts
	RobotExecutionEnvironment(location string, inject []string, full bool) []string
}

type Task interface {
	Commandline() []string
}

type robot struct {
	Tasks        map[string]*task `yaml:"tasks"`
	Devtasks     map[string]*task `yaml:"devTasks"`
	Conda        string           `yaml:"condaConfigFile,omitempty"`
	PreRun       []string         `yaml:"preRunScripts,omitempty"`
	Environments []string         `yaml:"environmentConfigs,omitempty"`
	Ignored      []string         `yaml:"ignoreFiles"`
	Artifacts    string           `yaml:"artifactsDir"`
	Path         []string         `yaml:"PATH"`
	Pythonpath   []string         `yaml:"PYTHONPATH"`
	Root         string
}

type task struct {
	Task    string   `yaml:"robotTaskName,omitempty"`
	Shell   string   `yaml:"shell,omitempty"`
	Command []string `yaml:"command,omitempty"`
	robot   *robot
}

func (it *robot) taskMap(note bool) map[string]*task {
	if common.DeveloperFlag {
		if note {
			pretty.Note("Operating in developer mode. Using 'devTasks:' instead of 'tasks:'.")
		}
		return it.Devtasks
	} else {
		return it.Tasks
	}
}

func (it *robot) relink() {
	for _, task := range it.Tasks {
		if task != nil {
			task.robot = it
		}
	}
	for _, task := range it.Devtasks {
		if task != nil {
			task.robot = it
		}
	}
}

func (it *robot) diagnoseTasks(diagnose common.Diagnoser) {
	if it.Tasks == nil {
		diagnose.Fail("", "Missing 'tasks:' from robot.yaml.")
		return
	}
	ok := true
	if len(it.Tasks) == 0 {
		diagnose.Fail("", "There must be at least one task defined in 'tasks:' section in robot.yaml.")
		ok = false
	} else {
		diagnose.Ok("Tasks are defined in robot.yaml")
	}
	for name, task := range it.Tasks {
		count := 0
		if len(task.Task) > 0 {
			count += 1
		}
		if len(task.Shell) > 0 {
			count += 1
		}
		if task.Command != nil && len(task.Command) > 0 {
			count += 1
		}
		if count != 1 {
			diagnose.Fail("", "In robot.yaml, task '%s' needs exactly one of robotTaskName/shell/command definition!", name)
			ok = false
		}
	}
	if ok {
		diagnose.Ok("Each task has exactly one definition.")
	}
}

func (it *robot) diagnoseVariousPaths(diagnose common.Diagnoser) {
	ok := true
	for _, path := range it.Path {
		if filepath.IsAbs(path) {
			diagnose.Fail("", "PATH entry %q seems to be absolute, which makes robot machine dependent.", path)
			ok = false
		}
	}
	if ok {
		diagnose.Ok("PATH settings in robot.yaml are ok.")
	}
	ok = true
	for _, path := range it.Pythonpath {
		if filepath.IsAbs(path) {
			diagnose.Fail("", "PYTHONPATH entry %q seems to be absolute, which makes robot machine dependent.", path)
			ok = false
		}
	}
	if ok {
		diagnose.Ok("PYTHONPATH settings in robot.yaml are ok.")
	}
	ok = true
	if it.Ignored == nil || len(it.Ignored) == 0 {
		diagnose.Warning("", "No ignoreFiles defined, so everything ends up inside robot.zip file.")
		ok = false
	} else {
		for at, path := range it.Ignored {
			if len(strings.TrimSpace(path)) == 0 {
				diagnose.Fail("", "there is empty entry in ignoreFiles at position %d", at+1)
				ok = false
				continue
			}
			if filepath.IsAbs(path) {
				diagnose.Fail("", "ignoreFiles entry %q seems to be absolute, which makes robot machine dependent.", path)
				ok = false
			}
		}
		for _, path := range it.IgnoreFiles() {
			if !pathlib.IsFile(path) {
				diagnose.Fail("", "ignoreFiles entry %q is not a file.", path)
				ok = false
			}
		}
	}
	if ok {
		diagnose.Ok("ignoreFiles settings in robot.yaml are ok.")
	}
}

func (it *robot) Diagnostics(target *common.DiagnosticStatus, production bool) {
	diagnose := target.Diagnose("Robot")
	it.diagnoseTasks(diagnose)
	it.diagnoseVariousPaths(diagnose)
	if it.Artifacts == "" {
		diagnose.Fail("", "In robot.yaml, 'artifactsDir:' is required!")
	} else {
		if filepath.IsAbs(it.Artifacts) {
			diagnose.Fail("", "artifactDir %q seems to be absolute, which makes robot machine dependent.", it.Artifacts)
		} else {
			diagnose.Ok("Artifacts directory defined in robot.yaml")
		}
	}
	if it.Conda == "" {
		diagnose.Ok("In robot.yaml, 'condaConfigFile:' is missing. So this is shell robot.")
	} else {
		if filepath.IsAbs(it.Conda) {
			diagnose.Fail("", "condaConfigFile %q seems to be absolute, which makes robot machine dependent.", it.Artifacts)
		} else {
			diagnose.Ok("In robot.yaml, 'condaConfigFile:' is present. So this is python robot.")
			condaEnv, err := conda.ReadCondaYaml(it.CondaConfigFile())
			if err != nil {
				diagnose.Fail("", "From robot.yaml, loading conda.yaml failed with: %v", err)
			} else {
				condaEnv.Diagnostics(target, production)
			}
		}
	}
	target.Details["robot-use-conda"] = fmt.Sprintf("%v", it.UsesConda())
	target.Details["robot-conda-file"] = it.CondaConfigFile()
	target.Details["hololib.zip"] = it.Holozip()
	target.Details["robot-root-directory"] = it.RootDirectory()
	target.Details["robot-working-directory"] = it.WorkingDirectory()
	target.Details["robot-artifact-directory"] = it.ArtifactDirectory()
	target.Details["robot-paths"] = strings.Join(it.Paths(), ", ")
	target.Details["robot-python-paths"] = strings.Join(it.PythonPaths(), ", ")
	dependencies, ok := it.DependenciesFile()
	if !ok {
		dependencies = "missing"
	} else {
		if it.VerifyCondaDependencies() {
			diagnose.Ok("Dependencies in conda.yaml and dependencies.yaml match.")
		}
	}
	target.Details["robot-dependencies-yaml"] = dependencies
}

func (it *robot) Validate() (bool, error) {
	if it.Tasks == nil {
		return false, errors.New("In robot.yaml, 'tasks:' is required!")
	}
	if len(it.Tasks) == 0 {
		return false, errors.New("In robot.yaml, 'tasks:' must have at least one task defined!")
	}
	if it.Artifacts == "" {
		return false, errors.New("In robot.yaml, 'artifactsDir:' is required!")
	}
	for name, task := range it.Tasks {
		count := 0
		if len(task.Task) > 0 {
			count += 1
		}
		if len(task.Shell) > 0 {
			count += 1
		}
		if task.Command != nil && len(task.Command) > 0 {
			count += 1
		}
		if count != 1 {
			return false, fmt.Errorf("In robot.yaml, task '%s' needs exactly one of robotTaskName/shell/command definition!", name)
		}
	}
	return true, nil
}

func (it *robot) DependenciesFile() (string, bool) {
	filename := filepath.Join(it.Root, "dependencies.yaml")
	return filename, pathlib.IsFile(filename)
}

func (it *robot) VerifyCondaDependencies() bool {
	wanted, ok := it.DependenciesFile()
	if !ok {
		return true
	}
	dependencies := conda.LoadWantedDependencies(wanted)
	if len(dependencies) == 0 {
		return true
	}
	condaEnv, err := conda.ReadCondaYaml(it.CondaConfigFile())
	if err != nil {
		return true
	}
	ideal, ok := condaEnv.FromDependencies(dependencies)
	if !ok {
		body, err := ideal.AsYaml()
		if err == nil {
			fmt.Println("IDEAL:", body)
		}
	}
	return ok
}

func (it *robot) RootDirectory() string {
	return it.Root
}

func (it *robot) HasHolozip() bool {
	return len(it.Holozip()) > 0
}

func (it *robot) Holozip() string {
	zippath := filepath.Join(it.Root, "hololib.zip")
	if pathlib.IsFile(zippath) {
		return zippath
	}
	return ""
}

func (it *robot) IgnoreFiles() []string {
	if it.Ignored == nil {
		return []string{}
	}
	result := make([]string, 0, len(it.Ignored))
	for at, entry := range it.Ignored {
		if len(strings.TrimSpace(entry)) == 0 {
			pretty.Warning("Ignore file entry at position %d is empty string!", at+1)
			continue
		}
		fullpath := filepath.Join(it.Root, entry)
		if !pathlib.IsFile(fullpath) {
			pretty.Warning("Ignore file %q is not a file!", fullpath)
		}
		result = append(result, fullpath)
	}
	return result
}

func (it *robot) AvailableTasks() []string {
	tasks := it.taskMap(false)
	result := make([]string, 0, len(tasks))
	for name, _ := range tasks {
		result = append(result, fmt.Sprintf("%q", name))
	}
	sort.Strings(result)
	return result
}

func (it *robot) DefaultTask() Task {
	tasks := it.taskMap(true)
	if len(tasks) != 1 {
		return nil
	}
	var result *task
	for _, value := range tasks {
		result = value
		break
	}
	return result
}

func (it *robot) TaskByName(name string) Task {
	if len(name) == 0 {
		return it.DefaultTask()
	}
	tasks := it.taskMap(true)
	key := strings.Trim(name, "\t\r\n\"' ")
	found, ok := tasks[key]
	if ok {
		return found
	}
	caseless := strings.ToLower(key)
	for name, value := range tasks {
		if caseless == strings.ToLower(strings.TrimSpace(name)) {
			return value
		}
	}
	return nil
}

func (it *robot) UsesConda() bool {
	return len(it.Conda) > 0 || len(it.availableEnvironmentConfigurations(common.Platform())) > 0
}

func (it *robot) CondaConfigFile() string {
	available := it.availableEnvironmentConfigurations(common.Platform())
	if len(available) > 0 {
		return available[0]
	}
	return filepath.Join(it.Root, it.Conda)
}

func (it *robot) PreRunScripts() []string {
	return it.PreRun
}

func (it *robot) WorkingDirectory() string {
	return it.Root
}

func freezeFileBasename() string {
	return fmt.Sprintf("environment_%s_freeze.yaml", common.Platform())
}

func submatch(pattern *regexp.Regexp, expected, text string) bool {
	match := pattern.FindStringSubmatch(text)
	return match == nil || len(match) == 0 || match[0] == expected
}

func PlatformAcceptableFile(architecture, operatingSystem, filename string) bool {
	return submatch(GoarchPattern, architecture, filename) && submatch(GoosPattern, operatingSystem, filename)
}

func (it *robot) availableEnvironmentConfigurations(marker string) []string {
	result := make([]string, 0, len(it.Environments))
	common.Trace("Available environment configurations:")
	for _, part := range it.Environments {
		underscored := strings.Count(part, "_") > 2
		freezed := strings.Contains(strings.ToLower(part), "freeze")
		marked := strings.Contains(part, marker)
		if (underscored || freezed) && !marked {
			continue
		}
		if !PlatformAcceptableFile(runtime.GOARCH, runtime.GOOS, part) {
			continue
		}
		fullpath := filepath.Join(it.Root, part)
		if !pathlib.IsFile(fullpath) {
			continue
		}
		common.Trace("- %s", fullpath)
		result = append(result, fullpath)
	}
	if len(result) == 0 {
		common.Trace("- nothing")
	}
	return result
}

func (it *robot) FreezeFilename() string {
	return filepath.Join(it.ArtifactDirectory(), freezeFileBasename())
}

func (it *robot) ArtifactDirectory() string {
	return filepath.Join(it.Root, it.Artifacts)
}

func pathBuilder(root string, tails []string) pathlib.PathParts {
	result := make([]string, 0, len(tails))
	for _, part := range tails {
		if filepath.IsAbs(part) && pathlib.IsDir(part) {
			result = append(result, part)
			continue
		}
		fullpath := filepath.Join(root, part)
		realpath, err := filepath.Abs(fullpath)
		if err == nil {
			result = append(result, realpath)
		}
	}
	return pathlib.PathFrom(result...)
}

func (it *robot) Paths() pathlib.PathParts {
	if it == nil {
		return pathlib.PathFrom()
	}
	return pathBuilder(it.Root, it.Path)
}

func (it *robot) PythonPaths() pathlib.PathParts {
	if it == nil {
		return pathlib.PathFrom()
	}
	return pathBuilder(it.Root, it.Pythonpath)
}

func (it *robot) SearchPath(location string) pathlib.PathParts {
	return conda.FindPath(location).Prepend(it.Paths()...)
}

func (it *robot) RobotExecutionEnvironment(location string, inject []string, full bool) []string {
	environment := conda.CondaExecutionEnvironment(location, inject, full)
	return append(environment,
		it.SearchPath(location).AsEnvironmental("PATH"),
		it.PythonPaths().AsEnvironmental("PYTHONPATH"),
		fmt.Sprintf("ROBOT_ROOT=%s", it.WorkingDirectory()),
		fmt.Sprintf("ROBOT_ARTIFACTS=%s", it.ArtifactDirectory()),
	)
}

func (it *task) shellCommand() []string {
	result, err := shell.Split(it.Shell)
	if err != nil {
		common.Log("Shell parsing failure: %v with command %v", err, it.Shell)
		return []string{}
	}
	return result
}

func (it *task) taskCommand() []string {
	output := "output"
	if it.robot != nil {
		output = it.robot.Artifacts
	}
	return []string{
		"python",
		"-m", "robot",
		"--report", "NONE",
		"--outputdir", output,
		"--logtitle", "Task log",
		"--task", it.Task,
		".",
	}
}

func (it *task) Commandline() []string {
	if len(it.Task) > 0 {
		return it.taskCommand()
	}
	if len(it.Shell) > 0 {
		return it.shellCommand()
	}
	return it.Command
}

func robotFrom(content []byte) (*robot, error) {
	config := robot{}
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	config.relink()
	return &config, nil
}

func PlainEnvironment(inject []string, full bool) []string {
	environment := make([]string, 0, 100)
	if full {
		environment = append(environment, os.Environ()...)
	}
	environment = append(environment, inject...)
	return environment
}

func LoadRobotYaml(filename string, visible bool) (Robot, error) {
	fullpath, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", filename, err)
	}
	content, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", fullpath, err)
	}
	if visible {
		common.Log("%q as robot.yaml is:\n%s", fullpath, string(content))
	}
	robot, err := robotFrom(content)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", fullpath, err)
	}
	robot.Root = filepath.Dir(fullpath)
	return robot, nil
}

func DetectConfigurationName(directory string) string {
	robot, err := pathlib.FindNamedPath(directory, "robot.yaml")
	if err == nil && len(robot) > 0 {
		return robot
	}
	return filepath.Join(directory, "robot.yaml")
}
