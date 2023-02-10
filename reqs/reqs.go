package reqs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func Request(url, method, token string, data map[string]interface{}) (error, map[string]interface{}) {
	//Method: "OPTIONS" | "GET" | "HEAD" | "POST" | "PUT" | "DELETE" | "TRACE" | "CONNECT"
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err, nil
	}
	request, err := http.NewRequest(method, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return err, nil
	}
	// 设置请求头
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	if len(token) > 0 {
		request.Header.Add("token", token)
	}
	// 设置5秒的超时时间
	client.Timeout = 5 * time.Second

	// 发起请求
	resp, err := client.Do(request)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err, nil
	}
	return nil, res
}
