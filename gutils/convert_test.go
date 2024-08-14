package gutils

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"testing"
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
