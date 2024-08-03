package gutils

import (
	"container/list"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

func ToJsonString(v any) string {
	d, _ := jsoniter.MarshalToString(v)
	return d
}

func VToUint64(v any) uint64 {
	switch n := v.(type) {
	case int:
		return uint64(n)
	case int64:
		return uint64(n)
	case int8:
		return uint64(n)
	case int32:
		return uint64(n)
	case uint:
		return uint64(n)
	case uint64:
		return n
	case uint8:
		return uint64(n)
	case uint32:
		return uint64(n)
	case float64:
		return uint64(n)
	case float32:
		return uint64(n)
	case string:
		i, _ := strconv.ParseUint(n, 10, 64)
		return i
	}
	return 0
}

func LinkedListToArray(l *list.List) []interface{} {
	var array []interface{}
	for e := l.Front(); e != nil; e = e.Next() {
		array = append(array, e.Value)
	}
	return array
}
