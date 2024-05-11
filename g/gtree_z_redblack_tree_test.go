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
	"github.com/wesleywu/gcontainer/utils/comparators"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

func Test_RedBlackTree_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewRedBlackTree[string, string](comparators.ComparatorString)
		m.Put("key1", "val1")
		t.Assert(m.Keys(), []interface{}{"key1"})

		t.Assert(m.Get("key1"), "val1")
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrPut("key2", "val2"), "val2")
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

		m.Puts(map[string]string{"key3": "val3", "key1": "val1"})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)

		m2 := g.NewRedBlackTreeFrom(comparators.ComparatorString, map[string]string{"1": "1", "key1": "val1"})
		t.Assert(m2.Map(), map[string]string{"1": "1", "key1": "val1"})
	})
}

func Test_RedBlackTree_Set_Fun(t *testing.T) {
	//GetOrPutFunc lock or unlock
	gtest.C(t, func(t *gtest.T) {
		m := g.NewRedBlackTree[string, int](comparators.ComparatorString)
		t.Assert(m.GetOrPutFunc("fun", getValue), 3)
		t.Assert(m.GetOrPutFunc("fun", getValue), 3)
		t.Assert(m.GetOrPutFunc("funlock", getValue), 3)
		t.Assert(m.GetOrPutFunc("funlock", getValue), 3)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
	})
	//PutIfAbsentFunc lock or unlock
	gtest.C(t, func(t *gtest.T) {
		m := g.NewRedBlackTree[string, int](comparators.ComparatorString)
		t.Assert(m.PutIfAbsentFunc("fun", getValue), true)
		t.Assert(m.PutIfAbsentFunc("fun", getValue), false)
		t.Assert(m.PutIfAbsentFunc("funlock", getValue), true)
		t.Assert(m.PutIfAbsentFunc("funlock", getValue), false)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
	})

}

func Test_RedBlackTree_Batch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewRedBlackTree[string, string](comparators.ComparatorString)
		m.Puts(map[string]string{"1": "1", "key1": "val1", "key2": "val2", "key3": "val3"})
		t.Assert(m.Map(), map[string]string{"1": "1", "key1": "val1", "key2": "val2", "key3": "val3"})
		m.Removes([]string{"key1", "1"})
		t.Assert(m.Map(), map[string]string{"key2": "val2", "key3": "val3"})
	})
}

