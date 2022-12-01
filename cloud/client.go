package cloud

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/settings"
	"github.com/robocorp/rcc/xviper"
)

type internalClient struct {
	endpoint string
	client   *http.Client
}

type Request struct {
	Url              string
	Headers          map[string]string
	TransferEncoding string
	ContentLength    int64
	Body             io.Reader
	Stream           io.Writer
}

type Response struct {
	Status  int
	Err     error
	Body    []byte
	Elapsed common.Duration
}

type Client interface {
	Endpoint() string
	NewRequest(string) *Request
	Head(request *Request) *Response
	Get(request *Request) *Response
	Post(request *Request) *Response
	Put(request *Request) *Response
	Delete(request *Request) *Response
	NewClient(endpoint string) (Client, error)
	WithTimeout(time.Duration) Client
}

func EnsureHttps(endpoint string) (string, error) {
	nice := strings.TrimRight(strings.TrimSpace(endpoint), "/")
	parsed, err := url.Parse(nice)
	if err != nil {
		return "", err
	}
	if parsed.Host == "127.0.0.1" || strings.HasPrefix(parsed.Host, "127.0.0.1:") {
		return nice, nil
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("Endpoint '%s' must start with https:// prefix.", nice)
	}
	return nice, nil
}

func NewClient(endpoint string) (Client, error) {
	https, err := EnsureHttps(endpoint)
	if err != nil {
		return nil, err
	}
	return &internalClient{
		endpoint: https,
		client:   &http.Client{Transport: settings.Global.ConfiguredHttpTransport()},
	}, nil
}

func (it *internalClient) WithTimeout(timeout time.Duration) Client {
	return &internalClient{
		endpoint: it.endpoint,
		client: &http.Client{
			Transport: settings.Global.ConfiguredHttpTransport(),
			Timeout:   timeout,
		},
	}
}

func (it *internalClient) NewClient(endpoint string) (Client, error) {
	return NewClient(endpoint)
}

func (it *internalClient) Endpoint() string {
	return it.endpoint
}

func (it *internalClient) does(method string, request *Request) *Response {
	stopwatch := common.Stopwatch("stopwatch")
	response := new(Response)
	url := it.Endpoint() + request.Url
	common.Trace("Doing %s %s", method, url)
	defer func() {
		response.Elapsed = stopwatch.Elapsed()
		common.Trace("%s %s took %s", method, url, response.Elapsed)
	}()
	httpRequest, err := http.NewRequest(method, url, request.Body)
	if err != nil {
		response.Status = 9001
		response.Err = err
		return response
	}
	if request.ContentLength > 0 {
		httpRequest.ContentLength = request.ContentLength
	}
	if len(request.TransferEncoding) > 0 {
		httpRequest.TransferEncoding = []string{request.TransferEncoding}
	}
	httpRequest.Header.Add("robocorp-installation-id", xviper.TrackingIdentity())
	httpRequest.Header.Add("User-Agent", common.UserAgent())
	for name, value := range request.Headers {
		httpRequest.Header.Add(name, value)
	}
	httpResponse, err := it.client.Do(httpRequest)
	if err != nil {
		common.Error("http.Do", err)
		response.Status = 9002
		response.Err = err
		return response
	}
	defer httpResponse.Body.Close()
	response.Status = httpResponse.StatusCode
	if request.Stream != nil {
		io.Copy(request.Stream, httpResponse.Body)
	} else {
		response.Body, response.Err = io.ReadAll(httpResponse.Body)
	}
	if common.DebugFlag {
		body := "ignore"
		if response.Status > 399 {
			body = string(response.Body)
		}
		common.Debug("%v %v %v => %v (%v)", <-common.Identities, method, url, response.Status, body)
	}
	return response
}

func (it *internalClient) NewRequest(url string) *Request {
	return &Request{
		Url:     url,
		Headers: make(map[string]string),
	}
}

func (it *internalClient) Head(request *Request) *Response {
	return it.does("HEAD", request)
}

func (it *internalClient) Get(request *Request) *Response {
	return it.does("GET", request)
}

func (it *internalClient) Post(request *Request) *Response {
	return it.does("POST", request)
}

func (it *internalClient) Put(request *Request) *Response {
	return it.does("PUT", request)
}

func (it *internalClient) Delete(request *Request) *Response {
	return it.does("DELETE", request)
}

func Download(url, filename string) error {
	common.Timeline("start %s download", filename)
	defer common.Timeline("done %s download", filename)

	if pathlib.Exists(filename) {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
	}

	client := &http.Client{Transport: settings.Global.ConfiguredHttpTransport()}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Accept", "application/octet-stream")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("Downloading %q failed, reason: %q!", url, response.Status)
	}

	pathlib.EnsureDirectory(filepath.Dir(filename))
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	digest := sha256.New()
	many := io.MultiWriter(out, digest)

	common.Debug("Downloading %s <%s> -> %s", url, response.Status, filename)

	bytecount, err := io.Copy(many, response.Body)
	if err != nil {
		return err
	}

	common.Timeline("downloaded %d bytes to %s", bytecount, filename)

	err = out.Sync()
	if err != nil {
		return err
	}

	return common.Debug("%q SHA256 sum: %02x", filename, digest.Sum(nil))
}
