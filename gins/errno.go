package gins

import "fmt"

var (
	OK              = &Errno{Code: 0, Message: "OK"}
	InternalError   = &Errno{Code: 10001, Message: "Internal server error"}
	ErrTokenInvalid = &Errno{Code: 20001, Message: "The token was invalid."}
	ErrParam        = &Errno{Code: 30001, Message: "The parameter is error."}
	ErrPageParam    = &Errno{Code: 30002, Message: "The parameter of page no or page size is error"}
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Message struct {
	Style int         `json:"style"`
	Data  interface{} `json:"data"`
}

type ThirdResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type LogData struct {
	LogInfo       map[string]interface{} `json:"log_info"`
	RequestParams interface{}            `json:"request_params"`
	AccountInfo   interface{}            `json:"account_info"`
}

type Errno struct {
	Code    int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

// Err represents an error
type Err struct {
	Code    int
	Message string
	Err     error
}

func New(errno *Errno, err error) *Err {
	return &Err{Code: errno.Code, Message: errno.Message, Err: err}
}

func (err *Err) Add(message string) error {
	err.Message += " " + message
	return err
}

func (err *Err) AddAny(format string, args ...interface{}) error {
	err.Message += " " + fmt.Sprintf(format, args...)
	return err
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	return InternalError.Code, err.Error()
}
