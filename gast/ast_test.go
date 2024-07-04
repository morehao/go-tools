package gast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFile(t *testing.T) {
	file := "../glog/logger.go"
	res, err := ParseFile(file)
	assert.Nil(t, err)
	t.Log(res)
}
