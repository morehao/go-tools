package gast

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"testing"
)

func TestParseFile(t *testing.T) {
	file := "./instance.go"
	res, err := ParseFile(file)
	assert.Nil(t, err)
	t.Log(res)
}

func TestTrimFileTitle(t *testing.T) {
	file := "./instance.go"
	res, err := TrimFileTitle(file)
	assert.Nil(t, err)
	t.Log(res)
}

func TestFindMethodInFile(t *testing.T) {
	filePath := "./instance.go"

	method, ok, findErr := FindMethodInFile(filePath, "userImpl", "GetName")
	assert.Nil(t, findErr)
	assert.True(t, ok)

	src, readErr := os.ReadFile(filePath)
	assert.Nil(t, readErr)
	fileSet := token.NewFileSet()
	_, parseFileErr := parser.ParseFile(fileSet, filePath, src, parser.ParseComments)
	assert.Nil(t, parseFileErr)
	var buf bytes.Buffer
	printErr := printer.Fprint(&buf, fileSet, method)
	assert.Nil(t, printErr)

	t.Log(buf.String())
}

func TestFindFunctionInFile(t *testing.T) {
	filePath := "./instance.go"

	function, ok, findErr := FindFunctionInFile(filePath, "GetName")
	assert.Nil(t, findErr)
	assert.True(t, ok)

	src, readErr := os.ReadFile(filePath)
	assert.Nil(t, readErr)
	fileSet := token.NewFileSet()
	_, parseFileErr := parser.ParseFile(fileSet, filePath, src, parser.ParseComments)
	assert.Nil(t, parseFileErr)
	var buf bytes.Buffer
	printErr := printer.Fprint(&buf, fileSet, function)
	assert.Nil(t, printErr)

	t.Log(buf.String())
}

func TestAddMethodToInterface(t *testing.T) {
	filePath := "./instance.go"

	interfaceName := "User"
	methodName := "GetAge"

	err := AddMethodToInterfaceInFile(filePath, interfaceName, "userImpl", methodName)
	assert.Nil(t, err)
}

func TestAddContentToFunc(t *testing.T) {
	filePath := "./instance.go"
	content := "fmt.Println(1)"

	err := AddContentToFunc(content, "GetName", filePath)
	assert.Nil(t, err)
}
