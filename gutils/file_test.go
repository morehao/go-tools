package gutils

import (
	"fmt"
	"testing"
)

func TestGetFileExtension(t *testing.T) {
	name := "test.go"
	fmt.Println(GetFileExtension(name))
}

func TestTrimFileExtension(t *testing.T) {
	name := "test.go"
	fmt.Println(TrimFileExtension(name))
}
