package utils

import jsoniter "github.com/json-iterator/go"

func ToJson(v interface{}) string {
	d, _ := jsoniter.MarshalToString(v)
	return d
}
