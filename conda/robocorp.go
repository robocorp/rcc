package conda

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/shell"
	"github.com/robocorp/rcc/xviper"
)

var (
	ignoredPaths     = []string{"python", "conda"}
	hashPattern      = regexp.MustCompile("^[0-9a-f]{16}(?:\\.meta)?$")
	randomIdentifier string
)

func init() {
	randomIdentifier = fmt.Sprintf("%016x", rand.Uint64()^uint64(os.Getpid()))
}

func sorted(files []os.FileInfo) {
	sort.SliceStable(files, func(left, right int) bool {
		return files[left].Name() < files[right].Name()
	})
}

func ignoreDynamicDirectories(folder, entryName string) bool {
	base := strings.ToLower(filepath.Base(folder))
	name := strings.ToLower(entryName)
	return name == "__pycache__" || (name == "gen" && base == "comtypes")
}

func DigestFor(folder string, collect map[string]string) ([]byte, error) {
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
			if ignoreDynamicDirectories(folder, entry.Name()) {
				continue
			}
			digest, err := DigestFor(filepath.Join(folder, entry.Name()), collect)
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
	if collect != nil {
		key := fmt.Sprintf("%02x", result)
		collect[folder] = key
	}
	return result, nil
}

func hashedEntity(name string) bool {
	return hashPattern.MatchString(name)
}

func hasDatadir(basedir, metafile string) bool {
	if filepath.Ext(metafile) != ".meta" {
		return false
	}
	fullpath := filepath.Join(basedir, metafile)
	stat, err := os.Stat(fullpath[:len(fullpath)-5])
	return err == nil && stat.IsDir()
}

func hasMetafile(basedir, subdir string) bool {
	folder := filepath.Join(basedir, subdir)
	_, err := os.Stat(metafile(folder))
	return err == nil
}

func dirnamesFrom(location string) []string {
	result := make([]string, 0, 20)
	handle, err := os.Open(common.ExpandPath(location))
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

func orphansFrom(location string) []string {
	result := make([]string, 0, 20)
	handle, err := os.Open(common.ExpandPath(location))
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
		hashed := hashedEntity(child.Name())
		if hashed && child.IsDir() && hasMetafile(location, child.Name()) {
			continue
		}
		if hashed && !child.IsDir() && hasDatadir(location, child.Name()) {
			continue
		}
		result = append(result, filepath.Join(location, child.Name()))
	}

	return result
}

func FindPath(environment string) pathlib.PathParts {
	target := pathlib.TargetPath()
	target = target.Remove(ignoredPaths)
	target = target.Prepend(CondaPaths(environment)...)
	return target
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
	environment = append(environment,
		"CONDA_DEFAULT_ENV=rcc",
		"CONDA_PREFIX="+location,
		"CONDA_PROMPT_MODIFIER=(rcc) ",
		"CONDA_SHLVL=1",
		"PYTHONHOME=",
		"PYTHONSTARTUP=",
		"PYTHONEXECUTABLE=",
		"PYTHONNOUSERSITE=1",
		"PYTHONDONTWRITEBYTECODE=x",
		"PYTHONPYCACHEPREFIX="+RobocorpTemp(),
		"ROBOCORP_HOME="+common.RobocorpHome(),
		"RCC_ENVIRONMENT_HASH="+common.EnvironmentHash,
		"RCC_INSTALLATION_ID="+xviper.TrackingIdentity(),
		"RCC_TRACKING_ALLOWED="+fmt.Sprintf("%v", xviper.CanTrack()),
		"TEMP="+RobocorpTemp(),
		"TMP="+RobocorpTemp(),
		searchPath.AsEnvironmental("PATH"),
	)
	environment = append(environment, LoadActivationEnvironment(location)...)
	return environment
}

func EnvironmentFor(location string) []string {
	return append(os.Environ(), EnvironmentExtensionFor(location)...)
}

func MambaPackages() string {
	return common.ExpandPath(filepath.Join(common.RobocorpHome(), "pkgs"))
}

func MambaCache() string {
	return common.ExpandPath(filepath.Join(MambaPackages(), "cache"))
}

func asVersion(text string) (uint64, string) {
	text = strings.TrimSpace(text)
	multiline := strings.SplitN(text, "\n", 2)
	if len(multiline) > 0 {
		text = strings.TrimSpace(multiline[0])
	}
	parts := strings.SplitN(text, ".", 4)
	steps := len(parts)
	multipliers := []uint64{1000000, 1000, 1}
	version := uint64(0)
	for at, multiplier := range multipliers {
		if steps <= at {
			break
		}
		value, err := strconv.ParseUint(parts[at], 10, 64)
		if err != nil {
			break
		}
		version += multiplier * value
	}
	return version, text
}

func MicromambaVersion() string {
	versionText, _, err := shell.New(CondaEnvironment(), ".", BinMicromamba(), "--repodata-ttl", "90000", "--version").CaptureOutput()
	if err != nil {
		return err.Error()
	}
	_, versionText = asVersion(versionText)
	return versionText
}

func HasMicroMamba() bool {
	if !pathlib.IsFile(BinMicromamba()) {
		return false
	}
	version, versionText := asVersion(MicromambaVersion())
	goodEnough := version >= 14000
	common.Debug("%q version is %q -> %v (good enough: %v)", BinMicromamba(), versionText, version, goodEnough)
	common.Timeline("µmamba version is %q (at %q).", versionText, BinMicromamba())
	return goodEnough
}

func RobocorpTempRoot() string {
	return filepath.Join(common.RobocorpHome(), "temp")
}

func RobocorpTemp() string {
	tempLocation := filepath.Join(RobocorpTempRoot(), randomIdentifier)
	fullpath, err := pathlib.EnsureDirectory(tempLocation)
	if err != nil {
		common.Log("WARNING (%v) -> %v", tempLocation, err)
	}
	return fullpath
}

func MinicondaLocation() string {
	// Legacy function, but must remain until cleanup is done
	return filepath.Join(common.RobocorpHome(), "miniconda3")
}

func LocalChannel() (string, bool) {
	basefolder := filepath.Join(common.RobocorpHome(), "channel")
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
	return filepath.Join(common.BaseLocation(), hash)
}

func LiveFrom(hash string) string {
	if common.Stageonly {
		return common.StageFolder
	}
	return common.ExpandPath(filepath.Join(common.LiveLocation(), hash))
}

func TemplateList() []string {
	return dirnamesFrom(common.BaseLocation())
}

func LiveList() []string {
	return dirnamesFrom(common.LiveLocation())
}

func OrphanList() []string {
	result := orphansFrom(common.BaseLocation())
	result = append(result, orphansFrom(common.LiveLocation())...)
	return result
}

func InstallationPlan(hash string) (string, bool) {
	finalplan := filepath.Join(LiveFrom(hash), "rcc_plan.log")
	return finalplan, pathlib.IsFile(finalplan)
}
