// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package g_test

import (
	"testing"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/internal/gtest"
)

func getValue() int {
	return 3
}

func Test_Map_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m g.HashMap[int, int]
		m.Put(1, 11)
		t.Assert(m.Get(1), 11)
	})

	gtest.C(t, func(t *gtest.T) {
		var m g.HashMap[int, string]
		m.Put(1, "11")
		t.Assert(m.Get(1), "11")
	})
	gtest.C(t, func(t *gtest.T) {
		var m g.HashMap[string, string]
		m.Put("1", "11")
		t.Assert(m.Get("1"), "11")
	})
	gtest.C(t, func(t *gtest.T) {
		var m g.HashMap[string, int]
		m.Put("1", 11)
		t.Assert(m.Get("1"), 11)
	})
}

func Test_Map_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMap[string, string]()
		m.Put("key1", "val1")
		t.Assert(m.Keys(), []interface{}{"key1"})

		t.Assert(m.Get("key1"), "val1")
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrPut("key2", "val2"), "val2")
		t.Assert(m.PutIfAbsent("key2", "val2"), false)

		t.Assert(m.PutIfAbsent("key3", "val3"), true)

		val, removed := m.Remove("key2")
		t.Assert(val, "val2")
		t.Assert(removed, true)
		t.Assert(m.ContainsKey("key2"), false)

		t.AssertIN("key3", m.Keys())
		t.AssertIN("key1", m.Keys())
		t.AssertIN("val3", m.Values())
		t.AssertIN("val1", m.Values())

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
	})
}

func Test_Map_Set_Fun(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMap[string, int]()
		m.GetOrPutFunc("fun", getValue)
		m.GetOrPutFunc("funlock", getValue)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
		m.GetOrPutFunc("fun", getValue)
		t.Assert(m.PutIfAbsentFunc("fun", getValue), false)
		t.Assert(m.PutIfAbsentFunc("funlock", getValue), false)
	})
}

func Test_Map_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[string]string{"1": "1", "key1": "val1"}

		m := g.NewHashMapFrom[string, string](expect)
		m.Iterator(func(k string, v string) bool {
			t.Assert(expect[k], v)
			return true
		})
		// 断言返回值对遍历控制
		i := 0
		j := 0
		m.Iterator(func(k string, v string) bool {
			i++
			return true
		})
		m.Iterator(func(k string, v string) bool {
			j++
			return false
		})
		t.Assert(i, 2)
		t.Assert(j, 1)
	})
}

func Test_Map_Lock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[string]string{"1": "1", "key1": "val1"}
		m := g.NewHashMapFrom[string, string](expect)
		m.LockFunc(func(m map[string]string) {
			t.Assert(m, expect)
		})
		m.RLockFunc(func(m map[string]string) {
			t.Assert(m, expect)
		})
	})
}

func Test_Map_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// clone 方法是深克隆
		m := g.NewHashMapFrom[string, string](map[string]string{"1": "1", "key1": "val1"})
		m_clone := m.Clone()
		m.Remove("1")
		// 修改原 map,clone 后的 map 不影响
		t.AssertIN("1", m_clone.Keys())

		m_clone.Remove("key1")
		// 修改clone map,原 map 不影响
		t.AssertIN("key1", m.Keys())
	})
}

func Test_Map_Basic_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := g.NewHashMap[string, string]()
		m2 := g.NewHashMap[string, string]()
		m1.Put("key1", "val1")
		m2.Put("key2", "val2")
		m1.Merge(m2)
		t.Assert(m1.Map(), map[string]string{"key1": "val1", "key2": "val2"})
	})
}
