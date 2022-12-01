package conda

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"

	"gopkg.in/yaml.v2"
)

var (
	dependencyPattern = regexp.MustCompile("^([^<=~!> ]+)\\s*(?:([<=~!>]*)\\s*(\\S+.*?))?$")
)

type internalEnvironment struct {
	Name         string        `yaml:"name,omitempty"`
	Channels     []string      `yaml:"channels"`
	Dependencies []interface{} `yaml:"dependencies"`
	Prefix       string        `yaml:"prefix,omitempty"`
	PostInstall  []string      `yaml:"rccPostInstall,omitempty"`
}

type Environment struct {
	Name        string
	Prefix      string
	Channels    []string
	Conda       []*Dependency
	Pip         []*Dependency
	PostInstall []string
}

type Dependency struct {
	Original  string
	Name      string
	Qualifier string
	Versions  string
}

func AsDependency(value string) *Dependency {
	trimmed := strings.TrimSpace(value)
	parts := dependencyPattern.FindStringSubmatch(trimmed)
	if len(parts) != 4 {
		return nil
	}
	return &Dependency{
		Original:  parts[0],
		Name:      parts[1],
		Qualifier: parts[2],
		Versions:  parts[3],
	}
}

func (it *Dependency) Representation() string {
	parts := strings.SplitN(strings.ToLower(it.Name), "[", 2)
	return parts[0]
}

func (it *Dependency) IsExact() bool {
	return len(it.Qualifier)+len(it.Versions) > 0
}

func (it *Dependency) SameAs(right *Dependency) bool {
	return it.Name == right.Name
}

func (it *Dependency) ExactlySame(right *Dependency) bool {
	return it.Name == right.Name && it.Qualifier == right.Qualifier && it.Versions == right.Versions
}

func (it *Dependency) ChooseSpecific(right *Dependency) (*Dependency, error) {
	if !it.SameAs(right) {
		return nil, fmt.Errorf("Not same component: %v vs. %v", it.Name, right.Name)
	}
	if it.IsExact() && !right.IsExact() {
		return it, nil
	}
	if !it.IsExact() && right.IsExact() {
		return right, nil
	}
	if it.ExactlySame(right) {
		return it, nil
	}
	return nil, fmt.Errorf("Wont choose between dependencies: %v vs. %v", it.Original, right.Original)
}

func (it *Dependency) Index(others []*Dependency) int {
	for index, other := range others {
		if it.SameAs(other) {
			return index
		}
	}
	return -1
}

func (it *internalEnvironment) AsEnvironment() *Environment {
	result := &Environment{
		Name:        it.Name,
		Prefix:      it.Prefix,
		PostInstall: []string{},
	}
	seenScripts := make(map[string]bool)
	result.PostInstall = addItem(seenScripts, it.PostInstall, result.PostInstall)
	channel, ok := LocalChannel()
	if ok {
		pushChannels(result, []string{channel})
		common.Debug("Using local conda channel from: %v", channel)
	}
	pushChannels(result, it.Channels)
	pushConda(result, it.condaDependencies())
	pushPip(result, it.pipDependencies())
	result.pipPromote()
	return result
}

func remove(index int, target []*Dependency) []*Dependency {
	return append(target[:index], target[index+1:]...)
}

func SummonEnvironment(filename string) *Environment {
	if pathlib.IsFile(filename) {
		result, err := ReadCondaYaml(filename)
		if err == nil {
			return result
		}
	}
	return &Environment{
		Channels: []string{"conda-forge"},
		Conda:    []*Dependency{},
		Pip:      []*Dependency{},
	}
}