func Test_RedBlackTree_Iterator(t *testing.T) {
	keys := []string{"1", "key1", "key2", "key3", "key4"}
	keyLen := len(keys)
	index := 0

	expect := map[string]string{"key4": "val4", "1": "1", "key1": "val1", "key2": "val2", "key3": "val3"}
	m := g.NewRedBlackTreeFrom[string, string](comparators.ComparatorString, expect)

	gtest.C(t, func(t *gtest.T) {

		m.ForEach(func(k string, v string) bool {
			t.Assert(k, keys[index])
			index++
			t.Assert(expect[k], v)
			return true
		})

		m.ForEachDesc(func(k string, v string) bool {
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
		m.ForEach(func(k string, v string) bool {
			i++
			return true
		})
		m.ForEach(func(k string, v string) bool {
			j++
			return false
		})
		t.Assert(i, keyLen)
		t.Assert(j, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		i := 0
		j := 0
		m.ForEachDesc(func(k string, v string) bool {
			i++
			return true
		})
		m.ForEachDesc(func(k string, v string) bool {
			j++
			return false
		})
		t.Assert(i, keyLen)
		t.Assert(j, 1)
	})
}

func Test_RedBlackTree_IteratorFrom(t *testing.T) {
	m := make(map[int]int)
	for i := 1; i <= 10; i++ {
		if i == 2 || i == 8 {
			continue
		}
		m[i] = i * 10
	}
	tree := g.NewRedBlackTreeFrom[int, int](comparators.ComparatorInt, m)

	gtest.C(t, func(t *gtest.T) {
		n1 := 5
		tree.IteratorFrom(5, true, func(key, value int) bool {
			t.Assert(n1, key)
			t.Assert(n1*10, value)
			n1++
			if n1 == 2 || n1 == 8 {
				n1++
			}
			return true
		})

		n2 := 5
		tree.IteratorAscFrom(5, true, func(key, value int) bool {
			t.Assert(n2, key)
			t.Assert(n2*10, value)
			n2++
			if n2 == 2 || n2 == 8 {
				n2++
			}
			return true
		})

		n3 := 6
		tree.IteratorAscFrom(5, false, func(key, value int) bool {
			t.Assert(n3, key)
			t.Assert(n3*10, value)
			n3++
			if n3 == 2 || n3 == 8 {
				n3++
			}
			return true
		})

		n4 := 3
		tree.IteratorAscFrom(2, true, func(key, value int) bool {
			t.Assert(n4, key)
			t.Assert(n4*10, value)
			n4++
			if n4 == 2 || n4 == 8 {
				n4++
			}
			return true
		})

		n5 := 3
		tree.IteratorAscFrom(2, false, func(key, value int) bool {
			t.Assert(n5, key)
			t.Assert(n5*10, value)
			n5++
			if n5 == 2 || n5 == 8 {
				n5++
			}
			return true
		})

		n6 := 5
		tree.IteratorDescFrom(5, true, func(key, value int) bool {
			t.Assert(n6, key)
			t.Assert(n6*10, value)
			n6--
			if n6 == 2 || n6 == 8 {
				n6--
			}
			return true
		})

		n7 := 4
		tree.IteratorDescFrom(5, false, func(key, value int) bool {
			t.Assert(n7, key)
			t.Assert(n7*10, value)
			n7--
			if n7 == 2 || n7 == 8 {
				n7--
			}
			return true
		})

		n8 := 7
		tree.IteratorDescFrom(8, true, func(key, value int) bool {
			t.Assert(n8, key)
			t.Assert(n8*10, value)
			n8--
			if n8 == 2 || n8 == 8 {
				n8--
			}
			return true
		})

		n9 := 7
		tree.IteratorDescFrom(8, false, func(key, value int) bool {
			t.Assert(n9, key)
			t.Assert(n9*10, value)
			n9--
			if n9 == 2 || n8 == 8 {
				n9--
			}
			return true
		})
	})
}

func Test_RedBlackTree_SubMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string]int)
		for i := 0; i < 10; i++ {
			m["key"+gconv.String(i)] = i * 10
		}
		tree := g.NewRedBlackTreeFrom(comparators.ComparatorString, m)
		// both key exists in map
		t.Assert(tree.SubMap("key5", true, "key7", true).Values(), []int{50, 60, 70})
		t.Assert(tree.SubMap("key5", false, "key7", true).Values(), []int{60, 70})
		t.Assert(tree.SubMap("key5", true, "key7", false).Values(), []int{50, 60})
		t.Assert(tree.SubMap("key5", false, "key7", false).Values(), []int{60})
		// only fromKey exists in map
		t.Assert(tree.SubMap("key5.1", true, "key7", true).Values(), []int{60, 70})
		t.Assert(tree.SubMap("key5.1", false, "key7", true).Values(), []int{60, 70})
		t.Assert(tree.SubMap("key5.1", true, "key7", false).Values(), []int{60})
		t.Assert(tree.SubMap("key5.1", false, "key7", false).Values(), []int{60})
		// both key do not exist in map
		t.Assert(tree.SubMap("key5.1", true, "key7.1", true).Values(), []int{60, 70})
		t.Assert(tree.SubMap("key5.1", false, "key7.1", true).Values(), []int{60, 70})
		t.Assert(tree.SubMap("key5.1", true, "key7.1", false).Values(), []int{60, 70})
		t.Assert(tree.SubMap("key5.1", false, "key7.1", false).Values(), []int{60, 70})
		// fromKey out of upper bound
		t.Assert(tree.SubMap("zz", false, "key7.1", false).Values(), []int{})
		// fromKey out of lower bound
		t.Assert(tree.SubMap("aa", false, "key0.1", false).Values(), []int{0})
		// both key out of lower bound
		t.Assert(tree.SubMap("aa", false, "bb", false).Values(), []int{})
		// both key out of lower bound
		t.Assert(tree.SubMap("bb", false, "aa", false).Values(), []int{})
		// both key out of upper bound
		t.Assert(tree.SubMap("yy", false, "zz", false).Values(), []int{})
		// both key out of upper bound
		t.Assert(tree.SubMap("zz", false, "yy", false).Values(), []int{})
		// toKey out of upper bound
		t.Assert(tree.SubMap("key9", true, "zz", false).Values(), []int{90})
		// fromKey out of lower bound and toKey out of upper bound
		t.Assert(tree.SubMap("aa", true, "zz", false).Values(), []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90})
	})
}

