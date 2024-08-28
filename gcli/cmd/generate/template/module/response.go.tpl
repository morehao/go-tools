package dto{{.PackagePascalName}}

import (
	"{{.ProjectRootDir}}/internal/{{.ServiceName}}/object/objCommon"
	"{{.ProjectRootDir}}/internal/{{.ServiceName}}/object/obj{{.PackagePascalName}}"
)

type {{.StructName}}CreateResp struct {
	ID uint64 `json:"id"` // 数据自增id
}

type {{.StructName}}DetailResp struct {
	ID        uint64 `json:"id" validate:"required"` // 数据自增id
	obj{{.PackagePascalName}}.{{.StructName}}BaseInfo
	objCommon.OperatorBaseInfo

}

type {{.StructName}}PageListItem struct {
	ID        uint64 `json:"id" validate:"required"` // 数据自增id
	obj{{.PackagePascalName}}.{{.StructName}}BaseInfo
	objCommon.OperatorBaseInfo
}

type {{.StructName}}PageListResp struct {
	List  []{{.StructName}}PageListItem `json:"list"`  // 数据列表
	Total int64          `json:"total"` // 数据总条数
}
