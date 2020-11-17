package operations

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/robocorp/rcc/blobs"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
)

func unpack(content []byte, directory string) error {
	common.Debug("Initializing:")
	size := int64(len(content))
	byter := bytes.NewReader(content)
	reader, err := zip.NewReader(byter, size)
	if err != nil {
		return err
	}
	success := true
	for _, entry := range reader.File {
		if entry.FileInfo().IsDir() {
			continue
		}
		target := filepath.Join(directory, entry.Name)
		todo := WriteTarget{
			Source: entry,
			Target: target,
		}
		success = todo.Execute() && success
	}
	common.Debug("Done.")
	if !success {
		return errors.New(fmt.Sprintf("Problems while initializing robot. Use --debug to see details."))
	}
	return nil
}

func ListTemplates() []string {
	assets := blobs.AssetNames()
	result := make([]string, 0, len(assets))
	for _, name := range blobs.AssetNames() {
		if !strings.HasPrefix(name, "assets") || !strings.HasSuffix(name, ".zip") {
			continue
		}
		result = append(result, strings.TrimSuffix(filepath.Base(name), filepath.Ext(name)))
	}
	sort.Strings(result)
	return result
}

func InitializeWorkarea(directory, name string, force bool) error {
	content, err := blobs.Asset(fmt.Sprintf("assets/%s.zip", name))
	if err != nil {
		return err
	}
	fullpath, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	if force {
		err = pathlib.EnsureDirectoryExists(fullpath)
	} else {
		err = pathlib.EnsureEmptyDirectory(fullpath)
	}
	if err != nil {
		return err
	}
	UpdateRobot(fullpath)
	return unpack(content, fullpath)
}
