package autoCode

import (
	"fmt"
	"github.com/morehao/go-tools/utils"
	"gorm.io/gorm"
)

type mysqlImpl struct {
	db  *gorm.DB
	cfg *Cfg
}

func (m *mysqlImpl) Generate() error {
	dbName, getDbNameErr := getDbName(m.db)
	if getDbNameErr != nil {
		return getDbNameErr
	}

	columList, getColumnErr := getColumn(m.db, dbName, m.cfg.TableName)
	if getColumnErr != nil {
		fmt.Println(getColumnErr)
		return getColumnErr
	}
	fieldList := transferColumnFiled(columList)
	fmt.Println(utils.ToJson(fieldList))
	return nil
}
