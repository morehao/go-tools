package gerror

import (
	"errors"
	"fmt"
)

type Error struct {
	Code int
	Msg  string
}

type ErrorMap map[int]Error
type CodeMsgMap map[int]string

// Error 方法实现 error 接口
func (e Error) Error() string {
	return e.Msg
}

// Wrap 用于包装错误信息（使用标准库 fmt.Errorf）
func (e *Error) Wrap(err error) error {
	if err == nil {
		return nil
	}

	msg := e.Msg
	e.Msg = err.Error() // 可选：记录原始错误内容
	return fmt.Errorf("%s: %w", msg, err)
}

// Wrapf 用于格式化包装错误信息
func (e *Error) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	msg := e.Msg
	formattedMsg := fmt.Sprintf(format, args...)
	e.Msg = err.Error() // 可选：记录原始错误内容
	return fmt.Errorf("%s %s: %w", formattedMsg, msg, err)
}

// 获取错误码
func (e Error) GetCode() int {
	return e.Code
}

// 获取错误信息
func (e Error) GetMsg() string {
	return e.Msg
}

// ResetMsg 重置错误信息
func (e Error) ResetMsg(msg string) Error {
	e.Msg = msg
	return e
}

// AppendMsg 在原有错误信息上追加内容
func (e *Error) AppendMsg(args ...interface{}) {
	e.Msg += fmt.Sprint(args...)
}

// Is 判断是否为指定错误类型
func (e Error) Is(targetErr error) bool {
	return errors.Is(targetErr, e)
}

// As 判断是否为指定错误类型，并赋值给目标变量
func (e Error) As(targetErr interface{}) bool {
	return errors.As(e, targetErr)
}
