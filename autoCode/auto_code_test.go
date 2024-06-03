package autoCode

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	db, openErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.Nil(t, openErr)
	// 获取当前的运行路径
	workDir, getErr := os.Getwd()
	assert.Nil(t, getErr)
	tplDir := fmt.Sprintf("%s/tplExample", workDir)
	rootDir := fmt.Sprintf("%s/tmpAutoCode", workDir)
	cfg := &Cfg{
		TableName:   "user",
		PackageName: "user",
		TplDir:      tplDir,
		RootDir:     rootDir,
	}
	autoCodeTool := NewAutoCode(db, cfg)
	err := autoCodeTool.Generate()
	assert.Nil(t, err)
}
