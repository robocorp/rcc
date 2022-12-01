package operations

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/htfs"
	"github.com/robocorp/rcc/journal"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/robot"
	"github.com/robocorp/rcc/settings"
	"github.com/robocorp/rcc/xviper"
	"gopkg.in/yaml.v2"
)

const (
	canaryUrl      = `/canary.txt`
	pypiCanaryUrl  = `/jupyterlab-pygments/`
	condaCanaryUrl = `/conda-forge/linux-64/repodata.json`
	statusOk       = `ok`
	statusWarning  = `warning`
	statusFail     = `fail`
	statusFatal    = `fatal`
)

var (
	ignorePathContains = []string{".vscode", ".ipynb_checkpoints", ".virtual_documents"}
)

func shouldIgnorePath(fullpath string) bool {
	lowpath := strings.ToLower(fullpath)
	for _, ignore := range ignorePathContains {
		if strings.Contains(lowpath, ignore) {
			return true
		}
	}
	return false
}

type stringerr func() (string, error)

func justText(source stringerr) string {
	result, _ := source()
	return result
}

func RunDiagnostics() *common.DiagnosticStatus {
	result := &common.DiagnosticStatus{
		Details: make(map[string]string),
		Checks:  []*common.DiagnosticCheck{},
	}
	result.Details["executable"] = common.BinRcc()
	result.Details["rcc"] = common.Version
	result.Details["rcc.bin"] = common.BinRcc()
	result.Details["micromamba"] = conda.MicromambaVersion()
	result.Details["micromamba.bin"] = conda.BinMicromamba()
	result.Details["ROBOCORP_HOME"] = common.RobocorpHome()
	result.Details["ROBOCORP_OVERRIDE_SYSTEM_REQUIREMENTS"] = fmt.Sprintf("%v", common.OverrideSystemRequirements())
	result.Details["RCC_VERBOSE_ENVIRONMENT_BUILDING"] = fmt.Sprintf("%v", common.VerboseEnvironmentBuilding())
	result.Details["user-cache-dir"] = justText(os.UserCacheDir)
	result.Details["user-config-dir"] = justText(os.UserConfigDir)
	result.Details["user-home-dir"] = justText(os.UserHomeDir)
	result.Details["working-dir"] = justText(os.Getwd)
	result.Details["tempdir"] = os.TempDir()
	result.Details["controller"] = common.ControllerIdentity()
	result.Details["user-agent"] = common.UserAgent()
	result.Details["installationId"] = xviper.TrackingIdentity()
	result.Details["telemetry-enabled"] = fmt.Sprintf("%v", xviper.CanTrack())
	result.Details["config-piprc-used"] = fmt.Sprintf("%v", settings.Global.HasPipRc())
	result.Details["config-micromambarc-used"] = fmt.Sprintf("%v", settings.Global.HasMicroMambaRc())
	result.Details["config-settings-yaml-used"] = fmt.Sprintf("%v", pathlib.IsFile(common.SettingsFile()))
	result.Details["config-active-profile"] = settings.Global.Name()
	result.Details["config-https-proxy"] = settings.Global.HttpsProxy()
	result.Details["config-http-proxy"] = settings.Global.HttpProxy()
	result.Details["config-ssl-verify"] = fmt.Sprintf("%v", settings.Global.VerifySsl())
	result.Details["config-ssl-no-revoke"] = fmt.Sprintf("%v", settings.Global.NoRevocation())
	result.Details["os-holo-location"] = common.HoloLocation()
	result.Details["hololib-location"] = common.HololibLocation()
	result.Details["hololib-catalog-location"] = common.HololibCatalogLocation()
	result.Details["hololib-library-location"] = common.HololibLibraryLocation()
	result.Details["holotree-location"] = common.HolotreeLocation()
	result.Details["holotree-shared"] = fmt.Sprintf("%v", common.SharedHolotree)
	result.Details["holotree-user-id"] = common.UserHomeIdentity()
	result.Details["os"] = common.Platform()
	result.Details["cpus"] = fmt.Sprintf("%d", runtime.NumCPU())
	result.Details["when"] = time.Now().Format(time.RFC3339 + " (MST)")
	result.Details["no-build"] = fmt.Sprintf("%v", settings.Global.NoBuild())

	for name, filename := range lockfiles() {
		result.Details[name] = filename
	}

	who, err := user.Current()
	if err == nil {
		result.Details["uid:gid"] = fmt.Sprintf("%s:%s", who.Uid, who.Gid)
	}

	// checks
	if common.SharedHolotree {
		result.Checks = append(result.Checks, verifySharedDirectory(common.HoloLocation()))
		result.Checks = append(result.Checks, verifySharedDirectory(common.HololibLocation()))
		result.Checks = append(result.Checks, verifySharedDirectory(common.HololibCatalogLocation()))
		result.Checks = append(result.Checks, verifySharedDirectory(common.HololibLibraryLocation()))
	}
	result.Checks = append(result.Checks, robocorpHomeCheck())
	result.Checks = append(result.Checks, anyPathCheck("PYTHONPATH"))
	result.Checks = append(result.Checks, anyPathCheck("PLAYWRIGHT_BROWSERS_PATH"))
	result.Checks = append(result.Checks, anyPathCheck("NODE_OPTIONS"))
	result.Checks = append(result.Checks, anyPathCheck("NODE_PATH"))
	if !common.OverrideSystemRequirements() {
		result.Checks = append(result.Checks, longPathSupportCheck())
	}
	result.Checks = append(result.Checks, lockpidsCheck()...)
	result.Checks = append(result.Checks, lockfilesCheck()...)
	hostnames := settings.Global.Hostnames()
	dnsStopwatch := common.Stopwatch("DNS lookup time for %d hostnames was about", len(hostnames))
	for _, host := range hostnames {
		result.Checks = append(result.Checks, dnsLookupCheck(host))
	}
	result.Details["dns-lookup-time"] = dnsStopwatch.Text()
	result.Checks = append(result.Checks, canaryDownloadCheck())
	result.Checks = append(result.Checks, pypiHeadCheck())
	result.Checks = append(result.Checks, condaHeadCheck())
	return result
}

