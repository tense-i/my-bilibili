package xerr

import (
	"fmt"
)

// CodeError 自定义错误类型
type CodeError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

// NewCodeError 创建错误
func NewCodeError(code uint32, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

// NewCodeErrorWithMsg 根据错误码创建错误（使用预定义消息）
func NewCodeErrorWithMsg(code uint32) error {
	return &CodeError{Code: code, Msg: MapErrMsg(code)}
}

// Error 实现 error 接口
func (e *CodeError) Error() string {
	return fmt.Sprintf("Code: %d, Msg: %s", e.Code, e.Msg)
}

// GetErrCode 获取错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.Code
}

// GetErrMsg 获取错误消息
func (e *CodeError) GetErrMsg() string {
	return e.Msg
}

// Data 错误响应数据结构
func (e *CodeError) Data() *CodeError {
	return &CodeError{
		Code: e.Code,
		Msg:  e.Msg,
	}
}
