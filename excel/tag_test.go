package excel

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseFieldTags(t *testing.T) {
	type Dest struct {
		Name string `ex:"head:姓名,type:string,required,max:12"`
	}

	dest := Dest{}
	typ := reflect.TypeOf(dest)

	field, _ := typ.FieldByName("Name")
	tag := field.Tag.Get("ex")

	firstCtag, _ := parseFieldTags(tag)

	fmt.Println("Parsed Tags")
	for current := firstCtag; current != nil; current = current.next {
		fmt.Printf("Tag: %s, ParamsList: %s, Type: %d\n", current.tag, current.param, current.typeof)
	}
}
