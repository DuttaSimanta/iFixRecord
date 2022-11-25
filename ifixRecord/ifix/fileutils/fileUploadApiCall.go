package fileutils

import (
	"bytes"
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func FileUploadAPICall(ClientID int64, orgID int64, url string, path string) (string, string, error) {
	logger.Log.Println("FileUploadAPICall")
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	var originalFilename string
	var newFileName string
	file, err := os.Open(path)
	if err != nil {
		logger.Log.Println("Unable to Open File======>", err)
		return originalFilename, newFileName, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("myFile", filepath.Base(path))
	if err != nil {
		logger.Log.Println("Unable to write File in stream======>", err)
		return originalFilename, newFileName, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		logger.Log.Println("Unable to copy File in write buffer======>", err)
		return originalFilename, newFileName, err
	}
	fw, err := writer.CreateFormField("clientid")
	if err != nil {
		logger.Log.Println(err)
		return originalFilename, newFileName, err
	}
	_, err = io.Copy(fw, strings.NewReader(strconv.FormatInt(ClientID, 10)))
	if err != nil {
		logger.Log.Println(err)
		return originalFilename, newFileName, err
	}
	fw, err = writer.CreateFormField("mstorgnhirarchyid")
	if err != nil {
		logger.Log.Println(err)
		return originalFilename, newFileName, err
	}
	_, err = io.Copy(fw, strings.NewReader(strconv.FormatInt(orgID, 10)))
	if err != nil {
		logger.Log.Println(err)
		return originalFilename, newFileName, err
	}
	err = writer.Close()
	if err != nil {
		logger.Log.Println(err)
		return originalFilename, newFileName, err
	}
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := client.Do(req)
	var fileuploadResponse entities.FileuploadResponse
	respBody, err := ioutil.ReadAll(resp.Body)

	logger.Log.Println("respbody===>", string(respBody))
	if err != nil {
		logger.Log.Println(err)
		return originalFilename, newFileName, err
	}
	err = json.Unmarshal(respBody, &fileuploadResponse)
	logger.Log.Println("fileuploadResponse========>", fileuploadResponse)
	if err != nil {
		logger.Log.Println(err)
		return originalFilename, newFileName, err
	}
	//logger.Log.Println("Response====>", resp.Body)
	originalFilename = fileuploadResponse.Details.Originalfile
	newFileName = fileuploadResponse.Details.Filename
	return originalFilename, newFileName, nil
}
