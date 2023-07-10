// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray_test

import (
	"testing"

	"github.com/wesleywu/gcontainer/garray"
)

type anySortedArrayItem struct {
	priority int64
	value    interface{}
}

var (
	anyArray = garray.NewArray[int]()
)

func Benchmark_AnyArray_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		anyArray.Add(i)
	}
}
