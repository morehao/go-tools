package autoCode

import (
	"fmt"
	"github.com/morehao/go-tools/utils"
	"gorm.io/gorm"
	"text/template"
)

type mysqlImpl struct {
	db     *gorm.DB
	dbName string
	cfg    *Cfg
}

func (m *mysqlImpl) Generate() error {
	dbName, getDbNameErr := getDbName(m.db)
	if getDbNameErr != nil {
		return getDbNameErr
	}
	m.dbName = dbName
	tableList, getTableErr := getTableList(m.db, dbName)
	if getTableErr != nil {
		return getTableErr
	}
	tableMap := tableList.ToMap()
	if _, ok := tableMap[m.cfg.TableName]; !ok {
		return fmt.Errorf("table %s not exist", m.cfg.TableName)
	}

	modelFieldList, getFieldErr := m.getModelField()
	if getFieldErr != nil {
		return getFieldErr
	}
	fmt.Println(utils.ToJson(modelFieldList))

	// 获取模板文件
	tplFiles, getTplErr := getTmplFiles(m.cfg.TplDir)
	if getTplErr != nil {
		return getTplErr
	}
	var tplList []tplCfg
	for _, file := range tplFiles {
		targetFileName := file.originFilename
		if file.layerName != tplLayerNameDto {
			targetFileName = m.cfg.TableName + ".go"
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
			return parseErr
		}
		tplList[i].template = tmpl
	}

	packagePascalName := utils.SnakeToPascal(m.cfg.PackageName)
	structName := utils.SnakeToPascal(m.cfg.TableName)
	params := &tplParam{
		PackageName:       m.cfg.PackageName,
		TableName:         m.cfg.TableName,
		PackagePascalName: packagePascalName,
		StructName:        structName,
	}

	// 渲染模板
	for _, tplItem := range tplList {
		codeDir := tplItem.GetCodeDir(m.cfg.RootDir, structName)
		if err := createFile(codeDir, &tplItem, params); err != nil {
			return err
		}
	}

	return nil
}

func (m *mysqlImpl) getModelField() ([]ModelField, error) {
	var entities []mysqlTableColumn
	getColumnSql := fmt.Sprintf("SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s';", m.dbName, m.cfg.TableName)
	if err := m.db.Raw(getColumnSql).Scan(&entities).Error; err != nil {
		return nil, err
	}
	var modelFieldList []ModelField
	for _, v := range entities {
		item := ModelField{
			FiledName:  utils.SnakeToPascal(v.ColumnName),
			FieldType:  columnFieldTypeMap[v.DataType],
			ColumnName: v.ColumnName,
			ColumnType: v.DataType,
			Comment:    v.ColumnComment,
		}
		modelFieldList = append(modelFieldList, item)
	}
	return modelFieldList, nil
}
