package conda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/robocorp/rcc/common"
)

const (
	preformatMarker = "```"
	activateFile    = "rcc_activate.json"
)

func capturePreformatted(incoming string) ([]string, string) {
	lines := strings.SplitAfter(incoming, "\n")
	capture := false
	result := make([]string, 0, 5)
	pending := make([]string, 0, len(lines))
	other := make([]string, 0, len(lines))
	for _, line := range lines {
		flat := strings.TrimSpace(line)
		if strings.HasPrefix(flat, preformatMarker) {
			if len(pending) > 0 {
				result = append(result, strings.Join(pending, ""))
				pending = make([]string, 0, len(lines))
			}
			capture = !capture
			continue
		}
		if !capture {
			other = append(other, line)
			continue
		}
		pending = append(pending, line)
	}
	if len(pending) > 0 {
		result = append(result, strings.Join(pending, ""))
	}
	return result, strings.Join(other, "")
}

func createScript(targetFolder string) (string, error) {
	script := template.New("script")
	script, err := script.Parse(activateScript)
	if err != nil {
		return "", err
	}
	details := make(map[string]string)
	details["Rcc"] = common.BinRcc()
	details["Robocorphome"] = common.RobocorpHome()
	details["MambaRootPrefix"] = common.MambaRootPrefix()
	details["Micromamba"] = BinMicromamba()
	details["Live"] = targetFolder
	buffer := bytes.NewBuffer(nil)
	script.Execute(buffer, details)

	scriptfile := filepath.Join(targetFolder, fmt.Sprintf("rcc_activate%s", commandSuffix))
	err = os.WriteFile(scriptfile, buffer.Bytes(), 0o755)
	if err != nil {
		return "", err
	}
	return scriptfile, nil
}

func parseJson(content string) (map[string]string, error) {
	result := make(map[string]string)
	err := json.Unmarshal([]byte(content), &result)
	return result, err
}

func diffStringMaps(before, after map[string]string) map[string]string {
	result := make(map[string]string)
	for key, _ := range before {
		_, ok := after[key]
		if !ok {
			result[key] = ""
		}
	}
	for key, past := range before {
		future, ok := after[key]
		if ok && past != future {
			result[key] = future
		}
	}
	for key, value := range after {
		_, ok := before[key]
		if !ok {
			result[key] = value
		}
	}
	return result
}

func Activate(sink io.Writer, targetFolder string) error {
	envCommand := []string{common.BinRcc(), "internal", "env", "--label", "before"}
	out, _, err := LiveCapture(targetFolder, envCommand...)
	if err != nil {
		return err
	}
	parts, _ := capturePreformatted(out)
	if len(parts) == 0 {
		return fmt.Errorf("Could not detect environment details from 'before' output.")
	}
	before, err := parseJson(parts[0])
	if err != nil {
		return err
	}

	script, err := createScript(targetFolder)
	if err != nil {
		return err
	}

	out, _, err = LiveCapture(targetFolder, script)
	if err != nil {
		fmt.Fprintf(sink, "%v\n%s\n", err, out)
		return err
	}
	parts, other := capturePreformatted(out)
	fmt.Fprintf(sink, "%s\n", other)
	if len(parts) == 0 {
		return fmt.Errorf("Could not detect environment details from 'after' output.")
	}
	after, err := parseJson(parts[0])
	if err != nil {
		return err
	}
	difference := diffStringMaps(before, after)
	body, err := json.MarshalIndent(difference, "", "  ")
	if err != nil {
		return err
	}
	targetJson := filepath.Join(targetFolder, activateFile)
	err = os.WriteFile(targetJson, body, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func LoadActivationEnvironment(targetFolder string) []string {
	result := []string{}
	targetJson := filepath.Join(targetFolder, activateFile)
	content, err := os.ReadFile(targetJson)
	if err != nil {
		return result
	}
	var entries map[string]string
	err = json.Unmarshal(content, &entries)
	if err != nil {
		return result
	}
	for name, value := range entries {
		result = append(result, fmt.Sprintf("%s=%s", name, value))
	}
	common.Trace("Environment activation added %d variables.", len(result))
	return result
}
