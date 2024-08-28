package errorCode

import "github.com/morehao/go-tools/gerror"

var {{.ReceiverTypePascalName}}{{.FunctionName}}Err = gerror.Error{
	Code: 100100,
	Msg:  "{{.Description}}失败",
}
