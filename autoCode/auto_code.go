package autoCode

import (
	"gorm.io/gorm"
)

type Cfg struct {
	PackageName   string            // 包名
	TableName     string            // 表名
	ColumnTypeMap map[string]string // 字段类型映射
	TplDir        string            // 模板目录
	RootDir       string            // 生成文件的根目录
}

type AutoCode interface {
	GetTemplateParam() (*TemplateParams, error)
	CreateFile(param *CreateFileParam) error
}

func NewAutoCode(db *gorm.DB, cfg *Cfg) AutoCode {
	dbType := db.Dialector.Name()
	switch dbType {
	case dbTypeMysql:
		return &mysqlImpl{
			db:  db,
			cfg: cfg,
		}
	default:
		panic("unsupported database type")
	}
}