func (it *Environment) FreezeDependencies(fixed dependencies) *Environment {
	result := &Environment{
		Name:        it.Name,
		Prefix:      it.Prefix,
		Channels:    it.Channels,
		Conda:       []*Dependency{},
		Pip:         []*Dependency{},
		PostInstall: it.PostInstall,
	}
	used := make(map[string]bool)
	for _, dependency := range fixed {
		if dependency.Origin == "pypi" {
			continue
		}
		if used[dependency.Name] {
			continue
		}
		used[dependency.Name] = true
		result.Conda = append(result.Conda, &Dependency{
			Original:  fmt.Sprintf("%s=%s", dependency.Name, dependency.Version),
			Name:      dependency.Name,
			Qualifier: "=",
			Versions:  dependency.Version,
		})
	}
	for _, dependency := range fixed {
		if dependency.Origin != "pypi" {
			continue
		}
		if used[dependency.Name] {
			continue
		}
		used[dependency.Name] = true
		result.Pip = append(result.Pip, &Dependency{
			Original:  fmt.Sprintf("%s==%s", dependency.Name, dependency.Version),
			Name:      dependency.Name,
			Qualifier: "==",
			Versions:  dependency.Version,
		})
	}
	return result
}

func (it *Environment) FromDependencies(fixed dependencies) (*Environment, bool) {
	result := &Environment{
		Name:        it.Name,
		Prefix:      it.Prefix,
		Channels:    it.Channels,
		Conda:       []*Dependency{},
		Pip:         []*Dependency{},
		PostInstall: it.PostInstall,
	}
	same := true
	for _, dependency := range it.Conda {
		found, ok := fixed.Lookup(dependency.Name, false)
		if !ok {
			result.Conda = append(result.Conda, dependency)
			same = false
			common.Debug("Could not fix version for dependency %q from conda.", dependency.Name)
			continue
		}
		result.Conda = append(result.Conda, &Dependency{
			Original:  fmt.Sprintf("%s=%s", dependency.Name, found.Version),
			Name:      dependency.Name,
			Qualifier: "=",
			Versions:  found.Version,
		})
	}
	for _, dependency := range it.Pip {
		found, ok := fixed.Lookup(dependency.Name, true)
		if !ok {
			result.Conda = append(result.Pip, dependency)
			same = false
			common.Debug("Could not fix version for dependency %q from pypi.", dependency.Name)
			continue
		}
		result.Pip = append(result.Pip, &Dependency{
			Original:  fmt.Sprintf("%s==%s", dependency.Name, found.Version),
			Name:      dependency.Name,
			Qualifier: "==",
			Versions:  found.Version,
		})
	}
	return result, same
}

func (it *Environment) pipPromote() error {
	removed := make([]int, 0, len(it.Pip))

	for pindex, pip := range it.Pip {
		for cindex, conda := range it.Conda {
			if pip.SameAs(conda) {
				removed = append(removed, pindex)
				chosen, err := conda.ChooseSpecific(pip)
				if err != nil {
					return err
				}
				it.Conda[cindex] = chosen
			}
		}
	}

	for index, value := range removed {
		it.Pip = remove(value-index, it.Pip)
	}

	return nil
}

func (it *internalEnvironment) pipDependencies() []*Dependency {
	result := make([]*Dependency, 0, len(it.Dependencies))
	for _, item := range it.Dependencies {
		group, ok := item.(map[interface{}]interface{})
		if !ok {
			continue
		}
		for key, value := range group {
			if key != "pip" {
				continue
			}
			result = pipContent(result, value)
		}
	}
	return result
}

func (it *internalEnvironment) condaDependencies() []*Dependency {
	result := make([]*Dependency, 0, len(it.Dependencies))
	for _, item := range it.Dependencies {
		value, ok := item.(string)
		if !ok {
			continue
		}
		dependency := AsDependency(value)
		if dependency != nil {
			result = append(result, dependency)
		}
	}
	return result
}

func addItem(seen map[string]bool, source, target []string) []string {
	for _, item := range source {
		found := seen[item]
		if !found {
			seen[item] = true
			target = append(target, item)
		}
	}
	return target
}

func pushChannels(target *Environment, channels []string) {
	for _, channel := range channels {
		target.PushChannel(channel)
	}
}

func pushConda(target *Environment, dependencies []*Dependency) error {
	for _, dependency := range dependencies {
		err := target.PushConda(dependency)
		if err != nil {
			return err
		}
	}
	return nil
}

