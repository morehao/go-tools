package gutils

import (
	"fmt"
	"testing"
)

func TestSnakeToPascal(t *testing.T) {
	fmt.Println(SnakeToPascal("workflow"))
}

func TestFirstLetterToLower(t *testing.T) {
	fmt.Println(FirstLetterToLower("Workflow"))
}
