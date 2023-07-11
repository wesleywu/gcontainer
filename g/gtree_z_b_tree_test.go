// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package g_test

import (
	"fmt"
	"testing"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/internal/gtest"
	"github.com/wesleywu/gcontainer/utils/comparator"
)

func Test_BTree_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewBTree[string, string](3, comparator.ComparatorString)
		m.Put("key1", "val1")

		t.Assert(m.Height(), 1)

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

		m2 := g.NewBTreeFrom[string, string](3, comparator.ComparatorString, map[string]string{"1": "1", "key1": "val1"})
		t.Assert(m2.Map(), map[string]string{"1": "1", "key1": "val1"})
	})
}

func Test_BTree_Set_Fun(t *testing.T) {
	//GetOrPutFunc lock or unlock
	gtest.C(t, func(t *gtest.T) {
		m := g.NewBTree[string, int](3, comparator.ComparatorString)
		t.Assert(m.GetOrPutFunc("fun", getValue), 3)
		t.Assert(m.GetOrPutFunc("fun", getValue), 3)
		t.Assert(m.GetOrPutFunc("funlock", getValue), 3)
		t.Assert(m.GetOrPutFunc("funlock", getValue), 3)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
	})
	//PutIfAbsentFunc lock or unlock
	gtest.C(t, func(t *gtest.T) {
		m := g.NewBTree[string, int](3, comparator.ComparatorString)
		t.Assert(m.PutIfAbsentFunc("fun", getValue), true)
		t.Assert(m.PutIfAbsentFunc("fun", getValue), false)
		t.Assert(m.PutIfAbsentFunc("funlock", getValue), true)
		t.Assert(m.PutIfAbsentFunc("funlock", getValue), false)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
	})

}

func Test_BTree_Batch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewBTree[string, string](3, comparator.ComparatorString)
		m.Puts(map[string]string{"1": "1", "key1": "val1", "key2": "val2", "key3": "val3"})
		t.Assert(m.Map(), map[string]string{"1": "1", "key1": "val1", "key2": "val2", "key3": "val3"})
		m.Removes([]string{"key1", "1"})
		t.Assert(m.Map(), map[interface{}]interface{}{"key2": "val2", "key3": "val3"})
	})
}

func Test_BTree_Iterator(t *testing.T) {
	keys := []string{"1", "key1", "key2", "key3", "key4"}
	keyLen := len(keys)
	index := 0

	expect := map[string]string{"key4": "val4", "1": "1", "key1": "val1", "key2": "val2", "key3": "val3"}

	m := g.NewBTreeFrom[string, string](3, comparator.ComparatorString, expect)

	gtest.C(t, func(t *gtest.T) {
		m.Iterator(func(k string, v string) bool {
			t.Assert(k, keys[index])
			index++
			t.Assert(expect[k], v)
			return true
		})

		m.IteratorDesc(func(k string, v string) bool {
			index--
			t.Assert(k, keys[index])
			t.Assert(expect[k], v)
			return true
		})
	})

	m.Print()
	// 断言返回值对遍历控制
	gtest.C(t, func(t *gtest.T) {
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
		t.Assert(i, keyLen)
		t.Assert(j, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		i := 0
		j := 0
		m.IteratorDesc(func(k string, v string) bool {
			i++
			return true
		})
		m.IteratorDesc(func(k string, v string) bool {
			j++
			return false
		})
		t.Assert(i, keyLen)
		t.Assert(j, 1)
	})
}

func Test_BTree_IteratorFrom(t *testing.T) {
	m := make(map[int]int)
	for i := 1; i <= 10; i++ {
		m[i] = i * 10
	}
	tree := g.NewBTreeFrom[int, int](3, comparator.ComparatorInt, m)

	gtest.C(t, func(t *gtest.T) {
		n := 5
		tree.IteratorFrom(5, true, func(key, value int) bool {
			t.Assert(n, key)
			t.Assert(n*10, value)
			n++
			return true
		})

		i := 5
		tree.IteratorAscFrom(5, true, func(key, value int) bool {
			t.Assert(i, key)
			t.Assert(i*10, value)
			i++
			return true
		})

		j := 5
		tree.IteratorDescFrom(5, true, func(key, value int) bool {
			t.Assert(j, key)
			t.Assert(j*10, value)
			j--
			return true
		})
	})
}

func Test_BTree_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		//clone 方法是深克隆
		m := g.NewBTreeFrom[string, string](3, comparator.ComparatorString, map[string]string{"1": "1", "key1": "val1"})
		m_clone := m.Clone()
		m.Remove("1")
		//修改原 map,clone 后的 map 不影响
		t.AssertIN("1", m_clone.Keys())

		m_clone.Remove("key1")
		//修改clone map,原 map 不影响
		t.AssertIN("key1", m.Keys())
	})
}

func Test_BTree_LRNode(t *testing.T) {
	expect := map[string]string{"key4": "val4", "key1": "val1", "key2": "val2", "key3": "val3"}
	//safe
	gtest.C(t, func(t *gtest.T) {
		m := g.NewBTreeFrom[string, string](3, comparator.ComparatorString, expect)
		t.Assert(m.Left().Key(), "key1")
		t.Assert(m.Right().Key(), "key4")
	})
	//unsafe
	gtest.C(t, func(t *gtest.T) {
		m := g.NewBTreeFrom(3, comparator.ComparatorString, expect, true)
		t.Assert(m.Left().Key(), "key1")
		t.Assert(m.Right().Key(), "key4")
	})
}

func Test_BTree_Remove(t *testing.T) {
	m := g.NewBTree[int, string](3, comparator.ComparatorInt)
	for i := 1; i <= 100; i++ {
		m.Put(i, fmt.Sprintf("val%d", i))
	}
	expect := m.Map()
	gtest.C(t, func(t *gtest.T) {
		for k, v := range expect {
			m1 := m.Clone()
			val, removed := m1.Remove(k)
			t.Assert(val, v)
			t.Assert(removed, true)
			val, removed = m1.Remove(k)
			t.Assert(val, "")
			t.Assert(removed, false)
		}
	})
}