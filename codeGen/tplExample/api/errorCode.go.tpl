package errorCode

import "github.com/morehao/go-tools/gerror"

var {{.FunctionName}}Err = gerror.Error{
	Code: 50000,
	Msg:  "处理失败",
}
