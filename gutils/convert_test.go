package gutils

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedListToArray(t *testing.T) {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	var arr []int
	err := LinkedListToArray(l, &arr)
	assert.Nil(t, err)
	t.Log(ToJsonString(arr))
}

func TestToString(t *testing.T) {
	t.Log(ToString(123))
}
