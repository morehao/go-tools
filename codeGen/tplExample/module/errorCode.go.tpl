package errorCode

import "github.com/morehao/go-tools/gerror"

var {{.StructName}}CreateErr = gerror.Error{
	Code: 50000,
	Msg:  "创建失败",
}
