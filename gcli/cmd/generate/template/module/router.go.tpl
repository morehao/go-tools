package router

import (
	"{{.ProjectRootDir}}/internal/{{.ServiceName}}/controller/ctr{{.PackagePascalName}}"

	"github.com/gin-gonic/gin"
)

// {{.ReceiverTypeName}}Router 初始化{{.Description}}路由信息
func {{.ReceiverTypeName}}Router(routerGroup *gin.RouterGroup) {
	{{.ReceiverTypeName}}Ctr := ctr{{.PackagePascalName}}.New{{.ReceiverTypePascalName}}Ctr()
	{{.ReceiverTypeName}}Group := routerGroup.Group("{{.ApiGroup}}")
	{
		{{.ReceiverTypeName}}Group.POST("create", {{.ReceiverTypeName}}Ctr.Create)   // 新建{{.Description}}
		{{.ReceiverTypeName}}Group.POST("delete", {{.ReceiverTypeName}}Ctr.Delete)   // 删除{{.Description}}
		{{.ReceiverTypeName}}Group.POST("update", {{.ReceiverTypeName}}Ctr.Update)   // 更新{{.Description}}
		{{.ReceiverTypeName}}Group.GET("detail", {{.ReceiverTypeName}}Ctr.Detail)    // 根据ID获取{{.Description}}
        {{.ReceiverTypeName}}Group.GET("pageList", {{.ReceiverTypeName}}Ctr.PageList)  // 获取{{.Description}}列表
	}
}
