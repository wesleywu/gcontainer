// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package g_test

import (
	"testing"
	"time"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/internal/gtest"
	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

type exampleKey struct {
	Key string `json:"key"`
}

type exampleElement struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func getAny() interface{} {
	return 123
}

func Test_AnyAnyMap_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m = g.NewHashMap[exampleKey, *exampleElement]()
		element1 := &exampleElement{
			Code:    1,
			Message: "m1",
		}
		element2 := &exampleElement{
			Code:    2,
			Message: "m2",
		}
		element3 := &exampleElement{
			Code:    3,
			Message: "m3",
		}

		m.Put(exampleKey{Key: "1"}, element1)

		t.Assert(m.Get(exampleKey{Key: "1"}), element1)
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrPut(exampleKey{Key: "2"}, element2), element2)
		t.Assert(m.PutIfAbsent(exampleKey{Key: "2"}, element2), false)

		t.Assert(m.PutIfAbsent(exampleKey{Key: "3"}, element3), true)

		val, removed := m.Remove(exampleKey{Key: "2"})
		t.Assert(val, element2)
		t.Assert(removed, true)
		t.Assert(m.ContainsKey(exampleKey{Key: "2"}), false)

		t.AssertIN(exampleKey{Key: "3"}, m.Keys())
		t.AssertIN(exampleKey{Key: "1"}, m.Keys())
		t.AssertIN(element3, m.Values())
		t.AssertIN(element1, m.Values())

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
	})
	gtest.C(t, func(t *gtest.T) {
		var m = g.NewHashMap[string, string]()
		m.Put("1", "1")

		t.Assert(m.Get("1"), "1")
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrPut("2", "2"), "2")
		t.Assert(m.PutIfAbsent("2", "2"), false)

		t.Assert(m.PutIfAbsent("3", "3"), true)

		val, removed := m.Remove("2")
		t.Assert(val, "2")
		t.Assert(removed, true)
		t.Assert(m.ContainsKey("2"), false)

		t.AssertIN(3, m.Keys())
		t.AssertIN(1, m.Keys())
		t.AssertIN(3, m.Values())
		t.AssertIN(1, m.Values())

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)
	})
}

func Test_AnyAnyMap_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMap[int, int]()
		m.Put(1, 1)

		t.Assert(m.Get(1), 1)
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrPut(2, 2), 2)
		t.Assert(m.PutIfAbsent(2, 2), false)

		t.Assert(m.PutIfAbsent(3, 3), true)

		val, removed := m.Remove(2)
		t.Assert(val, 2)
		t.Assert(removed, true)
		t.Assert(m.ContainsKey(2), false)

		t.AssertIN(3, m.Keys())
		t.AssertIN(1, m.Keys())
		t.AssertIN(3, m.Values())
		t.AssertIN(1, m.Values())

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)

		m2 := g.NewHashMapFrom[int, int](map[int]int{1: 1, 2: 2})
		t.Assert(m2.Map(), map[int]int{1: 1, 2: 2})
	})
}

func Test_AnyAnyMap_Set_Fun(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMap[int, any]()

		m.GetOrPutFunc(1, getAny)
		m.GetOrPutFunc(2, getAny)
		t.Assert(m.Get(1), 123)
		t.Assert(m.Get(2), 123)

		t.Assert(m.PutIfAbsentFunc(1, getAny), false)
		t.Assert(m.PutIfAbsentFunc(3, getAny), true)

		t.Assert(m.PutIfAbsentFunc(2, getAny), false)
		t.Assert(m.PutIfAbsentFunc(4, getAny), true)
	})

}

func Test_AnyAnyMap_Batch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMap[int, int]()

		m.Puts(map[int]int{1: 1, 2: 2, 3: 3})
		t.Assert(m.Map(), map[int]int{1: 1, 2: 2, 3: 3})
		m.Removes([]int{1, 2})
		t.Assert(m.Map(), map[int]int{3: 3})
	})
}

func Test_AnyAnyMap_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[int]int{1: 1, 2: 2}
		m := g.NewHashMapFrom[int, int](expect)
		m.ForEach(func(k int, v int) bool {
			t.Assert(expect[k], v)
			return true
		})
		// 断言返回值对遍历控制
		i := 0
		j := 0
		m.ForEach(func(k int, v int) bool {
			i++
			return true
		})
		m.ForEach(func(k int, v int) bool {
			j++
			return false
		})
		t.Assert(i, 2)
		t.Assert(j, 1)
	})
}

func Test_AnyAnyMap_Lock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := map[int]int{1: 1, 2: 2}
		m := g.NewHashMapFrom[int, int](expect)
		m.LockFunc(func(m map[int]int) {
			t.Assert(m, expect)
		})
		m.RLockFunc(func(m map[int]int) {
			t.Assert(m, expect)
		})
	})
}

func Test_AnyAnyMap_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// clone 方法是深克隆
		m := g.NewHashMapFrom[int, int](map[int]int{1: 1, 2: 2})

		m_clone := m.Clone()
		m.Remove(1)
		// 修改原 map,clone 后的 map 不影响
		t.AssertIN(1, m_clone.Keys())

		m_clone.Remove(2)
		// 修改clone map,原 map 不影响
		t.AssertIN(2, m.Keys())
	})
}

