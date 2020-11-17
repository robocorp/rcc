package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
)

const (
	loadLink = "/robot-v1/workspaces/%s/robots/%s/file/%slink"
)

func linkFor(direction, workspaceId, robotId string) string {
	return fmt.Sprintf(loadLink, workspaceId, robotId, direction)
}

func fetchRobotToken(client cloud.Client, account *account, claims *Claims) (string, error) {
	data, err := AuthorizeCommand(client, account, claims)
	if err != nil {
		return "", err
	}
	token, ok := data["token"].(string)
	if ok {
		return token, nil
	}
	return "", errors.New("Could not get authorization token.")
}

func summonAssistantToken(client cloud.Client, account *account, workspaceId string) (string, error) {
	claims := AssistantClaims(30*60, workspaceId)
	token, ok := account.Cached(claims.Name, claims.Url)
	if ok {
		return token, nil
	}
	return fetchRobotToken(client, account, claims)
}

func summonRobotToken(client cloud.Client, account *account, workspaceId string) (string, error) {
	claims := RobotClaims(30*60, workspaceId)
	token, ok := account.Cached(claims.Name, claims.Url)
	if ok {
		return token, nil
	}
	return fetchRobotToken(client, account, claims)
}

func getAnyloadLink(client cloud.Client, cloudUrl, credentials string) (string, error) {
	request := client.NewRequest(cloudUrl)
	request.Headers[authorization] = BearerToken(credentials)
	response := client.Get(request)
	if response.Status != 200 {
		return "", errors.New(fmt.Sprintf("%d: %s", response.Status, response.Body))
	}
	token := make(Token)
	err := json.Unmarshal(response.Body, &token)
	if err != nil {
		return "", err
	}
	uri, ok := token["uri"]
	if !ok {
		return "", errors.New(fmt.Sprintf("Cannot find URI from %s.", response.Body))
	}
	converted, ok := uri.(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("Cannot find URI as string from %s.", response.Body))
	}
	return converted, nil
}

func putContent(client cloud.Client, awsUrl, zipfile string) error {
	handle, err := os.Open(zipfile)
	if err != nil {
		return err
	}
	defer handle.Close()
	stat, err := handle.Stat()
	if err != nil {
		return err
	}
	request := client.NewRequest(awsUrl)
	request.ContentLength = stat.Size()
	request.TransferEncoding = "identity"
	request.Body = handle
	response := client.Put(request)
	if response.Status != 200 {
		return errors.New(fmt.Sprintf("%d: %s", response.Status, response.Body))
	}
	return nil
}

func getContent(client cloud.Client, awsUrl, zipfile string) error {
	handle, err := os.Create(zipfile)
	if err != nil {
		return err
	}
	defer handle.Close()
	request := client.NewRequest(awsUrl)
	request.Stream = handle
	response := client.Get(request)
	if response.Status != 200 {
		return errors.New(fmt.Sprintf("%d: %s", response.Status, response.Body))
	}
	return nil
}

func UploadCommand(client cloud.Client, account *account, workspaceId, robotId, zipfile string, debug bool) error {
	token, err := summonRobotToken(client, account, workspaceId)
	if err != nil {
		return err
	}
	linkPath := linkFor("upload", workspaceId, robotId)
	targetUrl, err := getAnyloadLink(client, linkPath, token)
	if err != nil {
		return err
	}
	parsed, err := url.Parse(targetUrl)
	if err != nil {
		return err
	}
	awsClient, err := client.NewClient(fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host))
	if err != nil {
		return err
	}
	err = putContent(awsClient, parsed.RequestURI(), zipfile)
	if err != nil {
		return err
	}
	return CacheRobot(zipfile)
}

func DownloadCommand(client cloud.Client, account *account, workspaceId, robotId, zipfile string, debug bool) error {
	token, err := summonRobotToken(client, account, workspaceId)
	if err != nil {
		return err
	}
	linkPath := linkFor("download", workspaceId, robotId)
	targetUrl, err := getAnyloadLink(client, linkPath, token)
	if err != nil {
		return err
	}
	parsed, err := url.Parse(targetUrl)
	if err != nil {
		return err
	}
	awsClient, err := client.NewClient(fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host))
	if err != nil {
		return err
	}
	err = getContent(awsClient, parsed.RequestURI(), zipfile)
	if err != nil {
		return err
	}
	return CacheRobot(zipfile)
}

func SummonRobotZipfile(client cloud.Client, account *account, workspaceId, robotId, digest string) (string, error) {
	found, ok := LookupRobot(digest)
	if ok {
		return found, nil
	}
	zipfile := filepath.Join(os.TempDir(), fmt.Sprintf("summon%x.zip", time.Now().Unix()))
	err := DownloadCommand(client, account, workspaceId, robotId, zipfile, common.DebugFlag)
	if err != nil {
		return "", err
	}
	return zipfile, nil
}
