package conda

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/shell"
	"github.com/robocorp/rcc/xviper"
)

func chooseBestEnvironment(best int, selected, reference string, candidates []string) (int, string) {
	for _, candidate := range candidates {
		move, err := Distance(reference, candidate)
		if err != nil {
			continue
		}
		if move < best {
			best, selected = move, candidate
		}
	}

	return best, selected
}

func Hexdigest(raw []byte) string {
	return fmt.Sprintf("%02x", raw)
}

func metafile(folder string) string {
	return ExpandPath(folder + ".meta")
}

func metaLoad(location string) (string, error) {
	raw, err := ioutil.ReadFile(metafile(location))
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func metaSave(location, data string) error {
	return ioutil.WriteFile(metafile(location), []byte(data), 0644)
}

func touchMetafile(location string) {
	pathlib.TouchWhen(metafile(location), time.Now())
}

func LastUsed(location string) (time.Time, error) {
	return pathlib.Modtime(metafile(location))
}

func IsPristine(folder string) bool {
	digest, err := DigestFor(folder)
	if err != nil {
		return false
	}
	meta, err := metaLoad(folder)
	if err != nil {
		return false
	}
	return Hexdigest(digest) == meta
}

func reuseExistingLive(key string) bool {
	candidate := LiveFrom(key)
	if IsPristine(candidate) {
		touchMetafile(candidate)
		return true
	}
	removeClone(candidate)
	return false
}

func LiveExecution(liveFolder string, command ...string) error {
	searchPath := FindPath(liveFolder)
	commandName := command[0]
	task, ok := searchPath.Which(commandName, FileExtensions)
	if !ok {
		return errors.New(fmt.Sprintf("Cannot find command: %v", commandName))
	}
	common.Debug("Using %v as command %v.", task, commandName)
	command[0] = task
	environment := EnvironmentFor(liveFolder)
	_, err := shell.New(environment, ".", command...).Transparent()
	return err
}

func newLive(condaYaml, requirementsText, key string, force, freshInstall bool) bool {
	targetFolder := LiveFrom(key)
	when := time.Now()
	if force {
		when = when.Add(-20 * 24 * time.Hour)
	}
	if force || !freshInstall {
		common.Log("rcc touching conda cache. (Stamp: %v)", when)
		SilentTouch(CondaCache(), when)
	}
	common.Debug("Setting up new conda environment using %v to folder %v", condaYaml, targetFolder)
	command := []string{CondaExecutable(), "env", "create", "-q", "-f", condaYaml, "-p", targetFolder}
	if common.DebugFlag {
		command = []string{CondaExecutable(), "env", "create", "-f", condaYaml, "-p", targetFolder}
	}
	_, err := shell.New(nil, ".", command...).Transparent()
	if err != nil {
		common.Error("Conda error", err)
		return false
	}
	common.Debug("Updating new environment at %v with pip requirements from %v", targetFolder, requirementsText)
	pipCommand := []string{"pip", "install", "--no-color", "--disable-pip-version-check", "--prefer-binary", "--cache-dir", PipCache(), "--find-links", WheelCache(), "--requirement", requirementsText, "--quiet"}
	if common.DebugFlag {
		pipCommand = []string{"pip", "install", "--no-color", "--disable-pip-version-check", "--prefer-binary", "--cache-dir", PipCache(), "--find-links", WheelCache(), "--requirement", requirementsText}
	}
	err = LiveExecution(targetFolder, pipCommand...)
	if err != nil {
		common.Error("Pip error", err)
		return false
	}
	digest, err := DigestFor(targetFolder)
	if err != nil {
		common.Error("Digest", err)
		return false
	}
	return metaSave(targetFolder, Hexdigest(digest)) == nil
}

func temporaryConfig(condaYaml, requirementsText string, filenames ...string) (string, error) {
	var left, right *Environment
	var err error

	for _, filename := range filenames {
		left = right
		right, err = ReadCondaYaml(filename)
		if err != nil {
			return "", err
		}
		if left == nil {
			continue
		}
		right, err = left.Merge(right)
		if err != nil {
			return "", err
		}
	}
	yaml, err := right.AsYaml()
	if err != nil {
		return "", err
	}
	hash, err := LocalitySensitiveHash(AsUnifiedLines(yaml))
	if err != nil {
		return "", err
	}
	err = right.SaveAsRequirements(requirementsText)
	if err != nil {
		return "", err
	}
	pure := right.AsPureConda()
	return hash, pure.SaveAs(condaYaml)
}

func NewEnvironment(force bool, configurations ...string) (string, error) {
	requests := xviper.GetInt("stats.env.request") + 1
	hits := xviper.GetInt("stats.env.hit")
	dirty := xviper.GetInt("stats.env.dirty")
	misses := xviper.GetInt("stats.env.miss")
	failures := xviper.GetInt("stats.env.failures")
	merges := xviper.GetInt("stats.env.merges")
	templates := len(TemplateList())
	freshInstall := templates == 0

	defer func() {
		common.Log("####  Progress: 4/4  [Done.] [Stats: %d environments, %d requests, %d merges, %d hits, %d dirty, %d misses, %d failures | %s]", templates, requests, merges, hits, dirty, misses, failures, common.Version)
	}()
	common.Log("####  Progress: 0/4  [try use existing live same environment?] %v", xviper.TrackingIdentity())

	xviper.Set("stats.env.request", requests)

	if len(configurations) > 1 {
		merges += 1
		xviper.Set("stats.env.merges", merges)
	}

	marker := time.Now().Unix()
	condaYaml := filepath.Join(os.TempDir(), fmt.Sprintf("conda_%x.yaml", marker))
	requirementsText := filepath.Join(os.TempDir(), fmt.Sprintf("require_%x.txt", marker))
	common.Debug("Using temporary conda.yaml file: %v and requirement.txt file: %v", condaYaml, requirementsText)
	key, err := temporaryConfig(condaYaml, requirementsText, configurations...)
	if err != nil {
		failures += 1
		xviper.Set("stats.env.failures", failures)
		return "", err
	}
	defer os.Remove(condaYaml)
	defer os.Remove(requirementsText)

	liveFolder := LiveFrom(key)
	if reuseExistingLive(key) {
		hits += 1
		xviper.Set("stats.env.hit", hits)
		return liveFolder, nil
	}
	common.Log("####  Progress: 1/4  [try clone existing same template to live, key: %v]", key)
	if CloneFromTo(TemplateFrom(key), liveFolder) {
		dirty += 1
		xviper.Set("stats.env.dirty", dirty)
		return liveFolder, nil
	}
	common.Log("####  Progress: 2/4  [try create new environment from scratch]")
	if newLive(condaYaml, requirementsText, key, force, freshInstall) {
		misses += 1
		xviper.Set("stats.env.miss", misses)
		common.Log("####  Progress: 3/4  [backup new environment as template]")
		CloneFromTo(liveFolder, TemplateFrom(key))
		return liveFolder, nil
	}

	failures += 1
	xviper.Set("stats.env.failures", failures)
	return "", errors.New("Could not create environment.")
}

func RemoveEnvironment(label string) {
	removeClone(LiveFrom(label))
	removeClone(TemplateFrom(label))
}

func removeClone(location string) {
	os.Remove(metafile(location))
	os.RemoveAll(location)
}

func CloneFromTo(source, target string) bool {
	removeClone(target)
	os.MkdirAll(target, 0755)

	if !IsPristine(source) {
		removeClone(source)
		return false
	}
	expected, err := metaLoad(source)
	if err != nil {
		return false
	}
	success := cloneFolder(source, target, 8)
	if !success {
		removeClone(target)
		return false
	}
	digest, err := DigestFor(target)
	if err != nil || Hexdigest(digest) != expected {
		removeClone(target)
		return false
	}
	metaSave(target, expected)
	touchMetafile(source)
	return true
}

func cloneFolder(source, target string, workers int) bool {
	queue := make(chan copyRequest)
	done := make(chan bool)

	for x := 0; x < workers; x++ {
		go copyWorker(queue, done)
	}

	success := copyFolder(source, target, queue)
	close(queue)

	for x := 0; x < workers; x++ {
		<-done
	}

	return success
}

func SilentTouch(directory string, when time.Time) bool {
	handle, err := os.Open(directory)
	if err != nil {
		return false
	}
	entries, err := handle.Readdir(-1)
	handle.Close()
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			pathlib.TouchWhen(filepath.Join(directory, entry.Name()), when)
		}
	}
	return true
}

func copyFolder(source, target string, queue chan copyRequest) bool {
	os.Mkdir(target, 0755)

	handle, err := os.Open(source)
	if err != nil {
		common.Error("OPEN", err)
		return false
	}
	entries, err := handle.Readdir(-1)
	handle.Close()
	if err != nil {
		common.Error("DIR", err)
		return false
	}

	success := true
	expect := 0
	for _, entry := range entries {
		if entry.Name() == "__pycache__" {
			continue
		}
		newSource := filepath.Join(source, entry.Name())
		newTarget := filepath.Join(target, entry.Name())
		if entry.IsDir() {
			copyFolder(newSource, newTarget, queue)
		} else {
			queue <- copyRequest{newSource, newTarget}
			expect += 1
		}
	}

	return success
}

type copyRequest struct {
	source, target string
}

func copyWorker(tasks chan copyRequest, done chan bool) {
	for {
		task, ok := <-tasks
		if !ok {
			break
		}
		link, err := os.Readlink(task.source)
		if err != nil {
			pathlib.CopyFile(task.source, task.target, false)
			continue
		}
		err = os.Symlink(link, task.target)
		if err != nil {
			common.Error("LINK", err)
			continue
		}
	}

	done <- true
}
