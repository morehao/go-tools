package gutils

import "testing"

func TestSliceDiff(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 4}
	diff := SliceDiff(s1, s2)
	t.Log(ToJsonString(diff))
}

func TestSliceContain(t *testing.T) {
	s := []string{"a", "b"}
	t.Log(SliceContain(s, "a"))
}
