package ginRender

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gcontext"
	"github.com/morehao/go-tools/gerror"
	"github.com/pkg/errors"
)

func Success(ctx *gin.Context, data any) {
	r := gcontext.NewResponseRender()
	r.SetCode(0)
	r.SetMsg("success")
	r.SetData(data)
	ctx.JSON(http.StatusOK, r)
	return
}

func SuccessWithFormat(ctx *gin.Context, data any) {
	r := gcontext.NewResponseRender()
	r.SetCode(0)
	r.SetMsg("success")
	r.SetDataWithFormat(data)
	ctx.JSON(http.StatusOK, r)
	return
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
		msg = errors.Cause(err).Error()
	}

	r.SetCode(code)
	r.SetMsg(msg)
	r.SetData(gin.H{})
	ctx.JSON(http.StatusOK, r)
	return
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
		r.SetMsg(errors.Cause(err).Error())
		r.SetData(gin.H{})
	}
	ctx.AbortWithStatusJSON(http.StatusOK, r)

	return
}
