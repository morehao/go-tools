package autoCode

import (
	"gorm.io/gorm"
)

type AutoCode interface {
	GetTemplateParam() (*TemplateParams, error)
	CreateFile(param *CreateFileParam) error
}

type Cfg struct {
	TableName   string // 表名
	PackageName string // 包名
	TplDir      string // 模板目录
	RootDir     string // 生成文件的根目录
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
