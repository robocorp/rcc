package pathlib

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Forced func(os.FileInfo) bool
type Ignore func(os.FileInfo) bool
type Report func(string, string, os.FileInfo)

type IgnoreOlder time.Time

func (it IgnoreOlder) Ignore(candidate os.FileInfo) bool {
	return candidate.ModTime().Before(time.Time(it))
}

type IgnoreNewer time.Time

func (it IgnoreNewer) Ignore(candidate os.FileInfo) bool {
	return candidate.ModTime().After(time.Time(it))
}

func IgnoreNothing(_ os.FileInfo) bool {
	return false
}

func IgnoreDirectories(target os.FileInfo) bool {
	return target.IsDir()
}

func ForceNothing(_ os.FileInfo) bool {
	return false
}

func ForceFilename(filename string) Forced {
	return func(target os.FileInfo) bool {
		return !target.IsDir() && target.Name() == filename
	}
}

func NoReporting(string, string, os.FileInfo) {
}

func sorted(files []os.FileInfo) {
	sort.SliceStable(files, func(left, right int) bool {
		return files[left].Name() < files[right].Name()
	})
}

type composite []Ignore

func (it composite) Ignore(file os.FileInfo) bool {
	for _, ignore := range it {
		if ignore(file) {
			return true
		}
	}
	return false
}

type exactIgnore string

func (it exactIgnore) Ignore(file os.FileInfo) bool {
	return file.Name() == string(it)
}

type globIgnore string

func (it globIgnore) Ignore(file os.FileInfo) bool {
	name := file.Name()
	result, err := filepath.Match(string(it), name)
	if err == nil && result {
		return true
	}
	if file.IsDir() {
		result, err = filepath.Match(string(it), name+"/")
		return err == nil && result
	}
	return false
}

func CompositeIgnore(ignores ...Ignore) Ignore {
	return composite(ignores).Ignore
}

func IgnorePattern(text string) Ignore {
	return CompositeIgnore(exactIgnore(text).Ignore, globIgnore(text).Ignore)
}

func LoadIgnoreFile(filename string) (Ignore, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	result := make([]Ignore, 0, 10)
	for _, line := range strings.SplitAfter(string(content), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		result = append(result, IgnorePattern(line))
	}
	return CompositeIgnore(result...), nil
}

func LoadIgnoreFiles(filenames []string) (Ignore, error) {
	result := make([]Ignore, 0, len(filenames))
	for _, filename := range filenames {
		ignore, err := LoadIgnoreFile(filename)
		if err != nil {
			return nil, err
		}
		result = append(result, ignore)
	}
	return CompositeIgnore(result...), nil
}

func folderEntries(directory string) ([]os.FileInfo, error) {
	handle, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	entries, err := handle.Readdir(-1)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func recursiveDirWalk(here os.FileInfo, directory, prefix string, report Report) error {
	entries, err := folderEntries(directory)
	if err != nil {
		return err
	}
	sorted(entries)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		nextPrefix := filepath.Join(prefix, entry.Name())
		entryPath := filepath.Join(directory, entry.Name())
		recursiveDirWalk(entry, entryPath, nextPrefix, report)
	}
	report(directory, prefix, here)
	return nil
}

func recursiveWalk(directory, prefix string, force Forced, ignore Ignore, report Report) error {
	entries, err := folderEntries(directory)
	if err != nil {
		return err
	}
	sorted(entries)
	for _, entry := range entries {
		if !force(entry) && ignore(entry) {
			continue
		}
		nextPrefix := filepath.Join(prefix, entry.Name())
		entryPath := filepath.Join(directory, entry.Name())
		if entry.IsDir() {
			recursiveWalk(entryPath, nextPrefix, force, ignore, report)
		} else {
			report(entryPath, nextPrefix, entry)
		}
	}
	return nil
}

func ForceWalk(directory string, force Forced, ignore Ignore, report Report) error {
	fullpath, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	return recursiveWalk(fullpath, ".", force, ignore, report)
}

func DirWalk(directory string, report Report) error {
	fullpath, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	entry, err := os.Stat(fullpath)
	if err != nil {
		return err
	}
	return recursiveDirWalk(entry, fullpath, ".", report)
}

func Walk(directory string, ignore Ignore, report Report) error {
	return ForceWalk(directory, ForceNothing, ignore, report)
}

func Glob(directory string, pattern string) []string {
	result := []string{}
	ignore := func(entry os.FileInfo) bool {
		match, err := filepath.Match(pattern, entry.Name())
		return err != nil || !entry.IsDir() && !match
	}
	capture := func(_, localpath string, _ os.FileInfo) {
		result = append(result, localpath)
	}
	ForceWalk(directory, ForceNothing, ignore, capture)
	return result
}
