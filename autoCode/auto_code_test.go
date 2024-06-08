package autoCode

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
	rootDir := fmt.Sprintf("%s/tmpAutoCode", workDir)
	cfg := &ModuleCfg{
		PackageName: "user",
		TableName:   "user",
		TplDir:      tplDir,
		RootDir:     rootDir,
	}
	autoCodeTool := NewAutoCode()
	templateParam, getParamErr := autoCodeTool.GetModuleTemplateParam(db, cfg)
	assert.Nil(t, getParamErr)
	type Param struct {
		PackageName       string
		PackagePascalName string
		StructName        string
	}
	var params []CreateFileParamsItem
	for _, tplItem := range templateParam.TemplateList {
		params = append(params, CreateFileParamsItem{
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
	err := autoCodeTool.CreateFile(&CreateFileParam{
		Params: params,
	})
	assert.Nil(t, err)
}

func TestCreateApiFile(t *testing.T) {
	// 获取当前的运行路径
	workDir, getErr := os.Getwd()
	assert.Nil(t, getErr)
	tplDir := fmt.Sprintf("%s/tplExample/api", workDir)
	rootDir := fmt.Sprintf("%s/tmpAutoCode", workDir)
	cfg := &ApiCfg{
		PackageName:    "user",
		TargetFilename: "user.go",
		TplDir:         tplDir,
		RootDir:        rootDir,
	}
	autoCodeTool := NewAutoCode()
	templateParam, getParamErr := autoCodeTool.GetApiTemplateParam(cfg)
	assert.Nil(t, getParamErr)
	type Param struct {
		PackageName       string
		PackagePascalName string
		FunctionName      string
	}
	var params []CreateFileParamsItem
	for _, tplItem := range templateParam.TemplateList {
		params = append(params, CreateFileParamsItem{
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
	err := autoCodeTool.CreateFile(&CreateFileParam{
		Params: params,
	})
	assert.Nil(t, err)
}
