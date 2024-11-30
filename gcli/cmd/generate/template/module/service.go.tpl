package svc{{.PackagePascalName}}

import (
	"{{.ProjectRootDir}}/internal/app/dto/dto{{.PackagePascalName}}"
	"{{.ProjectRootDir}}/internal/app/model/dao{{.PackagePascalName}}"
	"{{.ProjectRootDir}}/internal/app/object/objCommon"
	"{{.ProjectRootDir}}/internal/app/object/obj{{.PackagePascalName}}"
	"{{.ProjectRootDir}}/internal/pkg/context"
	"{{.ProjectRootDir}}/internal/pkg/errorCode"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
)

type {{.ReceiverTypePascalName}}Svc interface {
	Create(c *gin.Context, req *dto{{.PackagePascalName}}.{{.StructName}}CreateReq) (*dto{{.PackagePascalName}}.{{.StructName}}CreateResp, error)
	Delete(c *gin.Context, req *dto{{.PackagePascalName}}.{{.StructName}}DeleteReq) error
	Update(c *gin.Context, req *dto{{.PackagePascalName}}.{{.StructName}}UpdateReq) error
	Detail(c *gin.Context, req *dto{{.PackagePascalName}}.{{.StructName}}DetailReq) (*dto{{.PackagePascalName}}.{{.StructName}}DetailResp, error)
	PageList(c *gin.Context, req *dto{{.PackagePascalName}}.{{.StructName}}PageListReq) (*dto{{.PackagePascalName}}.{{.StructName}}PageListResp, error)
}

type {{.ReceiverTypeName}}Svc struct {
}

var _ {{.ReceiverTypePascalName}}Svc = (*{{.ReceiverTypeName}}Svc)(nil)

func New{{.ReceiverTypePascalName}}Svc() {{.ReceiverTypePascalName}}Svc {
	return &{{.ReceiverTypeName}}Svc{}
}

