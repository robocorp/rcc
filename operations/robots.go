package operations

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/robocorp/rcc/common"
)

func UpdateRobot(directory string) error {
	fullpath, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	cache, err := SummonCache()
	if err != nil {
		return err
	}
	defer cache.Save()
	robot, ok := cache.Robots[fullpath]
	if !ok {
		robot = &Folder{
			Path:    fullpath,
			Created: common.When,
			Updated: common.When,
			Deleted: 0,
		}
		cache.Robots[fullpath] = robot
	}
	stat, err := os.Stat(fullpath)
	if err != nil || !stat.IsDir() {
		robot.Deleted = common.When
	}
	robot.Updated = common.When
	return nil
}

func sorted(folders []*Folder) {
	sort.SliceStable(folders, func(left, right int) bool {
		if folders[left].Deleted != folders[right].Deleted {
			return folders[left].Deleted < folders[right].Deleted
		}
		return folders[left].Updated > folders[right].Updated
	})
}

func detectDeadRobots() bool {
	cache, err := SummonCache()
	if err != nil {
		return false
	}
	changed := false
	for _, robot := range cache.Robots {
		stat, err := os.Stat(robot.Path)
		if err != nil || !stat.IsDir() {
			robot.Deleted = common.When
			changed = true
			continue
		}
		if robot.Deleted > 0 && stat.IsDir() {
			robot.Deleted = 0
			changed = true
		}
	}
	if changed {
		cache.Save()
	}
	return changed
}

func ListRobots() ([]*Folder, error) {
	detectDeadRobots()
	cache, err := SummonCache()
	if err != nil {
		return nil, err
	}
	result := make([]*Folder, 0, len(cache.Robots))
	for _, value := range cache.Robots {
		result = append(result, value)
	}
	sorted(result)
	return result, nil
}
