package gutils

import jsoniter "github.com/json-iterator/go"

func ToJsonString(v interface{}) string {
	d, _ := jsoniter.MarshalToString(v)
	return d
}
