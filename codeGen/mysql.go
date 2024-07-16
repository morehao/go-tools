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
	tplFiles, getTplErr := getTplFiles(cfg.TplDir, cfg.LayerNameMap, cfg.LayerPrefixMap)
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
		templateList = append(templateList, ModelTemplateParamsItem{
			TemplateParamsItemBase: tplParamsItem,
			ModelFields:            modelFieldList,
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
	columnTypeMap := mysqlDefaultColumnTypeMap
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

var mysqlDefaultColumnTypeMap = map[string]string{
	"tinyint":    "int8",
	"smallint":   "int16",
	"mediumint":  "int32",
	"int":        "int32",
	"integer":    "int32",
	"bigint":     "int64",
	"float":      "float32",
	"double":     "float64",
	"decimal":    "string", // 或者使用 "big.Rat" 或 "float64"，取决于精度需求
	"date":       "time.Time",
	"datetime":   "time.Time",
	"timestamp":  "time.Time",
	"time":       "time.Duration",
	"year":       "int16",
	"char":       "string",
	"varchar":    "string",
	"text":       "string",
	"tinytext":   "string",
	"mediumtext": "string",
	"longtext":   "string",
	"blob":       "[]byte",
	"tinyblob":   "[]byte",
	"mediumblob": "[]byte",
	"longblob":   "[]byte",
	"enum":       "string", // 或者自定义类型
	"set":        "string", // 或者自定义类型，可能是字符串切片
	"bit":        "[]byte", // 或者 "uint64"，取决于位数
	"binary":     "[]byte",
	"varbinary":  "[]byte",
	"json":       "json.RawMessage", // 或者 "map[string]interface{}" 或自定义结构体
	"bool":       "bool",
	"boolean":    "bool",
}
