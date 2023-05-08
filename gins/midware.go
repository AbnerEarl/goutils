package gins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var WhitelistAPI = map[string]bool{
	"/api/v1/account/login":    true,
	"/api/v1/account/vericode": true,
	"/ping":                    true,
	"/favicon.ico":             true,
}

func SendResponse(c *Context, err error, data interface{}) {
	code, message := DecodeErr(err)
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func sendResponse(c *gin.Context, err error, data interface{}) {
	code, message := DecodeErr(err)
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Validate(checkToken func(token, apiPath string, c *gin.Context) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !WhitelistAPI[c.Request.URL.Path] {
			//token校验
			err := jwtVerify(c, checkToken)
			if err != nil {
				c.Abort()
				sendResponse(c, err, nil)
				return
			}
		}
		c.Next()
	}
}

func Args(checkParamMethods ...func(c *gin.Context) (bool, error)) gin.HandlerFunc {
	return func(c *gin.Context) {

		var requestParams = make(map[string]interface{})
		form, err := c.MultipartForm()
		if err == nil {
			file := form.File
			for k, v := range file {
				var names []string
				for _, f := range v {
					names = append(names, f.Filename)
				}
				c.Set(k, names)
				requestParams[k] = names
			}
		}
		bys, err := ioutil.ReadAll(c.Request.Body)
		if err == nil && bys != nil {
			maps := make(map[string]interface{})
			err = json.Unmarshal(bys, &maps)
			if err == nil {
				for k, v := range maps {
					c.Set(k, v)
					requestParams[k] = v
				}
			} else {
				params := strings.Split(string(bys), "&")
				for _, param := range params {
					if strings.Contains(param, "=") {
						arr := strings.Split(param, "=")
						key, _ := url.QueryUnescape(arr[0])
						val, _ := url.QueryUnescape(arr[1])
						c.Set(key, val)
						requestParams[key] = val
					}
				}
			}
		}
		query := c.Request.Form.Encode()
		if len(query) > 1 {
			params := strings.Split(query, "&")
			for _, param := range params {
				if len(param) > 1 && strings.Contains(param, "=") {
					arr := strings.Split(param, "=")
					key, _ := url.QueryUnescape(arr[0])
					val, _ := url.QueryUnescape(arr[1])
					c.Set(key, val)
					requestParams[key] = val
				}
			}
		}
		pathArr := c.Params
		for i := 0; i < len(pathArr); i++ {
			c.Set(pathArr[i].Key, pathArr[i].Value)
			requestParams[pathArr[i].Key] = pathArr[i].Value
		}
		c.Request.ParseForm()
		for k, v := range c.Request.PostForm {
			val := strings.Join(v, ",")
			c.Set(k, val)
			requestParams[k] = val
		}

		for _, method := range checkParamMethods {
			ck, err := method(c)
			if !ck && err != nil {
				c.Abort()
				sendResponse(c, err, nil)
				return
			}
		}

		c.Set("request_params", requestParams)
		c.Next()
	}
}

func CheckPageParam(c *gin.Context) (bool, error) {
	data, exists := c.Get("page_no")
	if exists && data != nil {
		result, e := strconv.ParseUint(fmt.Sprint(data), 10, 64)
		if e == nil {
			if result < 1 {
				return false, ErrPageParam
			}
		} else {
			return false, ErrPageParam
		}
	}
	data, exists = c.Get("page_size")
	if exists && data != nil {
		result, e := strconv.ParseUint(fmt.Sprint(data), 10, 64)
		if e == nil {
			if result < 1 || result > 1000 {
				return false, ErrPageParam
			}
		} else {
			return false, ErrPageParam
		}
	}
	return true, nil
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			//接收客户端发送的origin
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session")
			//允许浏览器（客户端）可以解析的头部
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie
			c.Header("Access-Control-Allow-Credentials", "true")
			//设置content类型
			//c.Header("Content-Type", "application/json")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// 验证token
func jwtVerify(c *gin.Context, checkToken func(token, apiPath string, c *gin.Context) (string, error)) error {
	token1 := c.GetHeader("token")
	tk, _ := c.Get("token")
	token2 := fmt.Sprint(tk)
	if len(token1) < 1 && len(token2) < 1 {
		return ErrTokenInvalid
	}
	var token string
	if len(token1) > 0 {
		token = token1
	} else {
		token = token2
	}

	//验证token，并存储在请求中
	apiPath := strings.Split(c.FullPath(), "/:")[0]
	token, err := checkToken(token, apiPath, c)
	if err != nil {
		return err
	}
	if len(token) > 0 {
		c.Header("token", token)
	}
	return nil
}

func LogAop(dealLogInfo func(data LogData)) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		blw := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		endTime := time.Now()

		var logInfo = make(map[string]interface{})
		logInfo["method"] = c.Request.Method
		logInfo["execute_time"] = time.Now()
		logInfo["content_length"] = c.Request.ContentLength
		logInfo["content_type"] = c.ContentType()
		logInfo["cost_time"] = endTime.Sub(startTime).Milliseconds()
		logInfo["request_url"] = c.Request.RequestURI
		logInfo["status_code"] = c.Writer.Status()
		logInfo["request_host"] = c.Request.Host
		logInfo["user_agent"] = c.Request.UserAgent()
		ip := c.Request.Header.Get("X-Real-IP")
		logInfo["remote_ip"] = ip
		logInfo["remote_addr"] = c.Request.RemoteAddr
		apiPath := strings.Split(c.FullPath(), "/:")[0]
		logInfo["api_path"] = apiPath
		logInfo["referer"] = c.Request.Referer()
		respData := blw.body.String()
		logInfo["response_data"] = respData
		requestParams, _ := c.Get("request_params")
		params, _ := json.Marshal(requestParams)
		logInfo["request_params"] = string(params)
		reqToken := c.GetHeader("token")
		logInfo["request_token"] = reqToken

		logData := LogData{
			LogInfo:       logInfo,
			RequestParams: requestParams.(map[string]interface{}),
			RequestToken:  reqToken,
			ResponseData:  respData,
		}
		dealLogInfo(logData)
	}
}

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GetString(s interface{}) string {
	if s != nil {
		return fmt.Sprint(s)
	} else {
		return ""
	}
}

func GetStatus(s interface{}) string {
	result, e := strconv.Atoi(fmt.Sprint(s))
	if e == nil {
		if result == 0 {
			return "停用"
		} else {
			return "启用"
		}
	}
	return "未知"
}

func GetAnyString(desc string, keys ...interface{}) string {
	var result []string
	for _, key := range keys {
		if key != nil {
			result = append(result, fmt.Sprint(key))
		}
	}
	if len(result) > 0 {
		return fmt.Sprintf("，%s%s", desc, strings.Join(result, "、"))
	} else {
		return ""
	}
}

func GetKeyword(keys ...interface{}) string {
	var result []string
	for _, key := range keys {
		if key != nil {
			result = append(result, fmt.Sprint(key))
		}
	}
	if len(result) > 0 {
		return fmt.Sprintf("，关键词为：%s", strings.Join(result, "、"))
	} else {
		return ""
	}
}

func GetJoinString(ss ...interface{}) string {
	var result []string
	for _, s := range ss {
		if s != nil {
			result = append(result, fmt.Sprint(s))
		}

	}
	if len(result) > 0 {
		return strings.Join(result, "、")
	} else {
		return "空"
	}
}
