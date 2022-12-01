package robot

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Setup map[string]string

func (it Setup) AsEnvironment() []string {
	if it == nil {
		return []string{}
	}
	result := make([]string, 0, len(it))
	for key, value := range it {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}

func EnvironmentSetupFrom(content []byte) (Setup, error) {
	setup := make(Setup)
	err := yaml.Unmarshal(content, &setup)
	if err != nil {
		return nil, err
	}
	return setup, nil
}

func LoadEnvironmentSetup(filename string) (Setup, error) {
	if filename == "" {
		return nil, nil
	}
	fullpath, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", filename, err)
	}
	content, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", fullpath, err)
	}
	return EnvironmentSetupFrom(content)
}
