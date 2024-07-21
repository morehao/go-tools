package gast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddMethodToInterfaceInFile(t *testing.T) {
	filePath := "./test.go"

	err := AddMethodToInterfaceInFile(filePath, "userImpl", "GetAge", "User")
	assert.Nil(t, err)
}

func TestAddContentToFunc(t *testing.T) {
	filePath := "./test.go"
	content := `routerGroup.POST("test")`

	err := AddContentToFunc(filePath, "platformRouter", content)
	assert.Nil(t, err)
}

func TestAddFunction(t *testing.T) {
	content := `
func NewFunction() {
	fmt.Println("Hello, World!")
}
`
	err := AddFunction("test.go", content, "gast")
	assert.Nil(t, err)
}

func TestAddMethodToInterface(t *testing.T) {
	filePath := "test.go"
	content, err := getMethodDeclaration(filePath, "userImpl", "GetAge")
	assert.Nil(t, err)
	t.Log(content)
	interfaceName := "User"
	err = AddMethodToInterface(filePath, "userImpl", "GetAge", interfaceName)
	assert.Nil(t, err)
}

func TestAddContentToFuncWithLineNumber(t *testing.T) {
	filePath := "./test.go"
	content := `routerGroup.POST("test1")`
	err := AddContentToFuncWithLineNumber(filePath, "platformRouter", content, -2)
	assert.Nil(t, err)
}