func lockfiles() map[string]string {
	result := make(map[string]string)
	result["lock-config"] = xviper.Lockfile()
	result["lock-cache"] = cacheLockFile()
	result["lock-holotree"] = common.HolotreeLock()
	result["lock-robocorp"] = common.RobocorpLock()
	result["lock-userlock"] = htfs.UserHolotreeLockfile()
	return result
}

func longPathSupportCheck() *common.DiagnosticCheck {
	supportLongPathUrl := settings.Global.DocsLink("troubleshooting/windows-long-path")
	if conda.HasLongPathSupport() {
		return &common.DiagnosticCheck{
			Type:    "OS",
			Status:  statusOk,
			Message: "Supports long enough paths.",
			Link:    supportLongPathUrl,
		}
	}
	return &common.DiagnosticCheck{
		Type:    "OS",
		Status:  statusFail,
		Message: "Does not support long path names!",
		Link:    supportLongPathUrl,
	}
}

func lockfilesCheck() []*common.DiagnosticCheck {
	content := []byte(fmt.Sprintf("lock check %s @%d", common.Version, common.When))
	files := lockfiles()
	result := make([]*common.DiagnosticCheck, 0, len(files))
	support := settings.Global.DocsLink("troubleshooting")
	failed := false
	for identity, filename := range files {
		err := os.WriteFile(filename, content, 0o666)
		if err != nil {
			result = append(result, &common.DiagnosticCheck{
				Type:    "OS",
				Status:  statusFail,
				Message: fmt.Sprintf("Lock file %q write failed, reason: %v", identity, err),
				Link:    support,
			})
			failed = true
		}
	}
	if !failed {
		result = append(result, &common.DiagnosticCheck{
			Type:    "OS",
			Status:  statusOk,
			Message: fmt.Sprintf("%d lockfiles all seem to work correctly (for this user).", len(files)),
			Link:    support,
		})
	}
	return result
}

