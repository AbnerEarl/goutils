package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
)

type BaseGinEngine struct {
	Token  string
	Router *gin.Engine
}

func (b *BaseGinEngine) Get(uri string, param map[string]string) *httptest.ResponseRecorder {
	if param != nil {
		uri = uri + MapToStr(param)
	}
	// 构造请求
	req := httptest.NewRequest("GET", uri, nil)
	req.Header.Set("token", b.Token)
	// 初始化响应
	res := httptest.NewRecorder()
	// 调用相应的handler接口
	b.Router.ServeHTTP(res, req)
	return res
}

func (b *BaseGinEngine) Delete(uri string, param map[string]string) *httptest.ResponseRecorder {
	if param != nil {
		uri = uri + MapToStr(param)
	}
	// 构造请求
	req := httptest.NewRequest("DELETE", uri, nil)
	req.Header.Set("token", b.Token)
	// 初始化响应
	res := httptest.NewRecorder()
	// 调用相应的handler接口
	b.Router.ServeHTTP(res, req)
	return res
}

func (b *BaseGinEngine) PostForm(uri string, param map[string]string) *httptest.ResponseRecorder {
	if param != nil {
		uri = uri + MapToStr(param)
	}
	req := httptest.NewRequest("POST", uri, nil)
	req.Header.Set("token", b.Token)
	// 初始化响应
	res := httptest.NewRecorder()
	// 调用相应handler接口
	b.Router.ServeHTTP(res, req)
	return res
}

func (b *BaseGinEngine) PostJson(uri string, param map[string]interface{}) *httptest.ResponseRecorder {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(param)
	// 构造请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))
	req.Header.Set("token", b.Token)
	// 初始化响应
	res := httptest.NewRecorder()
	// 调用相应的handler接口
	b.Router.ServeHTTP(res, req)
	return res
}

func (b *BaseGinEngine) PostFile(uri string, param map[string]interface{}, files map[string]string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	//defer writer.Close()

	for k, v := range files {
		file, err := os.Open(v)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		part, err := writer.CreateFormFile(k, filepath.Base(v))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		io.Copy(part, file)
	}

	for k, v := range param {
		err := writer.WriteField(k, fmt.Sprint(v))
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}
	writer.Close()
	// 构造请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, body)
	req.Header.Set("token", b.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 初始化响应
	res := httptest.NewRecorder()
	// 调用相应的handler接口
	b.Router.ServeHTTP(res, req)
	return res
}

func (b *BaseGinEngine) PutForm(uri string, param map[string]string) *httptest.ResponseRecorder {
	if param != nil {
		uri = uri + MapToStr(param)
	}
	req := httptest.NewRequest("PUT", uri, nil)
	req.Header.Set("token", b.Token)
	// 初始化响应
	res := httptest.NewRecorder()
	// 调用相应handler接口
	b.Router.ServeHTTP(res, req)
	return res
}

func (b *BaseGinEngine) PutJson(uri string, param map[string]interface{}) *httptest.ResponseRecorder {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(param)
	// 构造请求，json数据以请求body的形式传递
	req := httptest.NewRequest("PUT", uri, bytes.NewReader(jsonByte))
	req.Header.Set("token", b.Token)
	// 初始化响应
	res := httptest.NewRecorder()
	// 调用相应的handler接口
	b.Router.ServeHTTP(res, req)
	return res
}
