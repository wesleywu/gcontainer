package equal

import (
	"reflect"

	"github.com/wesleywu/gcontainer/utils/empty"
)

func Equals[T any](a, b T) bool {
	aIsNil := empty.IsNil(a)
	bIsNil := empty.IsNil(b)
	if aIsNil && bIsNil {
		return true
	} else if aIsNil || bIsNil {
		return false
	}
	v := reflect.ValueOf(a)
	if v.Comparable() {
		u := reflect.ValueOf(b)
		return v.Equal(u)
	}
	panic("type cannot be compared")
}
