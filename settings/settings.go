package settings

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/robocorp/rcc/blobs"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
)

const (
	pypiDefault  = "https://pypi.org/simple/"
	condaDefault = "https://conda.anaconda.org/"
)

var (
	httpTransport  *http.Transport
	cachedSettings *Settings
	Global         gateway
	chain          SettingsLayers
)

func cacheSettings(result *Settings) (*Settings, error) {
	if result != nil {
		cachedSettings = result
	}
	return result, nil
}

func HasCustomSettings() bool {
	return pathlib.IsFile(common.SettingsFile())
}

func DefaultSettings() ([]byte, error) {
	return blobs.Asset("assets/settings.yaml")
}

func DefaultSettingsLayer() *Settings {
	content, err := DefaultSettings()
	pretty.Guard(err == nil, 111, "Could not read default settings, reason: %v", err)
	config, err := FromBytes(content)
	pretty.Guard(err == nil, 111, "Could not parse default settings, reason: %v", err)
	return config
}

func CustomSettingsLayer() *Settings {
	if !HasCustomSettings() {
		return nil
	}
	config, err := LoadSetting(common.SettingsFile())
	pretty.Guard(err == nil, 111, "Could not load/parse custom settings, reason: %v", err)
	return config
}

func LoadSetting(filename string) (*Settings, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := FromBytes(content)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func TemporalSettingsLayer(filename string) error {
	config, err := LoadSetting(filename)
	if err != nil {
		return err
	}
	chain[2] = config
	cachedSettings = nil
	return nil
}

func SummonSettings() (*Settings, error) {
	if cachedSettings != nil {
		return cachedSettings, nil
	}
	return cacheSettings(chain.Effective())
}

func showDiagnosticsChecks(sink io.Writer, details *common.DiagnosticStatus) {
	fmt.Fprintln(sink, "Checks:")
	for _, check := range details.Checks {
		fmt.Fprintf(sink, " - %-8s %-8s %s\n", check.Type, check.Status, check.Message)
	}
}

func CriticalEnvironmentSettingsCheck() {
	config, err := SummonSettings()
	pretty.Guard(err == nil, 80, "Aborting! Could not even get setting, reason: %v", err)
	result := &common.DiagnosticStatus{
		Details: make(map[string]string),
		Checks:  []*common.DiagnosticCheck{},
	}
	config.CriticalEnvironmentDiagnostics(result)
	diagnose := result.Diagnose("Settings")
	if HasCustomSettings() {
		diagnose.Ok("Uses custom settings at %q.", common.SettingsFile())
	} else {
		diagnose.Ok("Uses builtin settings.")
	}
	fatal, fail, _, _ := result.Counts()
	if (fatal + fail) > 0 {
		showDiagnosticsChecks(os.Stderr, result)
		pretty.Guard(false, 111, "\nBroken settings.yaml. Cannot continue!")
	}
}

func resolveLink(link, page string) string {
	docs, err := url.Parse(link)
	if err != nil {
		return page
	}
	local, err := url.Parse(page)
	if err != nil {
		return page
	}
	return docs.ResolveReference(local).String()
}

type gateway bool

func (it gateway) settings() *Settings {
	config, err := SummonSettings()
	pretty.Guard(err == nil, 111, "Could not get settings, reason: %v", err)
	return config
}

func (it gateway) Name() string {
	return it.settings().Meta.Name
}

func (it gateway) Description() string {
	return it.settings().Meta.Description
}

func (it gateway) TemplatesYamlURL() string {
	return it.settings().Autoupdates["templates"]
}

func (it gateway) Diagnostics(target *common.DiagnosticStatus) {
	it.settings().Diagnostics(target)
}

func (it gateway) Endpoint(key string) string {
	return it.settings().Endpoints[key]
}

func (it gateway) Option(key string) bool {
	value, ok := it.settings().Options[key]
	return ok && value
}

func (it gateway) DefaultEndpoint() string {
	return it.Endpoint("cloud-api")
}

func (it gateway) IssuesURL() string {
	return it.Endpoint("issues")
}

func (it gateway) TelemetryURL() string {
	return it.Endpoint("telemetry")
}

func (it gateway) PypiURL() string {
	return it.Endpoint("pypi")
}

func (it gateway) PypiTrustedHost() string {
	return justHostAndPort(it.Endpoint("pypi-trusted"))
}

func (it gateway) CondaURL() string {
	return it.Endpoint("conda")
}

func (it gateway) HttpsProxy() string {
	return it.settings().Network.HttpsProxy
}

func (it gateway) HttpProxy() string {
	return it.settings().Network.HttpProxy
}
func (it gateway) HasPipRc() bool {
	return pathlib.IsFile(common.PipRcFile())
}

func (it gateway) HasMicroMambaRc() bool {
	return pathlib.IsFile(common.MicroMambaRcFile())
}

func (it gateway) HasCaBundle() bool {
	return pathlib.IsFile(common.CaBundleFile())
}

func (it gateway) DownloadsLink(resource string) string {
	return resolveLink(it.Endpoint("downloads"), resource)
}

func (it gateway) DocsLink(page string) string {
	return resolveLink(it.Endpoint("docs"), page)
}

func (it gateway) PypiLink(page string) string {
	endpoint := it.Endpoint("pypi")
	if len(endpoint) == 0 {
		endpoint = pypiDefault
	}
	return resolveLink(endpoint, page)
}

func (it gateway) CondaLink(page string) string {
	endpoint := it.Endpoint("conda")
	if len(endpoint) == 0 {
		endpoint = condaDefault
	}
	return resolveLink(endpoint, page)
}

func (it gateway) Hostnames() []string {
	return it.settings().Hostnames()
}
func (it gateway) VerifySsl() bool {
	return it.settings().Certificates.VerifySsl
}

func (it gateway) NoRevocation() bool {
	return it.settings().Certificates.SslNoRevoke
}

func (it gateway) NoBuild() bool {
	nobuild := len(os.Getenv("RCC_NO_BUILD")) > 0
	return nobuild || common.NoBuild || it.Option("no-build")
}

func (it gateway) ConfiguredHttpTransport() *http.Transport {
	return httpTransport
}

func (it gateway) loadRootCAs() *x509.CertPool {
	if !it.HasCaBundle() {
		return nil
	}
	certificates, err := os.ReadFile(common.CaBundleFile())
	if err != nil {
		common.Log("Warning! Problem reading %q, reason: %v", common.CaBundleFile(), err)
		return nil
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(certificates)
	if !ok {
		common.Log("Warning! Problem appending sertificated from %q.", common.CaBundleFile())
		return nil
	}
	return roots
}

func initProtection() {
	status := recover()
	if status != nil {
		fmt.Fprintf(os.Stderr, "Fatal failure with '%v', see: %v\n", common.SettingsFile(), status)
		fmt.Fprintln(os.Stderr, "Sorry. Cannot recover, will exit now!")
		os.Exit(111)
	}
}

func init() {
	defer initProtection()

	chain = SettingsLayers{
		DefaultSettingsLayer(),
		CustomSettingsLayer(),
		nil,
	}
	verifySsl := true
	Global = gateway(true)
	httpTransport = http.DefaultTransport.(*http.Transport).Clone()
	settings, err := SummonSettings()
	if err == nil && settings.Certificates != nil {
		verifySsl = settings.Certificates.VerifySsl
	}
	proxyUrl := ""
	if len(Global.HttpProxy()) > 0 {
		proxyUrl = Global.HttpProxy()
	}
	if len(Global.HttpsProxy()) > 0 {
		proxyUrl = Global.HttpsProxy()
	}
	if len(proxyUrl) > 0 {
		link, err := url.Parse(proxyUrl)
		if err != nil {
			common.Log("Warning! Problem parsing proxy URL %q, reason: %v.", proxyUrl, err)
		} else {
			httpTransport.Proxy = http.ProxyURL(link)
		}
	}
	httpTransport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: !verifySsl,
		RootCAs:            Global.loadRootCAs(),
	}
}
