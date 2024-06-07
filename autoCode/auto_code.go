package autoCode

import (
	"gorm.io/gorm"
)

type ModuleCfg struct {
	PackageName   string            // 包名
	TableName     string            // 表名
	ColumnTypeMap map[string]string // 字段类型映射
	TplDir        string            // 模板目录
	RootDir       string            // 生成文件的根目录
}

type AutoCode interface {
	GetModuleTemplateParam(db *gorm.DB, cfg *ModuleCfg) (*TemplateParams, error)
	CreateFile(param *CreateFileParam) error
}

func NewAutoCode(dbType DbType) AutoCode {
	switch dbType {
	case DbTypeMysql:
		return &mysqlImpl{}
	default:
		panic("unsupported database type")
	}
}
