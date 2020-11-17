package operations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
)

const (
	listAssistantsApi = `/assistant-v1/workspaces/%s/assistants`
	startAssistantApi = `/assistant-v1/workspaces/%s/assistants/%s/runs`
	stopAssistantApi  = `/assistant-v1/workspaces/%s/assistants/%s/runs/%s/complete`
	beatAssistantApi  = `/assistant-v1/workspaces/%s/assistants/%s/runs/%s/heartbeat`
)

type awsPostInfo struct {
	Url    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

type awsResponse struct {
	ArtifactId string       `json:"artifactId"`
	PostInfo   *awsPostInfo `json:"postInfo"`
}

type awsWrapper struct {
	Response *awsResponse `json:"response"`
}

type AssistantRobot struct {
	WorkspaceId string
	AssistantId string
	RunId       string
	RobotId     string
	TaskName    string
	Zipfile     string
	Environment map[string]string
	Details     map[string]interface{}
	Config      map[string]interface{}
	ArtifactURL string
}

func (it *AssistantRobot) extractEnvironment(source Token) {
	blob, ok := source["config"]
	if !ok {
		return
	}
	config, ok := blob.(map[string]interface{})
	if !ok {
		return
	}
	blob, ok = config["environment"]
	if !ok {
		return
	}
	environment, ok := blob.(map[string]interface{})
	if !ok {
		return
	}
	for key, value := range environment {
		it.Environment[key] = fmt.Sprintf("%v", value)
	}
}

type ArtifactPublisher struct {
	Client          cloud.Client
	ArtifactPostURL string
	ErrorCount      int
}

func (it *ArtifactPublisher) NewClient(targetUrl string) (cloud.Client, *url.URL, error) {
	parsed, err := url.Parse(targetUrl)
	if err != nil {
		return nil, nil, err
	}
	newClient, err := it.Client.NewClient(fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host))
	if err != nil {
		return nil, nil, err
	}
	return newClient, parsed, nil
}

func (it *ArtifactPublisher) Publish(fullpath, relativepath string, details os.FileInfo) {
	common.Debug("- publishing %s", relativepath)
	size, ok := pathlib.Size(fullpath)
	if !ok {
		it.ErrorCount += 1
		return //errors.New(fmt.Sprintf("Could not publish file %v, reason: could not determine size!", fullpath))
	}
	client, url, err := it.NewClient(it.ArtifactPostURL)
	if err != nil {
		it.ErrorCount += 1
		common.Error("Assistant", err)
		return //err
	}
	basename := filepath.Base(fullpath)
	request := client.NewRequest(url.RequestURI())
	request.Headers[contentType] = applicationJson
	data := make(Token)
	data["fileName"] = basename
	data["fileSize"] = fmt.Sprintf("%d", size)
	body, err := data.AsJson()
	if err != nil {
		it.ErrorCount += 1
		common.Error("Assistant", err)
		return //err
	}
	request.Body = strings.NewReader(body)
	response := client.Post(request)
	if response.Err != nil {
		it.ErrorCount += 1
		common.Error("Assistant", response.Err)
		return //err
	}
	if response.Status < 200 || 299 < response.Status {
		common.Log("ERR: status code %v", response.Status)
		return //err
	}
	var outcome awsWrapper
	err = json.Unmarshal(response.Body, &outcome)
	if err != nil {
		it.ErrorCount += 1
		common.Error("Assistant", err)
		return //err
	}
	if outcome.Response == nil {
		it.ErrorCount += 1
		common.Log("ERR: did not get correct response in reply from cloud.")
		return //err
	}
	if outcome.Response.PostInfo == nil {
		it.ErrorCount += 1
		common.Log("ERR: did not get correct response postinfo in reply from cloud.")
		return //err
	}
	err = multipartUpload(outcome.Response.PostInfo.Url, outcome.Response.PostInfo.Fields, basename, fullpath)
	if err != nil {
		it.ErrorCount += 1
		common.Error("Assistant/Last", err)
	}
}

func multipartUpload(url string, fields map[string]string, basename, fullpath string) error {
	buffer := new(bytes.Buffer)
	many := multipart.NewWriter(buffer)

	defer many.Close()

	for key, value := range fields {
		os.Stdout.Sync()
		many.WriteField(key, value)
	}
	sink, err := many.CreateFormFile("file", basename)
	if err != nil {
		return err
	}
	source, err := os.Open(fullpath)
	if err != nil {
		return err
	}
	defer source.Close()
	_, err = io.Copy(sink, source)
	if err != nil {
		return err
	}
	err = many.Close()
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, buffer)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", many.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Warning: status: %d reason: %s", response.StatusCode, IoAsString(response.Body)))
	}
	return nil
}

func IoAsString(source io.Reader) string {
	body, err := ioutil.ReadAll(source)
	if err != nil {
		return ""
	}
	return string(body)
}

func AssistantTreeCommand(client cloud.Client, account *account, workspace string) (*WorkspaceTreeData, error) {
	response, err := WorkspaceTreeCommandRequest(client, account, workspace)
	if err != nil {
		return nil, err
	}
	treedata := new(WorkspaceTreeData)
	err = json.Unmarshal(response.Body, &treedata)
	if err != nil {
		return nil, err
	}
	return treedata, nil
}

