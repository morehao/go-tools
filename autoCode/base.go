package autoCode

import (
	"fmt"
	"github.com/morehao/go-tools/utils"
	"strings"
)

type baseImpl struct{}

func (impl *baseImpl) GetApiTemplateParam(cfg *ApiCfg) (*ApiTemplateParams, error) {
	if err := impl.checkCfg(cfg); err != nil {
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
	packagePascalName := utils.SnakeToPascal(cfg.PackageName)
	res := &ApiTemplateParams{
		PackageName:       cfg.PackageName,
		PackagePascalName: packagePascalName,
		TemplateList:      templateList,
	}
	return res, nil
}

func (impl *baseImpl) checkCfg(cfg *ApiCfg) error {
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
