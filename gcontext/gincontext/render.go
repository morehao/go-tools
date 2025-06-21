package gincontext

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morehao/golib/gcontext"
	"github.com/morehao/golib/gerror"
)

type DtoRender struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(ctx *gin.Context, data any) {
	r := gcontext.NewResponseRender()
	r.SetCode(0)
	r.SetMsg("success")
	r.SetData(data)
	ctx.JSON(http.StatusOK, r)
}

func SuccessWithFormat(ctx *gin.Context, data any) {
	r := gcontext.NewResponseRender()
	r.SetCode(0)
	r.SetMsg("success")
	r.SetDataWithFormat(data)
	ctx.JSON(http.StatusOK, r)
}

func Fail(ctx *gin.Context, err error) {
	r := gcontext.NewResponseRender()

	var code int
	var msg string

	var gErr gerror.Error
	if errors.As(err, &gErr) {
		code = gErr.Code
		msg = gErr.Msg
	} else {
		code = -1
		msg = cause(err).Error()
	}

	r.SetCode(code)
	r.SetMsg(msg)
	r.SetData(gin.H{})
	ctx.JSON(http.StatusOK, r)
}

func Abort(ctx *gin.Context, err error) {
	r := gcontext.NewResponseRender()

	var gErr gerror.Error
	if errors.As(err, &gErr) {
		r.SetCode(gErr.Code)
		r.SetMsg(gErr.Msg)
		r.SetData(gin.H{})
	} else {
		r.SetCode(-1)
		r.SetMsg(cause(err).Error())
		r.SetData(gin.H{})
	}
	ctx.AbortWithStatusJSON(http.StatusOK, r)
}

func cause(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}
