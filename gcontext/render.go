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

func newResponseRender() ResponseRender {
	if newRender == nil {
		newRender = newDefaultRender
	}
	return newRender()
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

	var code int
	var msg string

	var gErr gerror.Error
	if errors.As(err, &gErr) {
		code = gErr.Code
		msg = gErr.Msg
	} else {
		code = -1
		msg = errors.Cause(err).Error()
	}

	r.SetCode(code)
	r.SetMsg(msg)
	r.SetData(gin.H{})
	ctx.JSON(http.StatusOK, r)
	return
}

func RenderAbort(ctx *gin.Context, err error) {
	r := newResponseRender()

	var gErr gerror.Error
	if errors.As(err, &gErr) {
		r.SetCode(gErr.Code)
		r.SetMsg(gErr.Msg)
		r.SetData(gin.H{})
	} else {
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
