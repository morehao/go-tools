package gutils

import (
	"container/list"
	"fmt"
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
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

func LinkedListToArray(l *list.List, dest any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("destination must be a pointer to a slice")
	}

	sliceValue := destValue.Elem()
	elementType := sliceValue.Type().Elem()

	for e := l.Front(); e != nil; e = e.Next() {
		newElement := reflect.New(elementType).Elem()
		newElement.Set(reflect.ValueOf(e.Value).Convert(elementType))
		sliceValue = reflect.Append(sliceValue, newElement)
	}

	destValue.Elem().Set(sliceValue)
	return nil
}
