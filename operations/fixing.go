package operations

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/robot"
)

var (
	nonExecutableExtensions = make(map[string]bool)
)

func init() {
	nonExecutableExtensions[".md"] = true
	nonExecutableExtensions[".txt"] = true
	nonExecutableExtensions[".htm"] = true
	nonExecutableExtensions[".html"] = true
}

func ToUnix(content []byte) []byte {
	parts := bytes.Split(content, []byte{'\r'})
	return bytes.Join(parts, []byte{})
}

func fixShellFile(fullpath string) {
	content, err := ioutil.ReadFile(fullpath)
	if err != nil || bytes.IndexByte(content, '\r') < 0 {
		return
	}
	common.Debug("Fixing newlines in file: %v", fullpath)
	err = ioutil.WriteFile(fullpath, ToUnix(content), 0o755)
	if err != nil {
		common.Log("Failure %v while fixing newlines in %v!", err, fullpath)
	}
}

func makeExecutable(fullpath string, file os.FileInfo) {
	extension := strings.ToLower(filepath.Ext(file.Name()))
	ignore, ok := nonExecutableExtensions[extension]
	if ok && ignore || file.Mode() == 0o755 {
		return
	}
	os.Chmod(fullpath, 0o755)
	common.Debug("Making file executable: %v", fullpath)
}

func ensureFilesExecutable(dir string) {
	handle, err := os.Open(dir)
	if err != nil {
		return
	}
	defer handle.Close()
	files, err := handle.Readdir(-1)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fullpath := filepath.Join(dir, file.Name())
		extension := strings.ToLower(filepath.Ext(fullpath))
		if extension == ".sh" {
			fixShellFile(fullpath)
		}
		makeExecutable(fullpath, file)
	}
}

func FixRobot(robotFile string) error {
	config, err := robot.LoadYamlConfiguration(robotFile)
	if err != nil {
		return err
	}
	tasks := config.AvailableTasks()
	for _, task := range tasks {
		for _, path := range config.Paths(task) {
			ensureFilesExecutable(path)
		}
	}
	return nil
}

func FixDirectory(dir string) error {
	primary := filepath.Join(dir, "robot.yaml")
	if pathlib.IsFile(primary) {
		return FixRobot(primary)
	}
	secondary := filepath.Join(dir, "package.yaml")
	if pathlib.IsFile(secondary) {
		return FixRobot(secondary)
	}
	return nil
}
