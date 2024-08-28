package dao{{.PackagePascalName}}

import (
	"fmt"
	"{{.ProjectRootDir}}/internal/pkg/errorCode"
	"{{.ProjectRootDir}}/internal/{{.ServiceName}}/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gutils"
	"gorm.io/gorm"
)

// {{.StructName}}Entity {{.Description}}表结构体
type {{.StructName}}Entity struct {
{{- range .ModelFields}}
	{{- if .IsPrimaryKey}}
	{{.FieldName}} uint64 `gorm:"column:{{.ColumnName}};comment:{{.Comment}};primaryKey"`
	{{- else if eq .FieldName "DeletedAt"}}
	{{.FieldName}} gorm.DeletedAt `gorm:"column:{{.ColumnName}};comment:{{.Comment}}"`
	{{- else}}
	{{.FieldName}} {{.FieldType}} `gorm:"column:{{.ColumnName}};comment:{{.Comment}}"`
	{{- end}}
{{- end}}
}

type {{.StructName}}EntityList []{{.StructName}}Entity

const TblName{{.StructName}} = "{{.TableName}}"

func ({{.StructName}}Entity ) TableName() string {
  return TblName{{.StructName}}
}

type {{.StructName}}Cond struct {
	ID             uint64
	IDs            []uint64
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

func (dao *{{.StructName}}Dao) Insert(c *gin.Context, entity *{{.StructName}}Entity) error {
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(TblName{{.StructName}})
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[{{.StructName}}Dao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *{{.StructName}}Dao) BatchInsert(c *gin.Context, entityList {{.StructName}}EntityList) error {
	db := dao.Db(c).Table(TblName{{.StructName}})
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[{{.StructName}}Dao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *{{.StructName}}Dao) Update(c *gin.Context, entity *{{.StructName}}Entity) error {
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(TblName{{.StructName}})
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[{{.StructName}}Dao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *{{.StructName}}Dao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(TblName{{.StructName}})
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[{{.StructName}}Dao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *{{.StructName}}Dao) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(TblName{{.StructName}})
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by": deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[{{.StructName}}Dao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *{{.StructName}}Dao) GetById(c *gin.Context, id uint64) (*{{.StructName}}Entity, error) {
	var entity {{.StructName}}Entity
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(TblName{{.StructName}})
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *{{.StructName}}Dao) GetByCond(c *gin.Context,cond *{{.StructName}}Cond) (*{{.StructName}}Entity, error) {
	var entity {{.StructName}}Entity
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(TblName{{.StructName}})

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *{{.StructName}}Dao) GetListByCond(c *gin.Context,cond *{{.StructName}}Cond) ({{.StructName}}EntityList, error) {
	var entityList {{.StructName}}EntityList
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(TblName{{.StructName}})

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *{{.StructName}}Dao) GetPageListByCond(c *gin.Context, cond *{{.StructName}}Cond) ({{.StructName}}EntityList, int64, error) {
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(TblName{{.StructName}})

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list {{.StructName}}EntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l {{.StructName}}EntityList) ToMap() map[uint64]{{.StructName}}Entity {
	m := make(map[uint64]{{.StructName}}Entity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *{{.StructName}}Dao) CountByCond(c *gin.Context, cond *{{.StructName}}Cond) (int64, error) {
	db := dao.Db(c).Model(&{{.StructName}}Entity{})
	db = db.Table(TblName{{.StructName}})

	dao.BuildCondition(db, cond)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, errorCode.ErrorDbFind.Wrapf(err, "[{{.StructName}}Dao] CountByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return count, nil
}


func (dao *{{.StructName}}Dao) BuildCondition(db *gorm.DB, cond *{{.StructName}}Cond) {
	if cond.ID > 0 {
        query := fmt.Sprintf("%s.id = ?", TblName{{.StructName}})
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
	    query := fmt.Sprintf("%s.id in (?)", TblName{{.StructName}})
		db.Where(query, cond.IDs)
	}
    if cond.CreatedAtStart > 0 {
        query := fmt.Sprintf("%s.created_at >= ?", TblName{{.StructName}})
        db.Where(query, time.Unix(cond.CreatedAtStart, 0))
    }
    if cond.CreatedAtEnd > 0 {
        query := fmt.Sprintf("%s.created_at <= ?", TblName{{.StructName}})
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
