package ansible

import (
	"fmt"
)

// Response represent response printed to stdout and which will be processed by
// Ansible
type Response struct {
	Msg     string `json:"msg"`
	Changed bool   `json:"changed"`
	Failed  bool   `json:"failed"`
}

// FailResponse create fail response with message
func FailResponse(msg string) Response {
	return Response{
		Msg:    msg,
		Failed: true,
	}
}

// FailResponsef do same as FailResponse but allow formats message
func FailResponsef(format string, a ...interface{}) Response {
	return FailResponse(fmt.Sprintf(format, a...))
}

// ErrorResponse do same as FailResponse but accept error instance instead message
func ErrorResponse(err error) Response {
	return FailResponse(err.Error())
}
