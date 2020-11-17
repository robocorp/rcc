package pathlib

import (
	"io"
	"os"
	"path/filepath"

	"github.com/robocorp/rcc/common"
)

func CopyFile(source, target string, overwrite bool) error {
	targetDir := filepath.Dir(target)
	err := os.MkdirAll(targetDir, 0o755)
	if err != nil {
		return err
	}
	if overwrite && Exists(target) {
		err = os.Remove(target)
	}
	if err != nil {
		return err
	}
	readable, err := os.Open(source)
	if err != nil {
		return err
	}
	defer readable.Close()
	stats, err := readable.Stat()
	if err != nil {
		return err
	}
	writable, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_EXCL, stats.Mode())
	if err != nil {
		return err
	}
	defer writable.Close()

	_, err = io.Copy(writable, readable)
	if err != nil {
		common.Error("copy-file", err)
	}

	return err
}
