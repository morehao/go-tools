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

func TestFindMethodInFile(t *testing.T) {
	filePath := "./instance.go"
	src, readErr := os.ReadFile(filePath)
	assert.Nil(t, readErr)

	fileSet := token.NewFileSet()
	file, parseFileErr := parser.ParseFile(fileSet, filePath, src, parser.ParseComments)
	assert.Nil(t, parseFileErr)

	method, ok := FindMethodInFile(file, "userImpl", "GetName")
	assert.True(t, ok)
	var buf bytes.Buffer
	printErr := printer.Fprint(&buf, fileSet, method)
	assert.Nil(t, printErr)

	t.Log(buf.String())
}

func TestFindFunctionInFile(t *testing.T) {
	filePath := "./instance.go"
	src, readErr := os.ReadFile(filePath)
	assert.Nil(t, readErr)

	fileSet := token.NewFileSet()
	file, parseFileErr := parser.ParseFile(fileSet, filePath, src, parser.ParseComments)
	assert.Nil(t, parseFileErr)

	function, ok := FindFunctionInFile(file, "GetName")
	assert.True(t, ok)
	var buf bytes.Buffer
	printErr := printer.Fprint(&buf, fileSet, function)
	assert.Nil(t, printErr)

	t.Log(buf.String())
}

func TestAddMethodToInterface(t *testing.T) {
	filePath := "./instance.go"

	interfaceName := "User"
	methodName := "GetName"

	err := AddMethodToInterfaceInFile(filePath, interfaceName, "userImpl", methodName)
	assert.Nil(t, err)
}
