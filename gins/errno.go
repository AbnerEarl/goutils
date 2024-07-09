package gins

import (
	"fmt"
)

var (
	OK              = &Errno{Code: 0, Message: "success", Tips: "成功"}
	ParamError      = &Errno{Code: 1, Message: "request parameter error", Tips: "请求参数错误"}
	InternalError   = &Errno{Code: 10001, Message: "internal server error", Tips: "服务器内部错误"}
	ErrTokenInvalid = &Errno{Code: 20001, Message: "the token was invalid", Tips: "Token无效"}
	ErrPageParam    = &Errno{Code: 30001, Message: "the parameter of page_no or page_size is error", Tips: "分页参数错误"}
)

type Response struct {
	*Errno
	Data interface{} `json:"data"`
}

type ThirdResponse struct {
	*Errno
	Data map[string]interface{} `json:"data"`
}

type LogData struct {
	LogInfo       map[string]interface{} `json:"log_info"`
	RequestParams map[string]interface{} `json:"request_params"`
	RequestToken  string                 `json:"request_token"`
	ResponseData  string                 `json:"response_data"`
}

type Errno struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Tips    string `json:"tips"`
	Err     string `json:"err"`
}

func NewErr(errno *Errno, err error) *Errno {
	return &Errno{Code: errno.Code, Message: errno.Message, Tips: errno.Tips, Err: err.Error()}
}

func (err *Errno) Add(message string) {
	err.Message += " " + message
}

func (err *Errno) AddAny(format string, args ...interface{}) {
	err.Message += " " + fmt.Sprintf(format, args...)
}

func (err *Errno) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, tips:%s, err: %s", err.Code, err.Message, err.Tips, err.Err)
}

func DecodeErr(err error) *Errno {
	if err == nil {
		return OK
	}

	switch typed := err.(type) {
	case *Errno:
		return typed
	default:
	}
	e := new(Errno)
	*e = *InternalError
	e.Err = err.Error()
	return e
}
