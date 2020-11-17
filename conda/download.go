package conda

import (
	"crypto/sha256"
	"io"
	"net/http"
	"os"

	"github.com/robocorp/rcc/common"
)

func DownloadConda() error {
	url := DownloadLink()
	filename := DownloadTarget()
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	digest := sha256.New()
	many := io.MultiWriter(out, digest)

	common.Debug("Downloading %s <%s> -> %s", url, response.Status, filename)

	_, err = io.Copy(many, response.Body)
	if err != nil {
		return err
	}

	return common.Debug("SHA256 sum: %02x", digest.Sum(nil))
}
