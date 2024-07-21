package gerror

import (
	"testing"
)

func TestWrap(t *testing.T) {
	err1 := Error{
		Code: 1,
		Msg:  "test1",
	}
	err2 := Error{
		Code: 2,
		Msg:  "test2",
	}
	err := err1.Wrap(err2)
	t.Log(err)
}

func TestWrapf(t *testing.T) {
	err1 := Error{
		Code: 1,
		Msg:  "test1",
	}
	err2 := Error{
		Code: 2,
		Msg:  "test2",
	}
	err := err1.Wrapf(err2, "here is errMsg:%s", "123")
	t.Log(err)
}
