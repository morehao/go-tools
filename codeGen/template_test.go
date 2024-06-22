package codeGen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetTmplFiles(t *testing.T) {
	// 获取当前的运行路径
	dir, getErr := os.Getwd()
	assert.Nil(t, getErr)
	files, getFileErr := getTplFiles(fmt.Sprintf("%s/tplExample", dir))
	assert.Nil(t, getFileErr)
	fmt.Println(files)
}
