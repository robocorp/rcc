package operations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/conda"
	"github.com/robocorp/rcc/pathlib"
)

var (
	zipPattern = regexp.MustCompile("^[0-9a-f]{64}\\.zip$")
)

func CacheRobot(filename string) error {
	fullpath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	digest, err := pathlib.Sha256(fullpath)
	if err != nil {
		return err
	}
	common.Debug("Digest for %v is %v.", fullpath, digest)
	_, exists := LookupRobot(digest)
	if exists {
		return nil
	}
	target := cacheRobotFilename(digest)
	err = pathlib.CopyFile(fullpath, target, false)
	if err != nil {
		return err
	}
	verify, err := pathlib.Sha256(target)
	if err != nil {
		defer os.Remove(target)
		return err
	}
	if verify != digest {
		defer os.Remove(target)
		return errors.New(fmt.Sprintf("Could not cache %v, reason: digest mismatch.", fullpath))
	}
	go CleanupOldestRobot()
	return nil
}

func cacheRobotFilename(digest string) string {
	return filepath.Join(conda.RobotCache(), digest+".zip")
}

func LookupRobot(digest string) (string, bool) {
	target := cacheRobotFilename(digest)
	if pathlib.IsFile(target) {
		pathlib.TouchWhen(target, time.Now())
		return target, true
	}
	return "", false
}

func CleanupOldestRobot() {
	oldest, _ := OldestRobot()
	if pathlib.IsFile(oldest) {
		common.Debug("Removing oldest cached robot %v.", oldest)
		os.Remove(oldest)
	}
}

func OldestRobot() (string, time.Time) {
	oldest, stamp := "", time.Now()
	deadline := time.Now().Add(-35 * 24 * time.Hour)
	pathlib.Walk(conda.RobotCache(), pathlib.IgnoreNewer(deadline).Ignore, func(full, relative string, details os.FileInfo) {
		if zipPattern.MatchString(details.Name()) && details.ModTime().Before(stamp) {
			oldest, stamp = full, details.ModTime()
		}
	})
	return oldest, stamp
}
