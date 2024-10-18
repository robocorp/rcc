package htfs_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robocorp/rcc/common"

	"github.com/robocorp/rcc/hamlet"
	"github.com/robocorp/rcc/htfs"
)

func TestHTFSspecification(t *testing.T) {
	must, wont := hamlet.Specifications(t)

	filename := filepath.Join(os.TempDir(), "htfs_test.json")

	fs, err := htfs.NewRoot("..")
	must.Nil(err)
	wont.Nil(fs)
	wont.Nil(fs.Tree)

	must.Nil(fs.Lift())

	content, err := fs.AsJson()
	must.Nil(err)
	must.True(len(content) > 50000)

	must.Nil(fs.SaveAs(filename))

	reloaded, err := htfs.NewRoot(".")
	must.Nil(err)
	wont.Nil(reloaded)
	before, err := reloaded.AsJson()
	must.Nil(err)
	must.True(len(before) < 300)
	wont.Equal(fs.Path, reloaded.Path)
	must.Nil(reloaded.LoadFrom(filename))
	after, err := reloaded.AsJson()
	must.Nil(err)
	must.Equal(len(after), len(content))
	must.Equal(fs.Path, reloaded.Path)
}

// This test case depends on runtime.GOARCH being "amd64" - this is enforced
// when running unit tests with invoke, but if the test suite is run otherwise,
// for example directly from the IDE, GOARCH env variable needs to be set in order
// for this test to pass.
func TestZipLibrary(t *testing.T) {
	must, wont := hamlet.Specifications(t)

	platform := common.Platform()
	var zipFileName string

	switch {
	case strings.Contains(platform, "linux"):
		zipFileName = "simple_linux.zip"
	case strings.Contains(platform, "darwin"):
		zipFileName = "simple_darwin.zip"
	case strings.Contains(platform, "windows"):
		zipFileName = "simple_windows.zip"
	}

	_, blueprint, err := htfs.ComposeFinalBlueprint([]string{"testdata/simple.yaml"}, "", false)
	must.Nil(err)
	wont.Nil(blueprint)
	sut, err := htfs.ZipLibrary(fmt.Sprintf("testdata/%s", zipFileName))
	must.Nil(err)
	wont.Nil(sut)
	must.True(sut.HasBlueprint(blueprint))
}
