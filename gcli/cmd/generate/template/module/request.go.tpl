package dto{{.PackagePascalName}}

import (
	"{{.ProjectRootDir}}/internal/{{.ServiceName}}/object/objCommon"
	"{{.ProjectRootDir}}/internal/{{.ServiceName}}/object/obj{{.PackagePascalName}}"
)

type {{.StructName}}CreateReq struct {
	obj{{.PackagePascalName}}.{{.StructName}}BaseInfo
}

type {{.StructName}}UpdateReq struct {
	ID uint64 `json:"id" validate:"required" label:"数据自增id"` // 数据自增id
	obj{{.PackagePascalName}}.{{.StructName}}BaseInfo
}

type {{.StructName}}DetailReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type {{.StructName}}PageListReq struct {
	objCommon.PageQuery
}

type {{.StructName}}DeleteReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}
