package codeGen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestCreateModuleFile(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	db, openErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.Nil(t, openErr)
	// 获取当前的运行路径
	workDir, getErr := os.Getwd()
	assert.Nil(t, getErr)
	tplDir := fmt.Sprintf("%s/tplExample/module", workDir)
	rootDir := fmt.Sprintf("%s/tmp", workDir)
	cfg := &ModuleCfg{
		PackageName: "user",
		TableName:   "user",
		TplDir:      tplDir,
		RootDir:     rootDir,
	}
	autoCodeTool := NewGenerator()
	templateParam, getParamErr := autoCodeTool.GetModuleTemplateParam(db, cfg)
	assert.Nil(t, getParamErr)
	type Param struct {
		PackageName       string
		PackagePascalName string
		StructName        string
	}
	var params []GenParamsItem
	for _, tplItem := range templateParam.TemplateList {
		params = append(params, GenParamsItem{
			TargetDir:      tplItem.TargetDir,
			TargetFileName: tplItem.TargetFileName,
			Template:       tplItem.Template,
			Param: &Param{
				PackageName:       templateParam.PackageName,
				PackagePascalName: templateParam.PackagePascalName,
				StructName:        templateParam.StructName,
			},
		})
	}
	err := autoCodeTool.Gen(&GenParam{
		Params: params,
	})
	assert.Nil(t, err)
}

func TestCreateApiFile(t *testing.T) {
	// 获取当前的运行路径
	workDir, getErr := os.Getwd()
	assert.Nil(t, getErr)
	tplDir := fmt.Sprintf("%s/tplExample/api", workDir)
	rootDir := fmt.Sprintf("%s/tmp", workDir)
	cfg := &ControllerCfg{
		PackageName:    "user",
		TargetFilename: "user.go",
		TplDir:         tplDir,
		RootDir:        rootDir,
	}
	autoCodeTool := NewGenerator()
	templateParam, getParamErr := autoCodeTool.GetControllerTemplateParam(cfg)
	assert.Nil(t, getParamErr)
	type Param struct {
		PackageName       string
		PackagePascalName string
		FunctionName      string
	}
	var params []GenParamsItem
	for _, tplItem := range templateParam.TemplateList {
		params = append(params, GenParamsItem{
			TargetDir:      tplItem.TargetDir,
			TargetFileName: tplItem.TargetFileName,
			Template:       tplItem.Template,
			Param: &Param{
				PackageName:       templateParam.PackageName,
				PackagePascalName: templateParam.PackagePascalName,
				FunctionName:      "UserSetting",
			},
		})
	}
	err := autoCodeTool.Gen(&GenParam{
		Params: params,
	})
	assert.Nil(t, err)
}

func TestCreateModelFile(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True"
	db, openErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.Nil(t, openErr)
	// 获取当前的运行路径
	workDir, getErr := os.Getwd()
	assert.Nil(t, getErr)
	tplDir := fmt.Sprintf("%s/tplExample/model", workDir)
	rootDir := fmt.Sprintf("%s/tmp", workDir)
	cfg := &ModuleCfg{
		PackageName: "user",
		TableName:   "user",
		TplDir:      tplDir,
		RootDir:     rootDir,
	}
	autoCodeTool := NewGenerator()
	templateParam, getParamErr := autoCodeTool.GetModuleTemplateParam(db, cfg)
	assert.Nil(t, getParamErr)
	type ModelFieldItem struct {
		FieldName    string
		ColumnName   string
		Comment      string
		IsPrimaryKey bool
	}
	type Param struct {
		PackageName       string
		PackagePascalName string
		StructName        string
		TableName         string
		TableDescription  string
		ModelFields       []ModelFieldItem
	}

	var params []GenParamsItem
	for _, tplItem := range templateParam.TemplateList {
		var modelFields []ModelFieldItem

		for _, field := range tplItem.ModelFields {
			modelFields = append(modelFields, ModelFieldItem{
				FieldName:    field.FieldName,
				ColumnName:   field.ColumnName,
				Comment:      field.Comment,
				IsPrimaryKey: field.ColumnKey == "PRI",
			})
		}

		param := GenParamsItem{
			TargetDir:      tplItem.TargetDir,
			TargetFileName: tplItem.TargetFileName,
			Template:       tplItem.Template,
			Param: &Param{
				PackageName:       templateParam.PackageName,
				PackagePascalName: templateParam.PackagePascalName,
				StructName:        templateParam.StructName,
				ModelFields:       modelFields,
			},
		}
		params = append(params, param)
	}
	err := autoCodeTool.Gen(&GenParam{
		Params: params,
	})
	assert.Nil(t, err)
}
