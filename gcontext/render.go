package gcontext

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gerror"
	"github.com/pkg/errors"
	"net/http"
)

// ResponseRender 返回数据格式化
type ResponseRender interface {
	SetCode(int)
	SetMsg(string)
	SetData(interface{})
	SetDataWithFormat(interface{})
}

var newRender func() ResponseRender

func RegisterRender(s func() ResponseRender) {
	newRender = s
}

func newResponseRender() ResponseRender {
	if newRender == nil {
		newRender = newDefaultRender
	}
	return newRender()
}

func RenderJson(ctx *gin.Context, code int, msg string, data interface{}) {
	r := newResponseRender()
	r.SetCode(code)
	r.SetMsg(msg)
	r.SetData(data)
	ctx.JSON(http.StatusOK, r)
	return
}

func RenderSuccess(ctx *gin.Context, data interface{}) {
	r := newResponseRender()
	r.SetCode(0)
	r.SetMsg("success")
	r.SetData(data)
	ctx.JSON(http.StatusOK, r)
	return
}

func RenderSuccessWithFormat(ctx *gin.Context, data interface{}) {
	r := newResponseRender()
	r.SetCode(0)
	r.SetMsg("success")
	r.SetDataWithFormat(data)
	ctx.JSON(http.StatusOK, r)
	return
}

func RenderFail(ctx *gin.Context, err error) {
	r := newResponseRender()

	code, msg := -1, errors.Cause(err).Error()
	switch errors.Cause(err).(type) {
	case gerror.Error:
		code = errors.Cause(err).(gerror.Error).Code
		msg = errors.Cause(err).(gerror.Error).Msg
	default:
	}

	r.SetCode(code)
	r.SetMsg(msg)
	r.SetData(gin.H{})
	ctx.JSON(http.StatusOK, r)

	return
}

func RenderAbort(ctx *gin.Context, err error) {
	r := newResponseRender()

	switch errors.Cause(err).(type) {
	case gerror.Error:
		r.SetCode(errors.Cause(err).(gerror.Error).Code)
		r.SetMsg(errors.Cause(err).(gerror.Error).Msg)
		r.SetData(gin.H{})
	default:
		r.SetCode(-1)
		r.SetMsg(errors.Cause(err).Error())
		r.SetData(gin.H{})
	}
	ctx.AbortWithStatusJSON(http.StatusOK, r)

	return
}

var newDefaultRender = func() ResponseRender {
	return &DefaultResponseRender{}
}

type DefaultResponseRender struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *DefaultResponseRender) SetCode(code int) {
	r.Code = code
}
func (r *DefaultResponseRender) SetMsg(msg string) {
	r.Msg = msg
}
func (r *DefaultResponseRender) SetData(data interface{}) {
	r.Data = data
}

func (r *DefaultResponseRender) SetDataWithFormat(data interface{}) {
	ResponseFormat(data)
	r.Data = data
}