func Test_RedBlackTree_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		//clone 方法是深克隆
		m := g.NewRedBlackTreeFrom[string, string](comparators.ComparatorString, map[string]string{"1": "1", "key1": "val1"})
		m_clone := m.Clone()
		m.Remove("1")
		//修改原 map,clone 后的 map 不影响
		t.AssertIN(1, m_clone.Keys())

		m_clone.Remove("key1")
		//修改clone map,原 map 不影响
		t.AssertIN("key1", m.Keys())
	})
}

func Test_RedBlackTree_LRNode(t *testing.T) {
	expect := map[string]string{"key4": "val4", "key1": "val1", "key2": "val2", "key3": "val3"}
	//safe
	gtest.C(t, func(t *gtest.T) {
		m := g.NewRedBlackTreeFrom[string, string](comparators.ComparatorString, expect)
		t.Assert(m.Left().Key(), "key1")
		t.Assert(m.Right().Key(), "key4")
	})
	//unsafe
	gtest.C(t, func(t *gtest.T) {
		m := g.NewRedBlackTreeFrom[string, string](comparators.ComparatorString, expect, true)
		t.Assert(m.Left().Key(), "key1")
		t.Assert(m.Right().Key(), "key4")
	})
}

func Test_RedBlackTree_CeilingFloor(t *testing.T) {
	expect := map[int]string{
		20: "val20",
		6:  "val6",
		10: "val10",
		12: "val12",
		1:  "val1",
		15: "val15",
		19: "val19",
		8:  "val8",
		4:  "val4"}
	//found and eq
	gtest.C(t, func(t *gtest.T) {
		m := g.NewRedBlackTreeFrom[int, string](comparators.ComparatorInt, expect)
		c := m.CeilingEntry(8)
		t.Assert(c != nil, true)
		t.Assert(c.Value(), "val8")
		f := m.FloorEntry(20)
		t.Assert(f != nil, true)
		t.Assert(f.Value(), "val20")
	})
	//found and neq
	gtest.C(t, func(t *gtest.T) {
		m := g.NewRedBlackTreeFrom[int, string](comparators.ComparatorInt, expect)
		c := m.CeilingEntry(9)
		t.Assert(c != nil, true)
		t.Assert(c.Value(), "val10")
		f := m.FloorEntry(5)
		t.Assert(f != nil, true)
		t.Assert(f.Value(), "val4")
	})
	//nofound
	gtest.C(t, func(t *gtest.T) {
		m := g.NewRedBlackTreeFrom[int, string](comparators.ComparatorInt, expect)
		c := m.CeilingEntry(21)
		t.Assert(c, nil)
		f := m.FloorEntry(-1)
		t.Assert(f, nil)
	})
}

func Test_RedBlackTree_Remove(t *testing.T) {
	m := g.NewRedBlackTree[int, string](comparators.ComparatorInt)
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
