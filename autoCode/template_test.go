package autoCode

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTmplName(t *testing.T) {
	tmplName()
}

func TestGetTmplFiles(t *testing.T) {
	// 获取当前的运行路径
	dir, getErr := os.Getwd()
	assert.Nil(t, getErr)
	// fmt.Println(filepath.Dir(dir))
	files, getFileErr := getTmplFiles(fmt.Sprintf("%s/tplExample", dir))
	assert.Nil(t, getFileErr)
	fmt.Println(files)
}