func Test_AnyAnyMap_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := g.NewHashMap[int, int]()
		m2 := g.NewHashMap[int, int]()
		m1.Put(1, 1)
		m2.Put(2, 2)
		m1.Merge(m2)
		t.Assert(m1.Map(), map[int]int{1: 1, 2: 2})
		m3 := g.NewHashMapFrom[int, int](nil)
		m3.Merge(m2)
		t.Assert(m3.Map(), m2.Map())
	})
}

func Test_AnyAnyMap_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMap[int, int]()
		m.Put(1, 0)
		m.Put(2, 2)
		t.Assert(m.Get(1), 0)
		t.Assert(m.Get(2), 2)
		data := m.Map()
		t.Assert(data[1], 0)
		t.Assert(data[2], 2)
		data[3] = 3
		t.Assert(m.Get(3), 0)
		m.Put(4, 4)
		t.Assert(data[4], 0)
	})
}

func Test_AnyAnyMap_FilterEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMap[int, int]()
		m.Put(1, 0)
		m.Put(2, 2)
		t.Assert(m.Get(1), 0)
		t.Assert(m.Get(2), 2)
		m.FilterEmpty()
		t.Assert(m.Get(1), 0)
		t.Assert(m.Get(2), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMap[string, time.Time]()
		m.Put("time1", time.Time{})
		m.Put("time2", time.Now())
		t.Assert(m.Get("time1"), time.Time{})
		m.FilterEmpty()
		t.Assert(m.Get("time1"), nil)
		t.AssertNE(m.Get("time2"), nil)
	})
}

func Test_AnyAnyMap_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		data := map[string]string{
			"k1": "v1",
			"k2": "v2",
		}
		m1 := g.NewHashMapFrom[string, string](data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(gconv.Map(data))
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})
	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		data := map[string]string{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(gconv.Map(data))
		t.AssertNil(err)

		m := g.NewHashMap[string, string]()
		err = json.UnmarshalUseNumber(b, m)
		t.AssertNil(err)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
	gtest.C(t, func(t *gtest.T) {
		data := map[string]string{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(gconv.Map(data))
		t.AssertNil(err)

		var m g.HashMap[string, string]
		err = json.UnmarshalUseNumber(b, &m)
		t.AssertNil(err)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
}

func Test_AnyAnyMap_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMapFrom[string, string](map[string]string{
			"k1": "v1",
			"k2": "v2",
		})
		t.Assert(m.Size(), 2)

		k1, v1 := m.Pop()
		t.AssertIN(k1, []string{"k1", "k2"})
		t.AssertIN(v1, []string{"v1", "v2"})
		t.Assert(m.Size(), 1)
		k2, v2 := m.Pop()
		t.AssertIN(k2, []string{"k1", "k2"})
		t.AssertIN(v2, []string{"v1", "v2"})
		t.Assert(m.Size(), 0)

		t.AssertNE(k1, k2)
		t.AssertNE(v1, v2)

		k3, v3 := m.Pop()
		t.AssertNil(k3)
		t.AssertNil(v3)
	})
}

func Test_AnyAnyMap_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMapFrom[string, string](map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		})
		t.Assert(m.Size(), 3)

		kArray := g.NewArrayList[string]()
		vArray := g.NewArrayList[string]()
		for k, v := range m.Pops(1) {
			t.AssertIN(k, []string{"k1", "k2", "k3"})
			t.AssertIN(v, []string{"v1", "v2", "v3"})
			kArray.Add(k)
			vArray.Add(v)
		}
		t.Assert(m.Size(), 2)
		for k, v := range m.Pops(2) {
			t.AssertIN(k, []string{"k1", "k2", "k3"})
			t.AssertIN(v, []string{"v1", "v2", "v3"})
			kArray.Add(k)
			vArray.Add(v)
		}
		t.Assert(m.Size(), 0)

		t.Assert(kArray.Unique().Size(), 3)
		t.Assert(vArray.Unique().Size(), 3)

		v := m.Pops(1)
		t.AssertNil(v)
		v = m.Pops(-1)
		t.AssertNil(v)
	})
}

func TestStrStrMap_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Map  *g.HashMap[string, string]
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"k1":"v1","k2":"v2"}`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("k1"), "v1")
		t.Assert(v.Map.Get("k2"), "v2")
	})
	// HashMap
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("k1"), "v1")
		t.Assert(v.Map.Get("k2"), "v2")
	})
}

func TestStrIntMap_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Map  *g.HashMap[string, int]
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"k1":1,"k2":2}`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("k1"), 1)
		t.Assert(v.Map.Get("k2"), 2)
	})
	// HashMap
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": map[string]int{
				"k1": 1,
				"k2": 2,
			},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("k1"), 1)
		t.Assert(v.Map.Get("k2"), 2)
	})
}

func TestIntStrMap_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Map  *g.HashMap[int, string]
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"1":"v1","2":"v2"}`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get(1), "v1")
		t.Assert(v.Map.Get(2), "v2")
	})
	// HashMap
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"map": map[int]string{
				1: "v1",
				2: "v2",
			},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get(1), "v1")
		t.Assert(v.Map.Get(2), "v2")
	})
}

func Test_AnyAnyMap_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.NewHashMapFrom[string, string](map[string]string{
			"k1": "v1",
			"k2": "v2",
		})
		t.Assert(m.Size(), 2)

		n := m.DeepCopy().(*g.HashMap[string, string])
		n.Put("k1", "val1")
		t.AssertNE(m.Get("k1"), n.Get("k1"))
	})
}
