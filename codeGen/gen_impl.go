package codeGen

import (
	"fmt"
	"github.com/morehao/go-tools/gutils"
	"gorm.io/gorm"
	"strings"
)

type generatorImpl struct{}

func (impl *generatorImpl) GetModuleTemplateParam(db *gorm.DB, cfg *ModuleCfg) (*ModelTemplateParamsRes, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	if err := impl.checkModuleCfg(cfg); err != nil {
		return nil, err
	}
	dbType := db.Dialector.Name()
	switch dbType {
	case dbTypeMysql:
		inst := &mysqlImpl{}
		return inst.GetModuleTemplateParam(db, cfg)
	default:
		return nil, fmt.Errorf("unsupported database type")
	}
}

func (impl *generatorImpl) checkModuleCfg(cfg *ModuleCfg) error {
	if cfg == nil {
		return fmt.Errorf("cfg is nil")
	}
	requiredFields := map[string]string{
		"packageName": cfg.PackageName,
		"tableName":   cfg.TableName,
		"tplDir":      cfg.TplDir,
		"rootDir":     cfg.RootDir,
	}

	for field, value := range requiredFields {
		if value == "" {
			return fmt.Errorf("%s is required", field)
		}
	}
	return nil
}

func (impl *generatorImpl) GetControllerTemplateParam(cfg *ControllerCfg) (*ControllerTemplateParams, error) {
	if err := impl.checkControllerCfg(cfg); err != nil {
		return nil, err
	}
	// 获取模板文件
	tplFiles, getTplErr := getTplFiles(cfg.TplDir)
	if getTplErr != nil {
		return nil, getTplErr
	}
	tplList, buildErr := buildTplCfg(tplFiles, cfg.TargetFilename)
	if buildErr != nil {
		return nil, buildErr
	}
	// 构造模板参数
	var templateList []TemplateParamsItemBase
	for _, tplItem := range tplList {
		targetDir := tplItem.BuildTargetDir(cfg.RootDir, cfg.PackageName)
		if appendDir, ok := layerAppendDirMap[tplItem.layerName]; ok {
			targetDir = fmt.Sprintf("%s/%s", targetDir, appendDir)
		}
		templateList = append(templateList, TemplateParamsItemBase{
			Template:       tplItem.template,
			Filename:       tplItem.filename,
			Filepath:       tplItem.filepath,
			OriginFilename: tplItem.originFilename,
			TargetFileName: tplItem.targetFileName,
			TargetDir:      targetDir,
			LayerName:      tplItem.layerName,
			LayerPrefix:    tplItem.layerPrefix,
		})
	}
	packagePascalName := gutils.SnakeToPascal(cfg.PackageName)
	res := &ControllerTemplateParams{
		PackageName:       cfg.PackageName,
		PackagePascalName: packagePascalName,
		TemplateList:      templateList,
	}
	return res, nil
}

func (impl *generatorImpl) checkControllerCfg(cfg *ControllerCfg) error {
	if cfg == nil {
		return fmt.Errorf("cfg is nil")
	}
	requiredFields := map[string]string{
		"packageName":    cfg.PackageName,
		"targetFilename": cfg.TargetFilename,
		"functionName":   cfg.FunctionName,
		"tplDir":         cfg.TplDir,
		"rootDir":        cfg.RootDir,
	}

	for field, value := range requiredFields {
		if value == "" {
			return fmt.Errorf("%s is required", field)
		}
	}
	if !strings.HasSuffix(cfg.TargetFilename, goFileSuffix) {
		return fmt.Errorf("targetFilename should end with %s", goFileSuffix)
	}
	return nil
}

func (impl *generatorImpl) Gen(param *GenParam) error {
	if err := impl.checkGenParam(param); err != nil {
		return err
	}
	for _, v := range param.Params {
		if err := createFile(v.TargetDir, v.TargetFileName, v.Template, v.Param); err != nil {
			return err
		}
	}
	return nil
}

func (impl *generatorImpl) checkGenParam(param *GenParam) error {
	if param == nil {
		return fmt.Errorf("param is nil")
	}
	if len(param.Params) == 0 {
		return fmt.Errorf("params is required")
	}
	for _, v := range param.Params {
		if v.TargetDir == "" {
			return fmt.Errorf("target dir is required")
		}
		if v.TargetFileName == "" {
			return fmt.Errorf("target file name is required")
		}
		if v.Template == nil {
			return fmt.Errorf("template is required")
		}
		if v.Param == nil {
			return fmt.Errorf("param is required")
		}
	}
	return nil
}
