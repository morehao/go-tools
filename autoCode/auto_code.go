package autoCode

import (
	"gorm.io/gorm"
)

type AutoCode interface {
	Generate() error
}

type Cfg struct {
	TableName   string
	PackageName string
	PrefixName  string
	TplDir      string
	RootDir     string
}

func NewAutoCode(db *gorm.DB, cfg *Cfg) AutoCode {
	return &mysqlImpl{
		db:  db,
		cfg: cfg,
	}
}
