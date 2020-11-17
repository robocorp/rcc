package conda

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
)

const (
	ROBOCORP_HOME_VARIABLE = `ROBOCORP_HOME`
)

var (
	ignoredPaths = []string{"python", "conda"}
	pythonPaths  = []string{"resources", "libraries", "tasks", "variables"}
)

func sorted(files []os.FileInfo) {
	sort.SliceStable(files, func(left, right int) bool {
		return files[left].Name() < files[right].Name()
	})
}

func DigestFor(folder string) ([]byte, error) {
	handle, err := os.Open(folder)
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	entries, err := handle.Readdir(-1)
	if err != nil {
		return nil, err
	}
	digester := sha256.New()
	sorted(entries)
	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "__pycache__" {
				continue
			}
			digest, err := DigestFor(filepath.Join(folder, entry.Name()))
			if err != nil {
				return nil, err
			}
			digester.Write(digest)
			continue
		}
		repr := fmt.Sprintf("%s -- %x", entry.Name(), entry.Size())
		digester.Write([]byte(repr))
	}
	result := digester.Sum([]byte{})
	return result, nil
}

func hasMetafile(basedir, subdir string) bool {
	folder := filepath.Join(basedir, subdir)
	_, err := os.Stat(metafile(folder))
	return err == nil
}

func dirnamesFrom(location string) []string {
	result := make([]string, 0, 20)
	handle, err := os.Open(ExpandPath(location))
	if err != nil {
		common.Error("Warning", err)
		return result
	}
	defer handle.Close()
	children, err := handle.Readdir(-1)
	if err != nil {
		common.Error("Warning", err)
		return result
	}

	for _, child := range children {
		if child.IsDir() && hasMetafile(location, child.Name()) {
			result = append(result, child.Name())
		}
	}

	return result
}

func FindPath(environment string) pathlib.PathParts {
	target := pathlib.TargetPath()
	target = target.Remove(ignoredPaths)
	target = target.Prepend(CondaPaths(environment)...)
	return target
}

func PythonPath() pathlib.PathParts {
	return pathlib.PathFrom(pythonPaths...)
}

func EnvironmentExtensionFor(location string) []string {
	environment := make([]string, 0, 20)
	searchPath := FindPath(location)
	python, ok := searchPath.Which("python3", FileExtensions)
	if !ok {
		python, ok = searchPath.Which("python", FileExtensions)
	}
	if ok {
		environment = append(environment, "PYTHON_EXE="+python)
	}
	return append(environment,
		"CONDA_DEFAULT_ENV=rcc",
		"CONDA_EXE="+BinConda(),
		"CONDA_PREFIX="+location,
		"CONDA_PROMPT_MODIFIER=(rcc)",
		"CONDA_PYTHON_EXE="+BinPython(),
		"CONDA_SHLVL=1",
		"PYTHONHOME=",
		"PYTHONSTARTUP=",
		"PYTHONEXECUTABLE=",
		"PYTHONNOUSERSITE=1",
		"ROBOCORP_HOME="+RobocorpHome(),
		searchPath.AsEnvironmental("PATH"),
		PythonPath().AsEnvironmental("PYTHONPATH"),
	)
}

func EnvironmentFor(location string) []string {
	return append(os.Environ(), EnvironmentExtensionFor(location)...)
}

func CondaExecutable() string {
	return ExpandPath(filepath.Join(MinicondaLocation(), "condabin", "conda"))
}

func CondaCache() string {
	return ExpandPath(filepath.Join(MinicondaLocation(), "pkgs", "cache"))
}

func HasConda() bool {
	location := ExpandPath(filepath.Join(MinicondaLocation(), "condabin"))
	stat, err := os.Stat(location)
	if err == nil && stat.IsDir() {
		return true
	}
	return false
}

func RobocorpHome() string {
	home := os.Getenv(ROBOCORP_HOME_VARIABLE)
	if len(home) > 0 {
		return ExpandPath(home)
	}
	return ExpandPath(defaultRobocorpLocation)
}

func LiveLocation() string {
	return filepath.Join(RobocorpHome(), "live")
}

func TemplateLocation() string {
	return filepath.Join(RobocorpHome(), "base")
}

func MinicondaLocation() string {
	return filepath.Join(RobocorpHome(), "miniconda3")
}

func ensureDirectory(name string) string {
	pathlib.EnsureDirectoryExists(name)
	return name
}

func PipCache() string {
	return ensureDirectory(filepath.Join(RobocorpHome(), "pipcache"))
}

func WheelCache() string {
	return ensureDirectory(filepath.Join(RobocorpHome(), "wheels"))
}

func RobotCache() string {
	return ensureDirectory(filepath.Join(RobocorpHome(), "robots"))
}

func LocalChannel() (string, bool) {
	basefolder := filepath.Join(RobocorpHome(), "channel")
	fullpath := filepath.Join(basefolder, "channeldata.json")
	stats, err := os.Stat(fullpath)
	if err != nil {
		return "", false
	}
	if !stats.IsDir() {
		return basefolder, true
	}
	return "", false
}

func TemplateFrom(hash string) string {
	return filepath.Join(TemplateLocation(), hash)
}

func LiveFrom(hash string) string {
	return ExpandPath(filepath.Join(LiveLocation(), hash))
}

func TemplateList() []string {
	return dirnamesFrom(TemplateLocation())
}

func LiveList() []string {
	return dirnamesFrom(LiveLocation())
}
