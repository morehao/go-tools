package autoCode

import (
	"fmt"
	"github.com/morehao/go-tools/utils"
	"gorm.io/gorm"
	"text/template"
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
	tableList, getTableErr := getTableList(m.db, dbName)
	if getTableErr != nil {
		return getTableErr
	}
	tableMap := tableList.ToMap()
	if _, ok := tableMap[m.cfg.TableName]; !ok {
		return fmt.Errorf("table %s not exist", m.cfg.TableName)
	}

	columList, getColumnErr := getColumn(m.db, dbName, m.cfg.TableName)
	if getColumnErr != nil {
		fmt.Println(getColumnErr)
		return getColumnErr
	}
	fieldList := transferColumnFiled(columList)
	fmt.Println(utils.ToJson(fieldList))

	// 获取模板文件
	tplFiles, getTplErr := getTmplFiles(m.cfg.TplPath)
	if getTplErr != nil {
		return getTplErr
	}
	var tplList []tplCfg
	for _, file := range tplFiles {
		tplItem := tplCfg{
			tplFile:        file,
			targetFilepath: fmt.Sprintf("%s/%s/%s", m.cfg.PackageName, file.layerName, m.cfg.PackageName+".go"),
		}
		tplList = append(tplList, tplItem)
	}
	// 解析模板
	for i, tplItem := range tplList {
		tmpl, parseErr := template.New("test").Parse(tplItem.filepath)
		if parseErr != nil {
			return parseErr
		}
		tplList[i].template = *tmpl
	}
	fmt.Println(utils.ToJson(tplList))
	return nil
}