func ListAssistantsCommand(client cloud.Client, account *account, workspaceId string) ([]Token, error) {
	credentials, err := summonAssistantToken(client, account, workspaceId)
	if err != nil {
		return nil, err
	}
	request := client.NewRequest(fmt.Sprintf(listAssistantsApi, workspaceId))
	request.Headers[authorization] = WorkspaceToken(credentials)
	response := client.Get(request)
	if response.Status != 200 {
		return nil, errors.New(fmt.Sprintf("%d: %s", response.Status, response.Body))
	}
	tokens := make([]Token, 100)
	err = json.Unmarshal(response.Body, &tokens)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func BackgroundAssistantHeartbeat(cancel chan bool, client cloud.Client, account *account, workspaceId, assistantId, runId string) {
	var counter = 0
	for {
		select {
		case _ = <-cancel:
			common.Trace("Stopping assistant heartbeat.")
			return
		case <-time.After(60 * time.Second):
			counter += 1
			common.Trace("Sending assistant heartbeat #%d.", counter)
			go BeatAssistantRun(client, account, workspaceId, assistantId, runId, counter)
		}
	}
}

func BeatAssistantRun(client cloud.Client, account *account, workspaceId, assistantId, runId string, beat int) error {
	credentials, err := summonAssistantToken(client, account, workspaceId)
	if err != nil {
		return err
	}
	token := make(Token)
	token["seq"] = beat
	request := client.NewRequest(fmt.Sprintf(beatAssistantApi, workspaceId, assistantId, runId))
	request.Headers[authorization] = WorkspaceToken(credentials)
	blob, err := json.Marshal(token)
	if err == nil {
		request.Body = bytes.NewReader(blob)
	}
	response := client.Post(request)
	if response.Status != 200 {
		return errors.New(fmt.Sprintf("%d: %s", response.Status, response.Body))
	}
	return nil
}

func StopAssistantRun(client cloud.Client, account *account, workspaceId, assistantId, runId, status, reason string) error {
	credentials, err := summonAssistantToken(client, account, workspaceId)
	if err != nil {
		return err
	}
	token := make(Token)
	token["result"] = status
	token["error"] = reason
	request := client.NewRequest(fmt.Sprintf(stopAssistantApi, workspaceId, assistantId, runId))
	request.Headers[authorization] = WorkspaceToken(credentials)
	blob, err := json.Marshal(token)
	if err == nil {
		request.Body = bytes.NewReader(blob)
	}
	response := client.Put(request)
	if response.Status != 200 {
		return errors.New(fmt.Sprintf("%d: %s", response.Status, response.Body))
	}
	return nil
}

func StartAssistantRun(client cloud.Client, account *account, workspaceId, assistantId string) (*AssistantRobot, error) {
	credentials, err := summonAssistantToken(client, account, workspaceId)
	if err != nil {
		return nil, err
	}
	key, err := GenerateEphemeralKey()
	if err != nil {
		return nil, err
	}
	request := client.NewRequest(fmt.Sprintf(startAssistantApi, workspaceId, assistantId))
	request.Headers[authorization] = WorkspaceToken(credentials)
	request.Body, err = key.RequestBody(nil)
	if err != nil {
		return nil, err
	}
	response := client.Post(request)
	if response.Status != 200 {
		return nil, errors.New(fmt.Sprintf("%d: %s", response.Status, response.Body))
	}
	plaintext, err := key.Decode(response.Body)
	if err != nil {
		return nil, err
	}
	details := make(Token)
	err = json.Unmarshal(plaintext, &details)
	if err != nil {
		return nil, err
	}
	assistant := AssistantRobot{
		WorkspaceId: workspaceId,
		AssistantId: assistantId,
		Environment: make(map[string]string),
		Details:     details,
	}
	assistant.extractEnvironment(details)
	runId, ok := pickString(details, "id")
	if !ok {
		return nil, errors.New("Incorrect run-id. Cannot run without it.")
	}
	assistant.RunId = runId

	blob, ok := details["config"]
	if !ok {
		return nil, errors.New("Missing robot configuration. Cannot run without them.")
	}
	assistantConfig, ok := blob.(map[string]interface{})
	if !ok {
		return nil, errors.New("Incorrect robot configuration. Cannot run without them.")
	}
	assistant.Config = assistantConfig

	artifactURL, ok := pickString(assistantConfig, "artifactPostURL")
	if ok {
		assistant.ArtifactURL = artifactURL
	}

	robotblob, ok := details["robot"]
	if !ok {
		return nil, errors.New("Missing robot details. Cannot run without them.")
	}
	robot, ok := robotblob.(map[string]interface{})
	if !ok {
		return nil, errors.New("Incorrect robot details. Cannot run without them.")
	}
	robotId, ok := pickString(robot, "id")
	if !ok {
		return nil, errors.New("Missing robot identity in details. Cannot run without them.")
	}
	assistant.RobotId = robotId
	taskName, ok := pickString(robot, "task")
	if !ok {
		return nil, errors.New("Missing task name details. Cannot run without them.")
	}
	assistant.TaskName = taskName
	digest, _ := pickString(robot, "sha256")
	zipfile, err := SummonRobotZipfile(client, account, workspaceId, robotId, digest)
	if err != nil {
		return nil, err
	}
	assistant.Zipfile = zipfile
	return &assistant, nil
}

func pickString(from map[string]interface{}, key string) (string, bool) {
	value, ok := from[key]
	if !ok {
		return "", false
	}
	result, ok := value.(string)
	if !ok {
		return "", false
	}
	return result, true
}
