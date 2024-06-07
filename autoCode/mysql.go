package autoCode

import (
	"fmt"
	"github.com/morehao/go-tools/utils"
	"gorm.io/gorm"
	"text/template"
)

type mysqlImpl struct {
}

func (m *mysqlImpl) GetModuleTemplateParam(db *gorm.DB, cfg *ModuleCfg) (*TemplateParams, error) {
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

	modelFieldList, getFieldErr := m.getModelField(db, dbName, cfg)
	if getFieldErr != nil {
		return nil, getFieldErr
	}

	// 获取模板文件
	tplFiles, getTplErr := getTmplFiles(cfg.TplDir)
	if getTplErr != nil {
		return nil, getTplErr
	}
	var tplList []tplCfg
	for _, file := range tplFiles {
		targetFileName := file.originFilename
		if file.layerName != tplLayerNameDto {
			targetFileName = cfg.TableName + ".go"
		}
		tplItem := tplCfg{
			tplFile:        file,
			targetFileName: targetFileName,
		}
		tplList = append(tplList, tplItem)
	}
	// 解析模板
	for i, tplItem := range tplList {
		tmpl, parseErr := template.ParseFiles(tplItem.filepath)
		if parseErr != nil {
			return nil, parseErr
		}
		tplList[i].template = tmpl
	}
	packagePascalName := utils.SnakeToPascal(cfg.PackageName)
	structName := utils.SnakeToPascal(cfg.TableName)
	var templateList []TemplateItem
	for _, tplItem := range tplList {
		targetDir := tplItem.BuildTargetDir(cfg.RootDir, cfg.PackageName)
		if appendDir, ok := layerAppendDirMap[tplItem.layerName]; ok {
			targetDir = fmt.Sprintf("%s/%s", targetDir, appendDir)
		}
		templateList = append(templateList, TemplateItem{
			Template:       tplItem.template,
			Filename:       tplItem.filename,
			Filepath:       tplItem.filepath,
			OriginFilename: tplItem.originFilename,
			TargetFileName: tplItem.targetFileName,
			TargetDir:      targetDir,
			LayerName:      tplItem.layerName,
			LayerPrefix:    tplItem.layerPrefix,
			ModelFields:    modelFieldList,
		})
	}
	res := &TemplateParams{
		PackageName:       cfg.PackageName,
		PackagePascalName: packagePascalName,
		TableName:         cfg.TableName,
		StructName:        structName,
		TemplateList:      templateList,
	}
	return res, nil
}

func (m *mysqlImpl) CreateFile(param *CreateFileParam) error {
	for _, v := range param.Params {
		if err := createFile(v.TargetDir, v.TargetFileName, v.Template, v.Param); err != nil {
			return err
		}
	}
	return nil
}

func (m *mysqlImpl) getModelField(db *gorm.DB, dbName string, cfg *ModuleCfg) ([]ModelField, error) {
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
