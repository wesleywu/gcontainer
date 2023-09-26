// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package g_test

import (
	"testing"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/internal/gtest"
)

func Test_Array_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var array g.ArrayList[int]
		expect := []int{2, 3, 1}
		array.Add(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
	gtest.C(t, func(t *gtest.T) {
		var array g.ArrayList[int]
		expect := []int{2, 3, 1}
		array.Add(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
	gtest.C(t, func(t *gtest.T) {
		var array g.ArrayList[string]
		expect := []string{"b", "a"}
		array.Add("b", "a")
		t.Assert(array.Slice(), expect)
	})
}

func Test_IntArray_Unique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []int{1, 2, 3, 4, 5, 6}
		array := g.NewArrayList[int]()
		array.Add(1, 1, 2, 3, 3, 4, 4, 5, 5, 6, 6)
		array.Unique()
		t.Assert(array.Slice(), expect)
	})
}

func TestNewFromCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"100", "200", "300", "400", "500", "600"}
		array1 := g.NewArrayListFromCopy[string](a1)
		t.AssertIN(array1.PopRands(2), a1)
		t.Assert(len(array1.PopRands(1)), 1)
		t.Assert(len(array1.PopRands(9)), 3)
	})
}