func pushPip(target *Environment, dependencies []*Dependency) error {
	for _, dependency := range dependencies {
		err := target.PushPip(dependency)
		if err != nil {
			return err
		}
	}
	return nil
}

func (it *Environment) Merge(right *Environment) (*Environment, error) {
	result := new(Environment)
	if len(it.Name) > 0 || len(right.Name) > 0 {
		result.Name = it.Name + "+" + right.Name
	}

	seenChannels := make(map[string]bool)
	result.Channels = addItem(seenChannels, it.Channels, result.Channels)
	result.Channels = addItem(seenChannels, right.Channels, result.Channels)

	seenScripts := make(map[string]bool)
	result.PostInstall = addItem(seenScripts, it.PostInstall, result.PostInstall)
	result.PostInstall = addItem(seenScripts, right.PostInstall, result.PostInstall)

	err := pushConda(result, it.Conda)
	if err != nil {
		return nil, err
	}
	err = pushConda(result, right.Conda)
	if err != nil {
		return nil, err
	}
	err = pushPip(result, it.Pip)
	if err != nil {
		return nil, err
	}
	err = pushPip(result, right.Pip)
	if err != nil {
		return nil, err
	}
	err = result.pipPromote()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func semiSmartPush(target []*Dependency, candidate *Dependency) ([]*Dependency, error) {
	for index, value := range target {
		if value.SameAs(candidate) {
			chosen, err := value.ChooseSpecific(candidate)
			if err != nil {
				return nil, err
			}
			target[index] = chosen
			return target, nil
		}
	}
	target = append(target, candidate)
	return target, nil
}

func (it *Environment) PushChannel(channel string) {
	for _, existing := range it.Channels {
		if channel == existing {
			return
		}
	}
	it.Channels = append(it.Channels, channel)
}

func (it *Environment) PushConda(dependency *Dependency) error {
	result, err := semiSmartPush(it.Conda, dependency)
	if err != nil {
		return err
	}
	it.Conda = result
	return nil
}

func (it *Environment) PushPip(dependency *Dependency) error {
	result, err := semiSmartPush(it.Pip, dependency)
	if err != nil {
		return err
	}
	it.Pip = result
	return nil
}

func asList(source []*Dependency) []interface{} {
	result := make([]interface{}, 0, len(source)+1)
	for _, value := range source {
		result = append(result, value.Original)
	}
	return result
}

func (it *Environment) CondaList() []interface{} {
	return asList(it.Conda)
}

func (it *Environment) PipList() []interface{} {
	return asList(it.Pip)
}

func (it *Environment) PipMap() map[interface{}]interface{} {
	result := make(map[interface{}]interface{})
	result["pip"] = it.PipList()
	return result
}

func (it *Environment) AsPureConda() *Environment {
	return &Environment{
		Name:     it.Name,
		Prefix:   it.Prefix,
		Channels: it.Channels,
		Conda:    it.Conda,
		Pip:      []*Dependency{},
	}
}

func (it *Environment) SaveAs(filename string) error {
	content, err := it.AsYaml()
	if err != nil {
		return err
	}
	common.Trace("FINAL conda environment file as %v:\n---\n%v---", filename, content)
	return os.WriteFile(filename, []byte(content), 0o640)
}

func (it *Environment) SaveAsRequirements(filename string) error {
	content := it.AsRequirementsText()
	common.Trace("FINAL pip requirements as %v:\n---\n%v\n---", filename, content)
	return os.WriteFile(filename, []byte(content), 0o640)
}

func (it *Environment) AsYaml() (string, error) {
	result := new(internalEnvironment)
	result.Name = it.Name
	result.Prefix = it.Prefix
	result.Channels = it.Channels
	result.Dependencies = it.CondaList()
	seenScripts := make(map[string]bool)
	result.PostInstall = addItem(seenScripts, it.PostInstall, result.PostInstall)
	if len(it.Pip) > 0 {
		result.Dependencies = append(result.Dependencies, it.PipMap())
	}
	content, err := yaml.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (it *Environment) AsRequirementsText() string {
	lines := make([]string, 0, len(it.Pip))
	for _, entry := range it.Pip {
		lines = append(lines, entry.Original)
	}
	return strings.Join(lines, Newline)
}

func (it *Environment) Diagnostics(target *common.DiagnosticStatus, production bool) {
	diagnose := target.Diagnose("Conda")
	notice := diagnose.Warning
	if production {
		notice = diagnose.Fail
	}
	packages := make(map[string]bool)
	countChannels := len(it.Channels)
	defaultsPostion := -1
	floating := false
	ok := true
	for index, channel := range it.Channels {
		if channel == "defaults" {
			defaultsPostion = index
			diagnose.Warning("", "Try to avoid defaults channel, and prefer using conda-forge instead.")
			ok = false
		}
	}
	if defaultsPostion == 0 && countChannels > 1 {
		diagnose.Warning("", "Try to avoid putting defaults channel as first channel.")
		ok = false
	}
	if countChannels > 1 {
		diagnose.Warning("", "Try to avoid multiple channel. They may cause problems with code compatibility.")
		ok = false
	}
	if ok {
		diagnose.Ok("Channels in conda.yaml are ok.")
	}
	ok = true
	for _, dependency := range it.Conda {
		presentation := dependency.Representation()
		if packages[presentation] {
			notice("", "Dependency %q seems to be duplicate of previous dependency.", dependency.Original)
		}
		packages[presentation] = true
		if strings.Contains(dependency.Versions, "*") || len(dependency.Qualifier) == 0 || len(dependency.Versions) == 0 {
			notice("", "Floating conda dependency %q should be bound to exact version before taking robot into production.", dependency.Original)
			ok = false
			floating = true
		}
		if len(dependency.Qualifier) > 0 && !(dependency.Qualifier == "==" || dependency.Qualifier == "=") {
			diagnose.Fail("", "Conda dependency %q must use '==' or '=' for version declaration.", dependency.Original)
			ok = false
			floating = true
		}
	}
	if ok {
		diagnose.Ok("Conda dependencies in conda.yaml are ok.")
	}
	ok = true
	for _, dependency := range it.Pip {
		presentation := dependency.Representation()
		if packages[presentation] {
			notice("", "Dependency %q seems to be duplicate of previous dependency.", dependency.Original)
		}
		packages[presentation] = true
		if strings.Contains(dependency.Versions, "*") || len(dependency.Qualifier) == 0 || len(dependency.Versions) == 0 {
			notice("", "Floating pip dependency %q should be bound to exact version before taking robot into production.", dependency.Original)
			ok = false
			floating = true
		}
		if len(dependency.Qualifier) > 0 && dependency.Qualifier != "==" {
			diagnose.Fail("", "Pip dependency %q must use '==' for version declaration.", dependency.Original)
			ok = false
			floating = true
		}
	}
	if ok {
		diagnose.Ok("Pip dependencies in conda.yaml are ok.")
	}
	if floating {
		diagnose.Warning("", "Floating dependencies in Robocorp Cloud containers will be slow, because floating environments cannot be cached.")
	}
}

func CondaYamlFrom(content []byte) (*Environment, error) {
	result := new(internalEnvironment)
	err := yaml.Unmarshal(content, result)
	if err != nil {
		return nil, err
	}
	return result.AsEnvironment(), nil
}

func ReadCondaYaml(filename string) (*Environment, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", filename, err)
	}
	return CondaYamlFrom(content)
}

func pipContent(result []*Dependency, value interface{}) []*Dependency {
	values, ok := value.([]interface{})
	if !ok {
		return result
	}
	for _, entry := range values {
		item, ok := entry.(string)
		if !ok {
			continue
		}
		dependency := AsDependency(item)
		if dependency != nil {
			result = append(result, dependency)
		}
	}
	return result
}
