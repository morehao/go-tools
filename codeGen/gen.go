package codeGen

import (
	"fmt"
	"github.com/morehao/go-tools/gutils"
	"gorm.io/gorm"
	"strings"
)

type Generator interface {
	GetModuleTemplateParams(db *gorm.DB, cfg *ModuleCfg) (*ModelTemplateParamsRes, error)
	GetControllerTemplateParams(cfg *ControllerCfg) (*ControllerTemplateParamsRes, error)
	Gen(params *GenParams) error
}

func NewGenerator() Generator {
	return &generatorImpl{}
}

type generatorImpl struct{}

func (impl *generatorImpl) GetModuleTemplateParams(db *gorm.DB, cfg *ModuleCfg) (*ModelTemplateParamsRes, error) {
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

func (impl *generatorImpl) GetControllerTemplateParams(cfg *ControllerCfg) (*ControllerTemplateParamsRes, error) {
	if err := impl.checkControllerCfg(cfg); err != nil {
		return nil, err
	}
	// 获取模板文件
	tplFiles, getTplErr := getTplFiles(cfg.TplDir, cfg.LayerNameMap, cfg.LayerPrefixMap)
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
		rootDir := cfg.RootDir
		if layerDir, ok := cfg.LayerDirMap[tplItem.layerName]; ok {
			rootDir = layerDir
		}
		targetDir := tplItem.BuildTargetDir(rootDir, cfg.PackageName)
		tplParamsItem := TemplateParamsItemBase{
			Template:       tplItem.template,
			TplFilename:    tplItem.filename,
			TplFilepath:    tplItem.filepath,
			OriginFilename: tplItem.originFilename,
			TargetFilename: tplItem.targetFileName,
			TargetDir:      targetDir,
			LayerName:      tplItem.layerName,
			LayerPrefix:    tplItem.layerPrefix,
		}
		if gutils.FileExists(fmt.Sprintf("%s/%s", targetDir, tplItem.targetFileName)) {
			tplParamsItem.TargetFileExist = true
		}
		templateList = append(templateList, tplParamsItem)
	}
	packagePascalName := gutils.SnakeToPascal(cfg.PackageName)
	res := &ControllerTemplateParamsRes{
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

func (impl *generatorImpl) Gen(params *GenParams) error {
	if err := impl.checkGenParams(params); err != nil {
		return err
	}
	for _, v := range params.ParamsList {
		if err := createFile(v.TargetDir, v.TargetFileName, v.Template, v.ExtraParams); err != nil {
			return err
		}
	}
	return nil
}

func (impl *generatorImpl) checkGenParams(params *GenParams) error {
	if params == nil {
		return fmt.Errorf("params is nil")
	}
	if len(params.ParamsList) == 0 {
		return fmt.Errorf("params is required")
	}
	for _, v := range params.ParamsList {
		if v.TargetDir == "" {
			return fmt.Errorf("target dir is required")
		}
		if v.TargetFileName == "" {
			return fmt.Errorf("target file name is required")
		}
		if v.Template == nil {
			return fmt.Errorf("template is required")
		}
		if v.ExtraParams == nil {
			return fmt.Errorf("params is required")
		}
	}
	return nil
}
