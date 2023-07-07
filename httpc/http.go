/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2023/7/6 4:38 PM
 * @desc: about the role of class.
 */

package httpc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/YouAreOnlyOne/goutils/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Request(url, method string, params map[string]interface{}, headers map[string]string, timeout time.Duration) (map[string]interface{}, error) {
	//Method: "OPTIONS" | "GET" | "HEAD" | "POST" | "PUT" | "DELETE" | "TRACE" | "CONNECT"
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	} else {
		request.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	if timeout > 0 {
		client.Timeout = timeout * time.Second
	} else {
		client.Timeout = 5 * time.Second
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(body.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DownLoadFile(url, method, filename string, params map[string]interface{}, headers map[string]string) (string, error) {
	//Method: "OPTIONS" | "GET" | "HEAD" | "POST" | "PUT" | "DELETE" | "TRACE" | "CONNECT"
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest(method, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return "", err
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	} else {
		request.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	uid, _ := uuid.NewV4()
	key := uid.String()
	key = strings.ReplaceAll(key, "-", "")
	desPath := fmt.Sprintf("/tmp/%s/", key)
	os.MkdirAll(desPath, os.ModePerm)
	filePath := fmt.Sprintf("%s%s", desPath, filename)
	out, err := os.Create(filePath)
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return filePath, nil
}

func UpLoadFileBinary(url, method, filePath string, headers map[string]string) (map[string]interface{}, error) {
	//Method: "OPTIONS" | "GET" | "HEAD" | "POST" | "PUT" | "DELETE" | "TRACE" | "CONNECT"
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	request, err := http.NewRequest(method, url, file)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	} else {
		request.Header.Add("Content-Type", "binary/octet-stream")
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(body.Bytes(), &res)
	return res, err
}

func UpLoadFileWriter(url, method, filePath string, headers map[string]string) (map[string]interface{}, error) {
	//Method: "OPTIONS" | "GET" | "HEAD" | "POST" | "PUT" | "DELETE" | "TRACE" | "CONNECT"
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	defer bodyWriter.Close()
	fileWriter, _ := bodyWriter.CreateFormFile("files", filepath.Base(filePath))
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	io.Copy(fileWriter, file)
	request, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	} else {
		request.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(body.Bytes(), &res)
	return res, err
}

func UpLoadFilesWriter(url, method string, filePaths []string, headers map[string]string) (map[string]interface{}, error) {
	//Method: "OPTIONS" | "GET" | "HEAD" | "POST" | "PUT" | "DELETE" | "TRACE" | "CONNECT"
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	defer bodyWriter.Close()
	for _, filePath := range filePaths {
		fileWriter, _ := bodyWriter.CreateFormFile("files", filepath.Base(filePath))
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		io.Copy(fileWriter, file)
		file.Close()
	}
	request, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	} else {
		request.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(body.Bytes(), &res)
	return res, err
}

func UpLoadFileForm(url, method, filePath string, params map[string]interface{}, headers map[string]string) (map[string]interface{}, error) {
	//Method: "OPTIONS" | "GET" | "HEAD" | "POST" | "PUT" | "DELETE" | "TRACE" | "CONNECT"
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	defer bodyWriter.Close()
	fileWriter, _ := bodyWriter.CreateFormFile("files", filepath.Base(filePath))
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	io.Copy(fileWriter, file)
	for key, value := range params {
		bodyWriter.WriteField(key, fmt.Sprint(value))
	}
	request, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	} else {
		request.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(body.Bytes(), &res)
	return res, err
}

func UpLoadFilesForm(url, method string, filePaths []string, params map[string]interface{}, headers map[string]string) (map[string]interface{}, error) {
	//Method: "OPTIONS" | "GET" | "HEAD" | "POST" | "PUT" | "DELETE" | "TRACE" | "CONNECT"
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	defer bodyWriter.Close()
	for _, filePath := range filePaths {
		fileWriter, _ := bodyWriter.CreateFormFile("files", filepath.Base(filePath))
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		io.Copy(fileWriter, file)
		file.Close()
	}
	for key, value := range params {
		bodyWriter.WriteField(key, fmt.Sprint(value))
	}
	request, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	} else {
		request.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(body.Bytes(), &res)
	return res, err
}