func lockpidsCheck() []*common.DiagnosticCheck {
	support := settings.Global.DocsLink("troubleshooting")
	result := []*common.DiagnosticCheck{}
	entries, err := os.ReadDir(common.HololibPids())
	if err != nil {
		result = append(result, &common.DiagnosticCheck{
			Type:    "OS",
			Status:  statusWarning,
			Message: fmt.Sprintf("Problem with pids directory: %q, reason: %v", common.HololibPids(), err),
			Link:    support,
		})
		return result
	}
	deadline := time.Now().Add(-12 * time.Hour)
	for _, entry := range entries {
		level, qualifier := statusWarning, "Pending"
		info, err := entry.Info()
		if err == nil && info.ModTime().Before(deadline) {
			level, qualifier = statusOk, "Stale(?)"
		}
		result = append(result, &common.DiagnosticCheck{
			Type:    "OS",
			Status:  level,
			Message: fmt.Sprintf("%s lock file info: %q", qualifier, entry.Name()),
			Link:    support,
		})
	}
	if len(result) == 0 {
		result = append(result, &common.DiagnosticCheck{
			Type:    "OS",
			Status:  statusOk,
			Message: "No pending lock files detected.",
			Link:    support,
		})
	}
	return result
}

func anyPathCheck(key string) *common.DiagnosticCheck {
	supportGeneralUrl := settings.Global.DocsLink("troubleshooting")
	anyPath := os.Getenv(key)
	if len(anyPath) > 0 {
		return &common.DiagnosticCheck{
			Type:    "OS",
			Status:  statusWarning,
			Message: fmt.Sprintf("%s is set to %q. This may cause problems.", key, anyPath),
			Link:    supportGeneralUrl,
		}
	}
	return &common.DiagnosticCheck{
		Type:    "OS",
		Status:  statusOk,
		Message: fmt.Sprintf("%s is not set, which is good.", key),
		Link:    supportGeneralUrl,
	}
}

func verifySharedDirectory(fullpath string) *common.DiagnosticCheck {
	shared := pathlib.IsSharedDir(fullpath)
	supportGeneralUrl := settings.Global.DocsLink("troubleshooting")
	if !shared {
		return &common.DiagnosticCheck{
			Type:    "OS",
			Status:  statusWarning,
			Message: fmt.Sprintf("%q is not shared. This may cause problems.", fullpath),
			Link:    supportGeneralUrl,
		}
	}
	return &common.DiagnosticCheck{
		Type:    "OS",
		Status:  statusOk,
		Message: fmt.Sprintf("%q is shared, which is ok.", fullpath),
		Link:    supportGeneralUrl,
	}
}

func robocorpHomeCheck() *common.DiagnosticCheck {
	supportGeneralUrl := settings.Global.DocsLink("troubleshooting")
	if !conda.ValidLocation(common.RobocorpHome()) {
		return &common.DiagnosticCheck{
			Type:    "RPA",
			Status:  statusFatal,
			Message: fmt.Sprintf("ROBOCORP_HOME (%s) contains characters that makes RPA fail.", common.RobocorpHome()),
			Link:    supportGeneralUrl,
		}
	}
	return &common.DiagnosticCheck{
		Type:    "RPA",
		Status:  statusOk,
		Message: fmt.Sprintf("ROBOCORP_HOME (%s) is good enough.", common.RobocorpHome()),
		Link:    supportGeneralUrl,
	}
}

func dnsLookupCheck(site string) *common.DiagnosticCheck {
	supportNetworkUrl := settings.Global.DocsLink("troubleshooting/firewall-and-proxies")
	found, err := net.LookupHost(site)
	if err != nil {
		return &common.DiagnosticCheck{
			Type:    "network",
			Status:  statusFail,
			Message: fmt.Sprintf("DNS lookup %q failed: %v", site, err),
			Link:    supportNetworkUrl,
		}
	}
	return &common.DiagnosticCheck{
		Type:    "network",
		Status:  statusOk,
		Message: fmt.Sprintf("%s found: %v", site, found),
		Link:    supportNetworkUrl,
	}
}

