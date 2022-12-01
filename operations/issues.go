package operations

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/settings"
	"github.com/robocorp/rcc/xviper"
)

const (
	issueUrl = `/diagnostics-v1/issue`
)

func loadToken(reportFile string) (Token, error) {
	content, err := os.ReadFile(reportFile)
	if err != nil {
		return nil, err
	}
	token := make(Token)
	err = token.FromJson(content)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func createIssueZip(attachmentsFiles []string) (string, error) {
	zipfile := filepath.Join(common.RobocorpTemp(), "attachments.zip")
	zipper, err := newZipper(zipfile)
	if err != nil {
		return "", err
	}
	defer zipper.Close()
	for index, attachment := range attachmentsFiles {
		niceName := fmt.Sprintf("%x_%s", index+1, filepath.Base(attachment))
		zipper.Add(attachment, niceName, nil)
	}
	// getting settings.yaml is optional, it should not break issue reporting
	config, err := settings.SummonSettings()
	if err != nil {
		return zipfile, nil
	}
	blob, err := config.AsYaml()
	if err != nil {
		return zipfile, nil
	}
	niceName := fmt.Sprintf("%x_settings.yaml", len(attachmentsFiles)+1)
	zipper.AddBlob(niceName, blob)
	return zipfile, nil
}

func createDiagnosticsReport(robotfile string) (string, *common.DiagnosticStatus, error) {
	file := filepath.Join(common.RobocorpTemp(), "diagnostics.txt")
	diagnostics, err := ProduceDiagnostics(file, robotfile, false, false)
	if err != nil {
		return "", nil, err
	}
	return file, diagnostics, nil
}

func virtualName(filename string) (string, error) {
	digest, err := pathlib.Sha256(filename)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("attachments_%s.zip", digest[:16]), nil
}

func ReportIssue(email, robotFile, reportFile string, attachmentsFiles []string, dryrun bool) error {
	issueHost := settings.Global.IssuesURL()
	if len(issueHost) == 0 {
		return nil
	}
	cloud.BackgroundMetric(common.ControllerIdentity(), "rcc.submit.issue", common.Version)
	token, err := loadToken(reportFile)
	if err != nil {
		return err
	}
	diagnostics, data, err := createDiagnosticsReport(robotFile)
	if err == nil {
		attachmentsFiles = append(attachmentsFiles, diagnostics)
	}
	plan, ok := data.Details["robot-conda-plan"]
	if ok {
		attachmentsFiles = append(attachmentsFiles, plan)
	}
	attachmentsFiles = append(attachmentsFiles, reportFile)
	filename, err := createIssueZip(attachmentsFiles)
	if err != nil {
		return err
	}
	shortname, err := virtualName(filename)
	if err != nil {
		return err
	}
	installationId := xviper.TrackingIdentity()
	token["installationId"] = installationId
	token["account-email"] = email
	token["fileName"] = shortname
	token["controller"] = common.ControllerIdentity()
	_, ok = token["platform"]
	if !ok {
		token["platform"] = common.Platform()
	}
	issueReport, err := token.AsJson()
	if err != nil {
		return err
	}
	if dryrun {
		metaForm := make(Token)
		metaForm["report"] = token
		metaForm["zipfile"] = filename
		report, err := metaForm.AsJson()
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, report)
		return nil
	}
	common.Trace(issueReport)
	client, err := cloud.NewClient(issueHost)
	if err != nil {
		return err
	}
	request := client.NewRequest(issueUrl)
	request.Headers[contentType] = applicationJson
	request.Body = bytes.NewBuffer([]byte(issueReport))
	response := client.Post(request)
	json := make(Token)
	err = json.FromJson(response.Body)
	if err != nil {
		return err
	}
	postInfo, ok := json["attachmentPostInfo"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Could not get attachmentPostInfo!")
	}
	url, ok := postInfo["url"].(string)
	if !ok {
		return fmt.Errorf("Could not get URL from attachmentPostInfo!")
	}
	fields, ok := postInfo["fields"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Could not get fields from attachmentPostInfo!")
	}
	return MultipartUpload(url, toStringMap(fields), shortname, filename)
}

func toStringMap(entries map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range entries {
		text, ok := value.(string)
		if ok {
			result[key] = text
		}
	}
	return result
}
