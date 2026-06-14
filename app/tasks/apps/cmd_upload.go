package apps

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
)

type UploadCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.UploadCmd]
}

func (task *UploadCmdTask) MainCmd() error {
	logger.Process("Uploading file " + task.CmdArgs.PackageFile)
	blobToken, err := task.UploadFile(task.CmdArgs.PackageFile)
	if err != nil {
		return err
	}

	requestPacket := task.CreateDefPacket()
	requestPacket.TaskType = &proto.RequestPacket_AppsTaskType{AppsTaskType: proto.TaskAppsTypes_TYPE_APP_PKG_PARSE}
	requestPacket.ExtraData = &blobToken

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if len(response.WorkerResponse) < 1 {
		if response.ErrorCause != "" {
			return errors.New(response.ErrorCause)
		}
		return fmt.Errorf("package is corrupted or missing metadata")
	}

	packageMetadata := response.GetWorkerResponse()[0].GetWorkerInfo().GetInstalledApp()[0]
	logger.Info("Package identified: \"%s\" version %s", packageMetadata.GetApp().AppName, packageMetadata.GetCurrentVersion().GetVersionName())

	requestPacket.TaskType = &proto.RequestPacket_AppsTaskType{AppsTaskType: proto.TaskAppsTypes_TYPE_APP_PKG_CREATE}
	response, err = task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if len(response.WorkerResponse) < 1 {
		if response.GetExtraData() != "" {
			return errors.New(response.GetExtraData())
		}
		return fmt.Errorf("error occurred while creating worker target package on server-side")
	}

	var packageArch string
	for _, workers := range response.GetWorkerResponse() {
		packageArch += workers.GetWorkerInfo().GetCpuArch().String() + ", "
	}

	packageArch = strings.TrimSuffix(packageArch, ", ")
	logger.Info("Supported architectures: %s", packageArch)
	logger.Info("Successfully registered at server!")
	logger.Tip("Tip: Check available packages using \"staploy-cli apps -n %s\"", packageMetadata.GetApp().GetAppName())

	return nil
}

func (task *UploadCmdTask) UploadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(consts.BLOB_FIELD_PACKAGE, filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	targetURL := task.GetServerAddr()
	req, err := http.NewRequest("POST", targetURL, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set(consts.BLOB_REQ_TYPE, consts.BLOB_REQ_TYPE_UPLOAD)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: cmds.SkipValidation,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	responseData := &proto.ResponsePacket{}
	err = protojson.Unmarshal(respBytes, responseData)
	if err != nil {
		return "", err
	}
	return responseData.GetExtraData(), nil
}