func condaHeadCheck() *common.DiagnosticCheck {
	supportNetworkUrl := settings.Global.DocsLink("troubleshooting/firewall-and-proxies")
	client, err := cloud.NewClient(settings.Global.CondaLink(""))
	if err != nil {
		return &common.DiagnosticCheck{
			Type:    "network",
			Status:  statusWarning,
			Message: fmt.Sprintf("%v: %v", settings.Global.CondaLink(""), err),
			Link:    supportNetworkUrl,
		}
	}
	request := client.NewRequest(condaCanaryUrl)
	response := client.Head(request)
	if response.Status >= 400 {
		return &common.DiagnosticCheck{
			Type:    "network",
			Status:  statusWarning,
			Message: fmt.Sprintf("Conda canary download failed: %d", response.Status),
			Link:    supportNetworkUrl,
		}
	}
	return &common.DiagnosticCheck{
		Type:    "network",
		Status:  statusOk,
		Message: fmt.Sprintf("Conda canary download successful: %s", settings.Global.CondaLink(condaCanaryUrl)),
		Link:    supportNetworkUrl,
	}
}

func pypiHeadCheck() *common.DiagnosticCheck {
	supportNetworkUrl := settings.Global.DocsLink("troubleshooting/firewall-and-proxies")
	client, err := cloud.NewClient(settings.Global.PypiLink(""))
	if err != nil {
		return &common.DiagnosticCheck{
			Type:    "network",
			Status:  statusWarning,
			Message: fmt.Sprintf("%v: %v", settings.Global.PypiLink(""), err),
			Link:    supportNetworkUrl,
		}
	}
	request := client.NewRequest(pypiCanaryUrl)
	response := client.Head(request)
	if response.Status >= 400 {
		return &common.DiagnosticCheck{
			Type:    "network",
			Status:  statusWarning,
			Message: fmt.Sprintf("PyPI canary download failed: %d", response.Status),
			Link:    supportNetworkUrl,
		}
	}
	return &common.DiagnosticCheck{
		Type:    "network",
		Status:  statusOk,
		Message: fmt.Sprintf("PyPI canary download successful: %s", settings.Global.PypiLink(pypiCanaryUrl)),
		Link:    supportNetworkUrl,
	}
}

func canaryDownloadCheck() *common.DiagnosticCheck {
	supportNetworkUrl := settings.Global.DocsLink("troubleshooting/firewall-and-proxies")
	client, err := cloud.NewClient(settings.Global.DownloadsLink(""))
	if err != nil {
		return &common.DiagnosticCheck{
			Type:    "network",
			Status:  statusFail,
			Message: fmt.Sprintf("%v: %v", settings.Global.DownloadsLink(""), err),
			Link:    supportNetworkUrl,
		}
	}
	request := client.NewRequest(canaryUrl)
	response := client.Get(request)
	if response.Status != 200 || string(response.Body) != "Used to testing connections" {
		return &common.DiagnosticCheck{
			Type:    "network",
			Status:  statusFail,
			Message: fmt.Sprintf("Canary download failed: %d: %s", response.Status, response.Body),
			Link:    supportNetworkUrl,
		}
	}
	return &common.DiagnosticCheck{
		Type:    "network",
		Status:  statusOk,
		Message: fmt.Sprintf("Canary download successful: %s", settings.Global.DownloadsLink(canaryUrl)),
		Link:    supportNetworkUrl,
	}
}

func jsonDiagnostics(sink io.Writer, details *common.DiagnosticStatus) {
	form, err := details.AsJson()
	if err != nil {
		pretty.Exit(1, "Error: %s", err)
	}
	fmt.Fprintln(sink, form)
}

