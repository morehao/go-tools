package autoCode

import (
	"fmt"
	"github.com/morehao/go-tools/utils"
	"gorm.io/gorm"
)

type mysqlImpl struct {
	baseImpl
}

func (impl *mysqlImpl) GetModuleTemplateParam(db *gorm.DB, cfg *ModuleCfg) (*ModuleTemplateParams, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	if err := impl.checkCfg(cfg); err != nil {
		return nil, err
	}
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
	var templateList []ModuleTemplateParamsItem
	for _, tplItem := range tplList {
		targetDir := tplItem.BuildTargetDir(cfg.RootDir, cfg.PackageName)
		if appendDir, ok := layerAppendDirMap[tplItem.layerName]; ok {
			targetDir = fmt.Sprintf("%s/%s", targetDir, appendDir)
		}
		templateList = append(templateList, ModuleTemplateParamsItem{
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
	packagePascalName := utils.SnakeToPascal(cfg.PackageName)
	structName := utils.SnakeToPascal(cfg.TableName)
	res := &ModuleTemplateParams{
		PackageName:       cfg.PackageName,
		PackagePascalName: packagePascalName,
		TableName:         cfg.TableName,
		StructName:        structName,
		TemplateList:      templateList,
	}
	return res, nil
}

func (impl *mysqlImpl) checkCfg(cfg *ModuleCfg) error {
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

func (impl *mysqlImpl) CreateFile(param *CreateFileParam) error {
	if err := impl.checkCreateFileParam(param); err != nil {
		return err
	}
	for _, v := range param.Params {
		if err := createFile(v.TargetDir, v.TargetFileName, v.Template, v.Param); err != nil {
			return err
		}
	}
	return nil
}

func (impl *mysqlImpl) checkCreateFileParam(param *CreateFileParam) error {
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
			FiledName:  utils.SnakeToPascal(v.ColumnName),
			FieldType:  columnTypeMap[v.DataType],
			ColumnName: v.ColumnName,
			ColumnType: v.DataType,
			Comment:    v.ColumnComment,
		}
		modelFieldList = append(modelFieldList, item)
	}
	return modelFieldList, nil
}
