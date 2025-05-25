package gerror

import (
	"fmt"

	"github.com/pkg/errors"
)

type Error struct {
	Code int
	Msg  string
}

type ErrorMap map[int]Error
type CodeMsgMap map[int]string

// Error 方法实现了error接口，返回错误信息
func (e Error) Error() string {
	return e.Msg
}

// Wrap 用于包装错误信息
func (e Error) Wrap(err error) error {
	if err == nil {
		return nil
	}

	// 保存原始错误信息
	msg := e.Msg
	// 更新错误信息，方便后续使用
	e.Msg = err.Error()
	return errors.Wrap(err, msg)
}

func (e Error) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	// 保存原始错误信息
	msg := e.Msg
	formattedMsg := fmt.Sprintf(format, args...)
	e.Msg = err.Error()
	return errors.Wrap(err, fmt.Sprintf("%s %s", formattedMsg, msg))
}

func (e Error) GetCode() int {
	return e.Code
}

func (e Error) GetMsg() string {
	return e.Msg
}

// ResetMsg 重置错误信息
func (e Error) ResetMsg(msg string) Error {
	e.Msg = msg
	return e
}

// AppendMsg 在既有错误信息的基础上追加新的错误信息
func (e Error) AppendMsg(args ...interface{}) {
	e.Msg += fmt.Sprint(args...)
}

// Is 判断是否为指定错误类型
func (e Error) Is(targetErr error) bool {
	return errors.Is(targetErr, e)
}

// As 判断是否为指定错误类型，并将错误信息赋值给目标变量
func (e Error) As(targetErr interface{}) bool {
	return errors.As(e, targetErr)
}
