package autoCode

import (
	"fmt"
	"github.com/morehao/go-tools/utils"
	"gorm.io/gorm"
)

const (
	dbTypeMysql = "mysql"
)

type TableColumn struct {
	Field   string `json:"field" gorm:"column:field"`     // 数据库字段
	Type    string `json:"type" gorm:"column:type"`       // 数据库字段类型
	Comment string `json:"comment" gorm:"column:comment"` // 数据库字段描述
}

type ModelField struct {
	FieldName    string `json:"fieldName"`    // 字段名称
	FieldType    string `json:"fieldType"`    // 字段类型
	FieldComment string `json:"fieldComment"` // 字段描述
}

type TableList []string

func (l TableList) ToMap() map[string]struct{} {
	m := make(map[string]struct{}, len(l))
	for _, v := range l {
		m[v] = struct{}{}
	}
	return m
}

var columnFieldTypeMap = map[string]string{
	"int":      "int64",
	"tinyint":  "int8",
	"smallint": "int",
	"bigint":   "int64",
	"varchar":  "string",
	"longtext": "string",
	"text":     "string",
	"blob":     "string",
	"datetime": "time.Time",
}

func getTableList(db *gorm.DB, dbName string) (tableList TableList, err error) {
	getTableSql := fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s';", dbName)
	if err = db.Raw(getTableSql).Scan(&tableList).Error; err != nil {
		return nil, err
	}
	return tableList, nil
}

func getDbName(db *gorm.DB) (dbName string, err error) {
	var entity struct {
		DbName string `gorm:"column:db_name"`
	}
	if err = db.Raw("SELECT DATABASE() db_name").Scan(&entity).Error; err != nil {
		return "", err
	}
	return entity.DbName, nil
}

func getColumn(db *gorm.DB, dbName, tableName string) (data []TableColumn, err error) {
	var entities []TableColumn

	getColumnSql := fmt.Sprintf("SELECT COLUMN_NAME field, DATA_TYPE type, COLUMN_COMMENT comment FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s';", dbName, tableName)
	if err = db.Raw(getColumnSql).Scan(&entities).Error; err != nil {
		return nil, err
	}
	return entities, err
}

func transferColumnFiled(columnList []TableColumn) []ModelField {
	fieldList := make([]ModelField, 0, len(columnList))
	for _, v := range columnList {
		item := ModelField{
			FieldName:    utils.SnakeToPascal(v.Field),
			FieldType:    columnFieldTypeMap[v.Type],
			FieldComment: v.Comment,
		}

		fieldList = append(fieldList, item)
	}
	return fieldList
}
