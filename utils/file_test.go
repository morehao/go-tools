package utils

import (
	"fmt"
	"testing"
)

func TestGetFileSuffix(t *testing.T) {
	name := "test.go"
	fmt.Println(GetFileSuffix(name))
}
