package conda

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"gopkg.in/yaml.v2"
)

type (
	packageDependencies struct {
		CondaForge []string `yaml:"conda-forge,omitempty"`
		Pypi       []string `yaml:"pypi,omitempty"`
	}
	internalPackage struct {
		Dependencies    *packageDependencies `yaml:"dependencies"`
		DevDependencies *packageDependencies `yaml:"dev-dependencies"`
		PostInstall     []string             `yaml:"post-install,omitempty"`
	}
)

func (it *internalPackage) AsEnvironment(devDependencies bool) *Environment {
	result := &Environment{
		Channels:    []string{"conda-forge"},
		PostInstall: []string{},
	}
	seenScripts := make(map[string]bool)
	result.PostInstall = addItem(seenScripts, it.PostInstall, result.PostInstall)
	pushConda(result, it.condaDependencies(it.Dependencies))
	pushPip(result, it.pipDependencies(it.Dependencies))
	if devDependencies {
		pushConda(result, it.condaDependencies(it.DevDependencies))
		pushPip(result, it.pipDependencies(it.DevDependencies))
	}
	result.pipPromote()
	return result
}

func fixPipDependency(dependency *Dependency) *Dependency {
	if dependency != nil {
		if dependency.Qualifier == "=" {
			dependency.Original = fmt.Sprintf("%s==%s", dependency.Name, dependency.Versions)
			dependency.Qualifier = "=="
		}
	}
	return dependency
}

func (it *internalPackage) pipDependencies(useDependencies *packageDependencies) []*Dependency {
	if useDependencies != nil {
		result := make([]*Dependency, 0, len(useDependencies.Pypi))
		for _, item := range useDependencies.Pypi {
			dependency := AsDependency(item)
			if dependency != nil {
				result = append(result, fixPipDependency(dependency))
			}
		}
		return result
	} else {
		return []*Dependency{}
	}
}

func (it *internalPackage) condaDependencies(useDependencies *packageDependencies) []*Dependency {
	if useDependencies != nil {
		result := make([]*Dependency, 0, len(useDependencies.CondaForge))
		for _, item := range useDependencies.CondaForge {
			dependency := AsDependency(item)
			if dependency != nil {
				result = append(result, dependency)
			}
		}
		return result
	} else {
		return []*Dependency{}
	}
}

func packageYamlFrom(content []byte, devDependencies bool) (*Environment, error) {
	result := new(internalPackage)
	err := yaml.Unmarshal(content, result)
	if err != nil {
		return nil, err
	}
	return result.AsEnvironment(devDependencies), nil
}

func ReadPackageCondaYaml(filename string, devDependencies bool) (*Environment, error) {
	basename := strings.ToLower(filepath.Base(filename))
	if basename == "package.yaml" {
		environment, err := readPackageYaml(filename, devDependencies)
		if err == nil {
			return environment, nil
		}
	}
	if devDependencies {
		// error: only valid when dealing with a `package.yaml` file
		return nil, fmt.Errorf("'--devdeps' flag is only valid when dealing with a `package.yaml` file. Current file: %q", filename)
	}
	return readCondaYaml(filename)
}

func readPackageYaml(filename string, devDependencies bool) (*Environment, error) {
	if devDependencies {
		common.Debug("Reading file %q with dev dependencies", filename)
	}
	var content []byte
	var err error

	if pathlib.IsFile(filename) {
		content, err = os.ReadFile(filename)
	} else {
		content, err = cloud.ReadFile(filename)
	}
	if err != nil {
		return nil, fmt.Errorf("%q: %w", filename, err)
	}
	return packageYamlFrom(content, devDependencies)
}
