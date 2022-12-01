package conda

import (
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	linebreaks = regexp.MustCompile("\r?\n")
)

func UnifyLine(value string) string {
	return strings.Trim(value, " \t\r\n")
}

func SplitLines(value string) []string {
	return linebreaks.Split(value, -1)
}

func ReadConfig(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func AsUnifiedLines(value string) []string {
	parts := SplitLines(value)
	limit := len(parts)
	seen := make(map[string]bool, limit)
	result := make([]string, 0, limit)
	for _, part := range parts {
		unified := UnifyLine(part)
		if seen[unified] {
			continue
		}
		seen[unified] = true
		if len(unified) > 0 {
			result = append(result, unified)
		}
	}
	sort.Strings(result)
	return result
}
