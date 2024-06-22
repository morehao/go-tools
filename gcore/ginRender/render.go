package ginRender

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gcore"
	"github.com/morehao/go-tools/gerror"
	"github.com/pkg/errors"
	"net/http"
)

func RenderSuccess(ctx *gin.Context, data interface{}) {
	r := gcore.NewResponseRender()
	r.SetCode(0)
	r.SetMsg("success")
	r.SetData(data)
	ctx.JSON(http.StatusOK, r)
	return
}

func RenderSuccessWithFormat(ctx *gin.Context, data interface{}) {
	r := gcore.NewResponseRender()
	r.SetCode(0)
	r.SetMsg("success")
	r.SetDataWithFormat(data)
	ctx.JSON(http.StatusOK, r)
	return
}

func RenderFail(ctx *gin.Context, err error) {
	r := gcore.NewResponseRender()

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
	r := gcore.NewResponseRender()

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
