package pathlib_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/robocorp/rcc/hamlet"
	"github.com/robocorp/rcc/pathlib"
)

func TestCalculateSha256OfFiles(t *testing.T) {
	must, wont := hamlet.Specifications(t)

	digest, err := pathlib.Sha256("testdata/missing")
	wont.Nil(err)
	must.Equal("", digest)

	digest, err = pathlib.Sha256("testdata/empty")
	must.Nil(err)
	must.Equal("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", digest)

	contents, err := os.ReadFile("testdata/hello.txt")
	must.Nil(err)
	// check that the file has no CR characters
	must.Equal(-1, bytes.Index(contents, []byte("\r")))

	digest, err = pathlib.Sha256("testdata/hello.txt")
	must.Nil(err)
	must.Equal("d9014c4624844aa5bac314773d6b689ad467fa4e1d1a50a1b8a99d5a95f72ff5", digest)
}
