package codeGen

import (
	"gorm.io/gorm"
)

type ModuleCfg struct {
	PackageName   string            `validate:"required"` // 包名
	TableName     string            `validate:"required"` // 表名
	ColumnTypeMap map[string]string // 字段类型映射
	TplDir        string            `validate:"required"` // 模板目录
	RootDir       string            `validate:"required"` // 生成文件的根目录
}

type ControllerCfg struct {
	PackageName    string // 包名
	TargetFilename string // 目标文件名
	FunctionName   string // 函数名
	TplDir         string // 模板目录
	RootDir        string // 生成文件的根目录
}

type Generator interface {
	GetModuleTemplateParam(db *gorm.DB, cfg *ModuleCfg) (*ModelTemplateParamsRes, error)
	GetControllerTemplateParam(cfg *ControllerCfg) (*ControllerTemplateParams, error)
	Gen(param *GenParam) error
}

func NewGenerator() Generator {
	return &generatorImpl{}
}
