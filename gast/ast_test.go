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
	file := "./test.go"
	res, err := ParseFile(file)
	assert.Nil(t, err)
	t.Log(res)
}

func TestTrimFileTitle(t *testing.T) {
	file := "./test.go"
	res, err := TrimFileTitle(file)
	assert.Nil(t, err)
	t.Log(res)
}

func TestFindMethodInFile(t *testing.T) {
	filePath := "./test.go"

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
	filePath := "./test.go"

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
	filePath := "./test.go"

	interfaceName := "User"
	methodName := "GetAge"

	err := AddMethodToInterfaceInFile(filePath, interfaceName, "userImpl", methodName)
	assert.Nil(t, err)
}

func TestAddContentToFunc(t *testing.T) {
	filePath := "./test.go"
	content := "fmt.Println(1)"

	err := AddContentToFunc(content, "GetName", filePath)
	assert.Nil(t, err)
}

func TestAddFunction(t *testing.T) {
	content := `
func NewFunction() {
	fmt.Println("Hello, World!")
}
`
	err := AddFunction(content, "test.go", "gast")
	assert.Nil(t, err)
}
