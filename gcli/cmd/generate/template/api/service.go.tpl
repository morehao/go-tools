package svc{{.PackagePascalName}}

import (
    "{{.ProjectRootDir}}/internal/{{.ServiceName}}/dto/dto{{.PackagePascalName}}"

    "github.com/gin-gonic/gin"
)

{{if not .TargetFileExist}}
type {{.ReceiverTypePascalName}}Svc interface {
    {{.FunctionName}}(c *gin.Context, req *dto{{.PackagePascalName}}.{{.FunctionName}}Req) (*dto{{.PackagePascalName}}.{{.FunctionName}}Resp, error)
}

type {{.ReceiverTypeName}}Svc struct {
}

var _ {{.ReceiverTypePascalName}}Svc = (*{{.ReceiverTypeName}}Svc)(nil)

func New{{.ReceiverTypePascalName}}Svc() {{.ReceiverTypePascalName}}Svc {
    return &{{.ReceiverTypeName}}Svc{
    }
}
{{end}}
func (svc *{{.ReceiverTypeName}}Svc) {{.FunctionName}}(c *gin.Context, req *dto{{.PackagePascalName}}.{{.FunctionName}}Req) (*dto{{.PackagePascalName}}.{{.FunctionName}}Resp, error) {
    return &dto{{.PackagePascalName}}.{{.FunctionName}}Resp{}, nil
}
