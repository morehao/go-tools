package ctr{{.PackagePascalName}}

import (
	"{{.ProjectRootDir}}/internal/{{.ServiceName}}/dto/dto{{.PackagePascalName}}"
	"{{.ProjectRootDir}}/internal/{{.ServiceName}}/service/svc{{.PackagePascalName}}"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gcontext/ginRender"
)

type {{.ReceiverTypePascalName}}Ctr interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
	Detail(c *gin.Context)
	PageList(c *gin.Context)
}

type {{.ReceiverTypeName}}Ctr struct {
	{{.ReceiverTypeName}}Svc svc{{.PackagePascalName}}.{{.ReceiverTypePascalName}}Svc
}

var _ {{.ReceiverTypePascalName}}Ctr = (*{{.ReceiverTypeName}}Ctr)(nil)

func New{{.ReceiverTypePascalName}}Ctr() {{.ReceiverTypePascalName}}Ctr {
	return &{{.ReceiverTypeName}}Ctr{
		{{.ReceiverTypeName}}Svc: svc{{.PackagePascalName}}.New{{.ReceiverTypePascalName}}Svc(),
	}
}


// Create 创建{{.Description}}
// @Tags {{.ApiDocTag}}
// @Summary 创建{{.Description}}
// @accept application/json
// @Produce application/json
// @Param req body dto{{.PackagePascalName}}.{{.StructName}}CreateReq true "创建{{.Description}}"
// @Success 200 {object} dto.DefaultRender{data=dto{{.PackagePascalName}}.{{.StructName}}CreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router {{.ApiPrefix}}/create [post]
func (ctr *{{.ReceiverTypeName}}Ctr) Create(c *gin.Context) {
	var req dto{{.PackagePascalName}}.{{.StructName}}CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.{{.ReceiverTypeName}}Svc.Create(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// Delete 删除{{.Description}}
// @Tags {{.ApiDocTag}}
// @Summary 删除{{.Description}}
// @accept application/json
// @Produce application/json
// @Param req body dto{{.PackagePascalName}}.{{.StructName}}DeleteReq true "删除{{.Description}}"
// @Success 200 {object} dto.DefaultRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router {{.ApiPrefix}}/delete [post]
func (ctr *{{.ReceiverTypeName}}Ctr) Delete(c *gin.Context) {
	var req dto{{.PackagePascalName}}.{{.StructName}}DeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}

	if err := ctr.{{.ReceiverTypeName}}Svc.Delete(c, &req); err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, "删除成功")
	}
}

// Update 修改{{.Description}}
// @Tags {{.ApiDocTag}}
// @Summary 修改{{.Description}}
// @accept application/json
// @Produce application/json
// @Param req body dto{{.PackagePascalName}}.{{.StructName}}UpdateReq true "修改{{.Description}}"
// @Success 200 {object} dto.DefaultRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router {{.ApiPrefix}}/update [post]
func (ctr *{{.ReceiverTypeName}}Ctr) Update(c *gin.Context) {
	var req dto{{.PackagePascalName}}.{{.StructName}}UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	if err := ctr.{{.ReceiverTypeName}}Svc.Update(c, &req); err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, "修改成功")
	}
}

// Detail {{.Description}}详情
// @Tags {{.ApiDocTag}}
// @Summary {{.Description}}详情
// @accept application/json
// @Produce application/json
// @Param req query dto{{.PackagePascalName}}.{{.StructName}}DetailReq true "{{.Description}}详情"
// @Success 200 {object} dto.DefaultRender{data=dto{{.PackagePascalName}}.{{.StructName}}DetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router {{.ApiPrefix}}/detail [get]
func (ctr *{{.ReceiverTypeName}}Ctr) Detail(c *gin.Context) {
	var req dto{{.PackagePascalName}}.{{.StructName}}DetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.{{.ReceiverTypeName}}Svc.Detail(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// PageList {{.Description}}列表
// @Tags {{.ApiDocTag}}
// @Summary {{.Description}}列表分页
// @accept application/json
// @Produce application/json
// @Param req query dto{{.PackagePascalName}}.{{.StructName}}PageListReq true "{{.Description}}列表"
// @Success 200 {object} dto.DefaultRender{data=dto{{.PackagePascalName}}.{{.StructName}}PageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router {{.ApiPrefix}}/pageList [get]
func (ctr *{{.ReceiverTypeName}}Ctr) PageList(c *gin.Context) {
	var req dto{{.PackagePascalName}}.{{.StructName}}PageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.{{.ReceiverTypeName}}Svc.PageList(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}
