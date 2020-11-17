package conda

import (
	"errors"
	"fmt"
	"io/ioutil"
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
}

type Environment struct {
	Name     string
	Prefix   string
	Channels []string
	Conda    []*Dependency
	Pip      []*Dependency
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
		return nil, errors.New(fmt.Sprintf("Not same component: %v vs. %v", it.Name, right.Name))
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
	return nil, errors.New(fmt.Sprintf("Wont choose between dependencies: %v vs. %v", it.Original, right.Original))
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
		Name:   it.Name,
		Prefix: it.Prefix,
	}
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
		Channels: []string{"defaults", "conda-forge"},
		Conda:    []*Dependency{},
		Pip:      []*Dependency{},
	}
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

func addChannels(seen map[string]bool, source, target []string) []string {
	for _, channel := range source {
		found := seen[channel]
		if !found {
			seen[channel] = true
			target = append(target, channel)
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
	result.Name = it.Name + "+" + right.Name
	seenChannels := make(map[string]bool)
	result.Channels = addChannels(seenChannels, it.Channels, result.Channels)
	result.Channels = addChannels(seenChannels, right.Channels, result.Channels)
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
	return ioutil.WriteFile(filename, []byte(content), 0o640)
}

func (it *Environment) SaveAsRequirements(filename string) error {
	content := it.AsRequirementsText()
	return ioutil.WriteFile(filename, []byte(content), 0o640)
}

func (it *Environment) AsYaml() (string, error) {
	result := new(internalEnvironment)
	result.Name = it.Name
	result.Prefix = it.Prefix
	result.Channels = it.Channels
	result.Dependencies = it.CondaList()
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

func CondaYamlFrom(content []byte) (*Environment, error) {
	result := new(internalEnvironment)
	err := yaml.Unmarshal(content, result)
	if err != nil {
		return nil, err
	}
	return result.AsEnvironment(), nil
}

func ReadCondaYaml(filename string) (*Environment, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
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
