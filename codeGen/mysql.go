package codeGen

import (
	"fmt"
	"github.com/morehao/go-tools/gutils"
	"gorm.io/gorm"
)

type mysqlImpl struct {
}

func (impl *mysqlImpl) GetModuleTemplateParam(db *gorm.DB, cfg *ModuleCfg) (*ModelTemplateParamsRes, error) {
	dbName, getDbNameErr := getDbName(db)
	if getDbNameErr != nil {
		return nil, getDbNameErr
	}
	tableList, getTableErr := getTableList(db, dbName)
	if getTableErr != nil {
		return nil, getTableErr
	}
	tableMap := tableList.ToMap()
	if _, ok := tableMap[cfg.TableName]; !ok {
		return nil, fmt.Errorf("table %s not exist", cfg.TableName)
	}

	modelFieldList, getFieldErr := impl.getModelField(db, dbName, cfg)
	if getFieldErr != nil {
		return nil, getFieldErr
	}

	// 获取模板文件
	tplFiles, getTplErr := getTplFiles(cfg.TplDir)
	if getTplErr != nil {
		return nil, getTplErr
	}
	tplList, buildErr := buildTplCfg(tplFiles, fmt.Sprintf("%s%s", cfg.TableName, goFileSuffix))
	if buildErr != nil {
		return nil, buildErr
	}

	// 构造模板参数
	var templateList []ModelTemplateParamsItem
	for _, tplItem := range tplList {
		targetDir := tplItem.BuildTargetDir(cfg.RootDir, cfg.PackageName)
		if appendDir, ok := layerAppendDirMap[tplItem.layerName]; ok {
			targetDir = fmt.Sprintf("%s/%s", targetDir, appendDir)
		}
		templateList = append(templateList, ModelTemplateParamsItem{
			TemplateParamsItemBase: TemplateParamsItemBase{
				Template:       tplItem.template,
				Filename:       tplItem.filename,
				Filepath:       tplItem.filepath,
				OriginFilename: tplItem.originFilename,
				TargetFileName: tplItem.targetFileName,
				TargetDir:      targetDir,
				LayerName:      tplItem.layerName,
				LayerPrefix:    tplItem.layerPrefix,
			},
			ModelFields: modelFieldList,
		})
	}
	packagePascalName := gutils.SnakeToPascal(cfg.PackageName)
	structName := gutils.SnakeToPascal(cfg.TableName)
	res := &ModelTemplateParamsRes{
		PackageName:       cfg.PackageName,
		PackagePascalName: packagePascalName,
		TableName:         cfg.TableName,
		StructName:        structName,
		TemplateList:      templateList,
	}
	return res, nil
}

func (impl *mysqlImpl) getModelField(db *gorm.DB, dbName string, cfg *ModuleCfg) ([]ModelField, error) {
	var entities []mysqlTableColumn
	getColumnSql := fmt.Sprintf("SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s';", dbName, cfg.TableName)
	if err := db.Raw(getColumnSql).Scan(&entities).Error; err != nil {
		return nil, err
	}
	columnTypeMap := columnFieldTypeMap
	if len(cfg.ColumnTypeMap) > 0 {
		columnTypeMap = cfg.ColumnTypeMap
	}
	var modelFieldList []ModelField
	for _, v := range entities {
		item := ModelField{
			FieldName:  gutils.SnakeToPascal(v.ColumnName),
			FieldType:  columnTypeMap[v.DataType],
			ColumnName: v.ColumnName,
			ColumnType: v.DataType,
			ColumnKey:  v.ColumnKey,
			Comment:    v.ColumnComment,
		}
		modelFieldList = append(modelFieldList, item)
	}
	return modelFieldList, nil
}
