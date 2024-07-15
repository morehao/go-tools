package dao{{.PackagePascalName}}

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gutils"
	"go-api/errorCode"
	"gorm.io/gorm"
	"time"
)

// {{.StructName}} {{.TableDescription}}结构体
type {{.StructName}} struct {
    {{- range .ModelFields}}
    {{.FieldName}} string `gorm:"column:{{.ColumnName}};comment:{{.Comment}};{{if .IsPrimaryKey}}primarykey{{end}}"`
    {{- end}}
}

type {{.StructName}}List []{{.StructName}}

const {{.StructName}}TblName = "{{.TableName}}"

func ({{.StructName}}) TableName() string {
  return {{.StructName}}TblName
}

type {{.StructName}}Cond struct {
	Id             uint64
	Ids            []uint64
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type {{.StructName}}Dao struct {
    model.Base
}

func New{{.StructName}}Dao() *{{.StructName}}Dao {
	return &{{.StructName}}Dao{}
}

func (dao *{{.StructName}}Dao) WithTx(db *gorm.DB) *{{.StructName}}Dao {
	dao.Tx = db
	return dao
}

func (dao *{{.StructName}}Dao) Insert(ctx *gin.Context, entity *{{.StructName}}) error {
	db := dao.Db(ctx).Model(&{{.StructName}}{})
	db = db.Table({{.StructName}}TblName)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.WrapPrintf(err, "[{{.StructName}}Dao] Insert fail, entity:%s", utils.ToJson(entity))
	}
	return nil
}

func (dao *{{.StructName}}Dao) BatchInsert(ctx *gin.Context, entityList []{{.StructName}}) error {
	db := dao.Db(ctx).Table({{.StructName}}TblName)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.WrapPrintf(err, "[{{.StructName}}Dao] BatchInsert fail, entityList:%s", utils.ToJson(entityList))
	}
	return nil
}

func (dao *{{.StructName}}Dao) Update(ctx *gin.Context, entity *{{.StructName}}) error {
	db := dao.Db(ctx).Model(&{{.StructName}}{})
	db = db.Table({{.StructName}}TblName)
    if len(dao.Omit) > 0 {
        db = db.Omit(dao.Omit...)
    }
    if len(dao.Select) > 0 {
        db = db.Select(dao.Select)
    }
	if err := db.Where("id = ?", entity.Id).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.WrapPrintf(err, "[{{.StructName}}Dao] Update fail, entity:%s", utils.ToJson(entity))
	}
	return nil
}

func (dao *{{.StructName}}Dao) UpdateMap(ctx *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(ctx).Model(&{{.StructName}}{})
	db = db.Table({{.StructName}}TblName)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.WrapPrintf(err, "[{{.StructName}}Dao] UpdateMap fail, id:%d, updateMap:%s", id, utils.ToJson(updateMap))
	}
	return nil
}

func (dao *{{.StructName}}Dao) Delete(ctx *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(ctx).Model(&{{.StructName}}{})
	db = db.Table({{.StructName}}TblName)
	updatedField := map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.WrapPrintf(err, "[{{.StructName}}Dao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *{{.StructName}}Dao) GetById(ctx *gin.Context, id uint64) (*{{.StructName}}, error) {
	var entity {{.StructName}}
	db := dao.Db(ctx).Model(&{{.StructName}}{})
	db = db.Table({{.StructName}}TblName)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.WrapPrintf(err, "[{{.StructName}}Dao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *{{.StructName}}Dao) GetByCond(ctx *gin.Context,cond *{{.StructName}}Cond) (*{{.StructName}}, error) {
	var entity {{.StructName}}
	db := dao.Db(ctx).Model(&{{.StructName}}{})
	db = db.Table({{.StructName}}TblName)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.WrapPrintf(err, "[{{.StructName}}Dao] GetById fail, cond:%s", utils.ToJson(cond))
	}
	return &entity, nil
}

func (dao *{{.StructName}}Dao) GetListByCond(ctx *gin.Context,cond *{{.StructName}}Cond) ({{.StructName}}List, error) {
	var entityList {{.StructName}}List
	db := dao.Db(ctx).Model(&{{.StructName}}{})
	db = db.Table({{.StructName}}TblName)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.WrapPrintf(err, "[{{.StructName}}Dao] GetListByCond fail, cond:%s", utils.ToJson(cond))
	}
	return entityList, nil
}

func (dao *{{.StructName}}Dao) GetPageListByCond(ctx *gin.Context, cond *{{.StructName}}Cond) ({{.StructName}}List, int64, error) {
	db := dao.Db(ctx).Model(&{{.StructName}}{})
	db = db.Table({{.StructName}}TblName)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.WrapPrintf(err, "[{{.StructName}}Dao] GetPageListByCond count fail, cond:%s", utils.ToJson(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list {{.StructName}}List
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.WrapPrintf(err, "[{{.StructName}}Dao] GetPageListByCond find fail, cond:%s", utils.ToJson(cond))
	}
	return list, count, nil
}

func (l {{.StructName}}List) ToMap() map[uint64]{{.StructName}} {
	m := make(map[uint64]{{.StructName}})
	for _, v := range l {
		m[v.Id] = v
	}
	return m
}


func (dao *{{.StructName}}Dao) BuildCondition(db *gorm.DB, cond *{{.StructName}}Cond) {
	if cond.Id > 0 {
        query := fmt.Sprintf("%s.id = ?", {{.StructName}}TblName)
		db.Where(query, cond.Id)
	}
	if len(cond.Ids) > 0 {
	    query := fmt.Sprintf("%s.id in (?)", {{.StructName}}TblName)
		db.Where(query, cond.Ids)
	}
    if cond.CreatedAtStart > 0 {
        query := fmt.Sprintf("%s.created_at >= ?", {{.StructName}}TblName)
        db.Where(query, time.Unix(cond.CreatedAtStart, 0))
    }
    if cond.CreatedAtEnd > 0 {
        query := fmt.Sprintf("%s.created_at <= ?", {{.StructName}}TblName)
        db.Where(query, time.Unix(cond.CreatedAtEnd, 0))
    }
	if cond.IsDelete {
        db.Unscoped()
    }

	if cond.OrderField != "" {
		db.Order(cond.OrderField)
	}

	return
}
