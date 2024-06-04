package autoCode

import (
	"gorm.io/gorm"
)

type AutoCode interface {
	Generate() error
	GetTemplateParam() (*TemplateParams, error)
	CreateFile(param *CreateFileParam) error
}

type Cfg struct {
	TableName   string
	PackageName string
	TplDir      string
	RootDir     string
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
