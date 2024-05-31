package autoCode

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGenerate(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	db, openErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.Nil(t, openErr)

	cfg := &Cfg{
		TableName:   "user",
		PackageName: "user",
		PrefixName:  "user",
	}
	autoCodeTool := NewAutoCode(db, cfg)
	err := autoCodeTool.Generate()
	assert.Nil(t, err)
}