func humaneDiagnostics(sink io.Writer, details *common.DiagnosticStatus) {
	fmt.Fprintln(sink, "Diagnostics:")
	keys := make([]string, 0, len(details.Details))
	for key, _ := range details.Details {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := details.Details[key]
		fmt.Fprintf(sink, " - %-38s...  %q\n", key, value)
	}
	fmt.Fprintln(sink, "")
	fmt.Fprintln(sink, "Checks:")
	for _, check := range details.Checks {
		fmt.Fprintf(sink, " - %-8s %-8s %s\n", check.Type, check.Status, check.Message)
	}
	count, body := journal.MakeStatistics(12, false, false, false, false)
	if count > 4 {
		fmt.Fprintln(sink, "")
		fmt.Fprintln(sink, "Statistics:")
		fmt.Fprintln(sink, "")
		fmt.Fprintln(sink, string(body))
	}
}

func fileIt(filename string) (io.WriteCloser, error) {
	if len(filename) == 0 {
		return os.Stdout, nil
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func ProduceDiagnostics(filename, robotfile string, json, production bool) (*common.DiagnosticStatus, error) {
	file, err := fileIt(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := RunDiagnostics()
	if len(robotfile) > 0 {
		addRobotDiagnostics(robotfile, result, production)
	}
	settings.Global.Diagnostics(result)
	if json {
		jsonDiagnostics(file, result)
	} else {
		humaneDiagnostics(file, result)
	}
	return result, nil
}

type Unmarshaler func([]byte, interface{}) error

func diagnoseFilesUnmarshal(tool Unmarshaler, label, rootdir string, paths []string, target *common.DiagnosticStatus) {
	supportGeneralUrl := settings.Global.DocsLink("troubleshooting")
	target.Details[fmt.Sprintf("%s-file-count", strings.ToLower(label))] = fmt.Sprintf("%d file(s)", len(paths))
	diagnose := target.Diagnose(label)
	var canary interface{}
	success := true
	investigated := false
	for _, tail := range paths {
		investigated = true
		fullpath := filepath.Join(rootdir, tail)
		if shouldIgnorePath(fullpath) {
			continue
		}
		content, err := os.ReadFile(fullpath)
		if err != nil {
			diagnose.Fail(supportGeneralUrl, "Problem reading %s file %q: %v", label, tail, err)
			success = false
			continue
		}
		err = tool(content, &canary)
		if err != nil {
			diagnose.Fail(supportGeneralUrl, "Problem parsing %s file %q: %v", label, tail, err)
			success = false
		}
	}
	if investigated && success {
		diagnose.Ok("%s files are readable and can be parsed.", label)
	}
}

func addFileDiagnostics(rootdir string, target *common.DiagnosticStatus) {
	jsons := pathlib.Glob(rootdir, "*.json")
	diagnoseFilesUnmarshal(json.Unmarshal, "JSON", rootdir, jsons, target)
	yamls := pathlib.Glob(rootdir, "*.yaml")
	yamls = append(yamls, pathlib.Glob(rootdir, "*.yml")...)
	diagnoseFilesUnmarshal(yaml.Unmarshal, "YAML", rootdir, yamls, target)
}

func addRobotDiagnostics(robotfile string, target *common.DiagnosticStatus, production bool) {
	supportGeneralUrl := settings.Global.DocsLink("troubleshooting")
	config, err := robot.LoadRobotYaml(robotfile, false)
	diagnose := target.Diagnose("Robot")
	if err != nil {
		diagnose.Fail(supportGeneralUrl, "About robot.yaml: %v", err)
	} else {
		config.Diagnostics(target, production)
	}
	addFileDiagnostics(filepath.Dir(robotfile), target)
}

func RunRobotDiagnostics(robotfile string, production bool) *common.DiagnosticStatus {
	result := &common.DiagnosticStatus{
		Details: make(map[string]string),
		Checks:  []*common.DiagnosticCheck{},
	}
	addRobotDiagnostics(robotfile, result, production)
	return result
}

func PrintRobotDiagnostics(robotfile string, json, production bool) error {
	result := RunRobotDiagnostics(robotfile, production)
	if json {
		jsonDiagnostics(os.Stdout, result)
	} else {
		humaneDiagnostics(os.Stderr, result)
	}
	return nil
}
