package gast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFile(t *testing.T) {
	file := "./instance.go"
	res, err := ParseFile(file)
	assert.Nil(t, err)
	t.Log(res)
}

func TestAddMethodToInterface(t *testing.T) {
	// The file path of the Go source file you want to modify
	filePath := "./instance.go"

	// The interface and method names you want to add
	interfaceName := "User"
	methodName := "GetAge"

	// Add the method to the interface in the specified file
	err := addMethodToInterfaceInFile(filePath, interfaceName, methodName)
	assert.Nil(t, err)
}
