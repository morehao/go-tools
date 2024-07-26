package main

import (
	"testing"
)

func TestCopyAndReplaceGoFile(t *testing.T) {
	dstDir := "../tmp/newCutter"
	cutter(dstDir)
}
