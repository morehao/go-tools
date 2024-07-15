package codeGen

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

const (
	dbTypeMysql = "mysql"
)

// mysqlTableColumn represents a column in the INFORMATION_SCHEMA.COLUMNS table
type mysqlTableColumn struct {
	ColumnName             string         `gorm:"column:COLUMN_NAME"`              // 列名
	DataType               string         `gorm:"column:DATA_TYPE"`                // 列的数据类型，如int
	ColumnType             string         `gorm:"column:COLUMN_TYPE"`              // 列的完整类型定义，如varchar(255)
	IsNullable             string         `gorm:"column:IS_NULLABLE"`              // 列是否允许 NULL 值。可能的值为 YES 或 NO
	ColumnDefault          sql.NullString `gorm:"column:COLUMN_DEFAULT"`           // 列的默认值
	ColumnComment          string         `gorm:"column:COLUMN_COMMENT"`           // 列的注释
	CharacterMaximumLength sql.NullInt64  `gorm:"column:CHARACTER_MAXIMUM_LENGTH"` // 字符串列的最大长度
	NumericPrecision       sql.NullInt64  `gorm:"column:NUMERIC_PRECISION"`        // 数值列的精度
	NumericScale           sql.NullInt64  `gorm:"column:NUMERIC_SCALE"`            // 数值列的小数位数
	DatetimePrecision      sql.NullInt64  `gorm:"column:DATETIME_PRECISION"`       // 日期时间列的精度
	CharacterSetName       sql.NullString `gorm:"column:CHARACTER_SET_NAME"`       // 字符串列的字符集名称
	CollationName          sql.NullString `gorm:"column:COLLATION_NAME"`           // 字符串列的排序规则名称
	OrdinalPosition        int64          `gorm:"column:ORDINAL_POSITION"`         // 列在表中的位置，从 1 开始
	ColumnKey              string         `gorm:"column:COLUMN_KEY"`               // 表示列是否是索引的一部分,可能的值为 PRI（主键）, UNI（唯一索引）, MUL（非唯一索引）
	Extra                  string         `gorm:"column:EXTRA"`                    // 列的额外信息，如 auto_increment
	Privileges             string         `gorm:"column:PRIVILEGES"`               // 与列相关的权限，如 select,insert,update,references
	GenerationExpression   sql.NullString `gorm:"column:GENERATION_EXPRESSION"`    // 生成列的表达式
}

type ModelField struct {
	FieldName         string // 字段名称
	FieldType         string // 字段数据类型，如int、string
	ColumnName        string // 列名
	ColumnType        string // 列数据类型，如varchar(255)
	ColumnSize        int    // 字段长度
	IsNull            bool   // 是否允许为空
	DefaultValue      string // 默认值
	ColumnKey         string // 索引类型
	Comment           string // 字段注释
	NumericPrecision  int64  // 数值列的精度
	NumericScale      int64  // 数值列的小数位数
	DatetimePrecision int64  // 日期时间列的精度
}

type TableList []string

func (l TableList) ToMap() map[string]struct{} {
	m := make(map[string]struct{}, len(l))
	for _, v := range l {
		m[v] = struct{}{}
	}
	return m
}

const (
	// 数值类型
	columnTypeInt       = "int"
	columnTypeTinyint   = "tinyint"
	columnTypeSmallint  = "smallint"
	columnTypeMediumint = "mediumint"
	columnTypeBigint    = "bigint"
	columnTypeDecimal   = "decimal"
	columnTypeNumeric   = "numeric"
	columnTypeFloat     = "float"
	columnTypeDouble    = "double"
	columnTypeBit       = "bit"
	columnTypeBool      = "bool"
	columnTypeBoolean   = "boolean"

	// 日期和时间类型
	columnTypeDate      = "date"
	columnTypeDatetime  = "datetime"
	columnTypeTimestamp = "timestamp"
	columnTypeTime      = "time"
	columnTypeYear      = "year"

	// 字符串类型
	columnTypeChar       = "char"
	columnTypeVarchar    = "varchar"
	columnTypeBinary     = "binary"
	columnTypeVarbinary  = "varbinary"
	columnTypeTinyblob   = "tinyblob"
	columnTypeBlob       = "blob"
	columnTypeMediumblob = "mediumblob"
	columnTypeLongblob   = "longblob"
	columnTypeTinytext   = "tinytext"
	columnTypeText       = "text"
	columnTypeMediumtext = "mediumtext"
	columnTypeLongtext   = "longtext"
	columnTypeEnum       = "enum"
	columnTypeSet        = "set"

	// JSON 类型
	columnTypeJSON = "json"

	// 空间数据类型
	columnTypeGeometry           = "geometry"
	columnTypePoint              = "point"
	columnTypeLinestring         = "linestring"
	columnTypePolygon            = "polygon"
	columnTypeMultiPoint         = "multipoint"
	columnTypeMultiLinestring    = "multilinestring"
	columnTypeMultiPolygon       = "multipolygon"
	columnTypeGeometryCollection = "geometrycollection"
)

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