// Create 创建{{.Description}}
func (svc *{{.ReceiverTypeName}}Svc) Create(c *gin.Context, req *dto{{.PackagePascalName}}.{{.StructName}}CreateReq) (*dto{{.PackagePascalName}}.{{.StructName}}CreateResp, error) {
	userId := context.GetUserID(c)
	now := time.Now()
	insertEntity := &dao{{.PackagePascalName}}.{{.StructName}}Entity{
{{- range .ModelFields}}
	{{- if .IsPrimaryKey}}
		{{- continue}}
	{{- end}}
	{{- if isSysField .FieldName}}
		{{- continue}}
	{{- end}}
	{{- if eq .FieldType "time.Time"}}
		{{.FieldName}}: time.Unix(req.{{.FieldName}}, 0),
	{{- else}}
		{{.FieldName}}: req.{{.FieldName}},
	{{- end}}
{{- end}}
		CreatedBy: userId,
		CreatedAt: now,
		UpdatedBy: userId,
		UpdatedAt: now,
	}

	if err := dao{{.PackagePascalName}}.New{{.StructName}}Dao().Insert(c, insertEntity); err != nil {
		glog.Errorf(c, "[svc{{.PackagePascalName}}.{{.StructName}}Create] dao{{.StructName}} Create fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.{{.StructName}}CreateErr
	}
	return &dto{{.PackagePascalName}}.{{.StructName}}CreateResp{
		ID: insertEntity.ID,
	}, nil
}

// Delete 删除{{.Description}}
func (svc *{{.ReceiverTypeName}}Svc) Delete(c *gin.Context, req *dto{{.PackagePascalName}}.{{.StructName}}DeleteReq) error {
	deletedBy := context.GetUserID(c)

	if err := dao{{.PackagePascalName}}.New{{.StructName}}Dao().Delete(c, req.ID, deletedBy); err != nil {
		glog.Errorf(c, "[svc{{.PackagePascalName}}.Delete] dao{{.StructName}} Delete fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return errorCode.{{.StructName}}DeleteErr
	}
	return nil
}

// Update 更新{{.Description}}
func (svc *{{.ReceiverTypeName}}Svc) Update(c *gin.Context, req *dto{{.PackagePascalName}}.{{.StructName}}UpdateReq) error {
	updateEntity := &dao{{.PackagePascalName}}.{{.StructName}}Entity{
        ID:   req.ID,
    }
    if err := dao{{.PackagePascalName}}.New{{.StructName}}Dao().Update(c, updateEntity); err != nil {
        glog.Errorf(c, "[svc{{.PackagePascalName}}.{{.StructName}}Update] dao{{.StructName}} Update fail, err:%v, req:%s", err, gutils.ToJsonString(req))
        return errorCode.{{.StructName}}UpdateErr
    }
    return nil
}

// Detail 根据id获取{{.Description}}
func (svc *{{.ReceiverTypeName}}Svc) Detail(c *gin.Context, req *dto{{.PackagePascalName}}.{{.StructName}}DetailReq) (*dto{{.PackagePascalName}}.{{.StructName}}DetailResp, error) {
	detailEntity, err := dao{{.PackagePascalName}}.New{{.StructName}}Dao().GetById(c, req.ID)
	if err != nil {
		glog.Errorf(c, "[svc{{.PackagePascalName}}.{{.StructName}}Detail] dao{{.StructName}} GetById fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.{{.StructName}}GetDetailErr
	}
    // 判断是否存在
    if detailEntity == nil || detailEntity.ID == 0 {
        return nil, errorCode.{{.StructName}}NotExistErr
    }
	Resp := &dto{{.PackagePascalName}}.{{.StructName}}DetailResp{
		ID:   detailEntity.ID,
		{{.StructName}}BaseInfo: obj{{.PackagePascalName}}.{{.StructName}}BaseInfo{
	{{- range .ModelFields}}
		{{- if .IsPrimaryKey}}
			{{- continue}}
		{{- end}}
		{{- if isSysField .FieldName}}
			{{- continue}}
		{{- end}}
		{{- if eq .FieldType "time.Time"}}
			{{.FieldName}}: detailEntity.{{.FieldName}}.Unix(),
		{{- else}}
			{{.FieldName}}: detailEntity.{{.FieldName}},
		{{- end}}
	{{- end}}
		},
		OperatorBaseInfo: objCommon.OperatorBaseInfo{
        	CreatedBy: detailEntity.CreatedBy,
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedBy: detailEntity.UpdatedBy,
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}
	return Resp, nil
}

// PageList 分页获取{{.Description}}列表
func (svc *{{.ReceiverTypeName}}Svc) PageList(c *gin.Context, req *dto{{.PackagePascalName}}.{{.StructName}}PageListReq) (*dto{{.PackagePascalName}}.{{.StructName}}PageListResp, error) {
	cond := &dao{{.PackagePascalName}}.{{.StructName}}Cond{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	dataList, total, err := dao{{.PackagePascalName}}.New{{.StructName}}Dao().GetPageListByCond(c, cond)
	if err != nil {
		glog.Errorf(c, "[svc{{.PackagePascalName}}.{{.StructName}}PageList] dao{{.StructName}} GetPageListByCond fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.{{.StructName}}GetPageListErr
	}
	list := make([]dto{{.PackagePascalName}}.{{.StructName}}PageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dto{{.PackagePascalName}}.{{.StructName}}PageListItem{
			ID:   v.ID,
			{{.StructName}}BaseInfo: obj{{.PackagePascalName}}.{{.StructName}}BaseInfo{
		{{- range .ModelFields}}
			{{- if .IsPrimaryKey}}
				{{- continue}}
			{{- end}}
			{{- if isSysField .FieldName}}
				{{- continue}}
			{{- end}}
			{{- if eq .FieldType "time.Time"}}
				{{.FieldName}}: v.{{.FieldName}}.Unix(),
			{{- else}}
				{{.FieldName}}: v.{{.FieldName}},
			{{- end}}
		{{- end}}
			},
			OperatorBaseInfo: objCommon.OperatorBaseInfo{
				CreatedBy: v.CreatedBy,
				CreatedAt: v.CreatedAt.Unix(),
				UpdatedBy: v.UpdatedBy,
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dto{{.PackagePascalName}}.{{.StructName}}PageListResp{
		List:  list,
		Total: total,
	}, nil
}


