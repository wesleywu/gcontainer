// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package g_test

import (
	"fmt"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/utils/comparators"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

func ExampleTreeMap_SetComparator() {
	var tree g.TreeMap[string, string]
	tree.SetComparator(comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())
	fmt.Println(tree.Size())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// 6
}

func ExampleTreeMap_Clone() {
	b := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		b.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	tree := b.Clone()

	fmt.Println(tree.Map())
	fmt.Println(tree.Size())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// 6
}

func ExampleTreeMap_Put() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())
	fmt.Println(tree.Size())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// 6
}

func ExampleTreeMap_Puts() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)

	tree.Puts(map[string]string{
		"key1": "val1",
		"key2": "val2",
	})

	fmt.Println(tree.Map())
	fmt.Println(tree.Size())

	// Output:
	// map[key1:val1 key2:val2]
	// 2
}

func ExampleTreeMap_Get() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Get("key1"))
	fmt.Println(tree.Get("key10"))

	// Output:
	// val1
	//
}

func ExampleTreeMap_GetOrPut() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetOrPut("key1", "newVal1"))
	fmt.Println(tree.GetOrPut("key6", "val6"))

	// Output:
	// val1
	// val6
}

func ExampleTreeMap_GetOrPutFunc() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetOrPutFunc("key1", func() string {
		return "newVal1"
	}))
	fmt.Println(tree.GetOrPutFunc("key6", func() string {
		return "val6"
	}))

	// Output:
	// val1
	// val6
}

func ExampleTreeMap_PutIfAbsent() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.PutIfAbsent("key1", "newVal1"))
	fmt.Println(tree.PutIfAbsent("key6", "val6"))

	// Output:
	// false
	// true
}

func ExampleTreeMap_PutIfAbsentFunc() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.PutIfAbsentFunc("key1", func() string {
		return "newVal1"
	}))
	fmt.Println(tree.PutIfAbsentFunc("key6", func() string {
		return "val6"
	}))

	// Output:
	// false
	// true
}

func ExampleTreeMap_ContainsKey() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.ContainsKey("key1"))
	fmt.Println(tree.ContainsKey("key6"))

	// Output:
	// true
	// false
}

func ExampleTreeMap_Remove() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Remove("key1"))
	fmt.Println(tree.Remove("key6"))
	fmt.Println(tree.Map())

	// Output:
	// val1 true
	//  false
	// map[key0:val0 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleTreeMap_Removes() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	removeKeys := make([]string, 2)
	removeKeys = append(removeKeys, "key1")
	removeKeys = append(removeKeys, "key6")

	tree.Removes(removeKeys)

	fmt.Println(tree.Map())

	// Output:
	// map[key0:val0 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleTreeMap_IsEmpty() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)

	fmt.Println(tree.IsEmpty())

	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.IsEmpty())

	// Output:
	// true
	// false
}

func ExampleTreeMap_Size() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)

	fmt.Println(tree.Size())

	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Size())

	// Output:
	// 0
	// 6
}

func ExampleTreeMap_Keys() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 6; i > 0; i-- {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Keys())

	// Output:
	// [key1 key2 key3 key4 key5 key6]
}

func ExampleTreeMap_Values() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 6; i > 0; i-- {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Values())

	// Output:
	// [val1 val2 val3 val4 val5 val6]
}

func ExampleTreeMap_Map() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleTreeMap_MapStrAny() {
	tree := g.NewTreeMap[int, string](comparators.ComparatorInt)
	for i := 0; i < 6; i++ {
		tree.Put(1000+i, "val"+gconv.String(i))
	}

	fmt.Println(tree.MapStrAny())

	// Output:
	// map[1000:val0 1001:val1 1002:val2 1003:val3 1004:val4 1005:val5]
}

func ExampleTreeMap_Left() {
	tree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	for i := 1; i < 100; i++ {
		tree.Put(i, i)
	}
	fmt.Println(tree.Left().Key(), tree.Left().Value())

	emptyTree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	fmt.Println(emptyTree.Left())

	// Output:
	// 1 1
	// <nil>
}

func ExampleTreeMap_Right() {
	tree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	for i := 1; i < 100; i++ {
		tree.Put(i, i)
	}
	fmt.Println(tree.Right().Key(), tree.Right().Value())

	emptyTree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	fmt.Println(emptyTree.Left())

	// Output:
	// 99 99
	// <nil>
}

func ExampleTreeMap_FloorEntry() {
	tree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	for i := 1; i < 100; i++ {
		if i != 50 {
			tree.Put(i, i)
		}
	}

	node := tree.FloorEntry(95)
	if node != nil {
		fmt.Println("FloorEntry 95:", node.Key())
	}

	node = tree.FloorEntry(50)
	if node != nil {
		fmt.Println("FloorEntry 50:", node.Key())
	}

	node = tree.FloorEntry(100)
	if node != nil {
		fmt.Println("FloorEntry 100:", node.Key())
	}

	node = tree.FloorEntry(0)
	if node != nil {
		fmt.Println("FloorEntry 0:", node.Key())
	}

	// Output:
	// FloorEntry 95: 95
	// FloorEntry 50: 49
	// FloorEntry 100: 99
}

func ExampleTreeMap_CeilingEntry() {
	tree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	for i := 1; i < 100; i++ {
		if i != 50 {
			tree.Put(i, i)
		}
	}

	node := tree.CeilingEntry(1)
	if node != nil {
		fmt.Println("CeilingEntry 1:", node.Key())
	}

	node = tree.CeilingEntry(50)
	if node != nil {
		fmt.Println("CeilingEntry 50:", node.Key())
	}

	node = tree.CeilingEntry(100)
	if node != nil {
		fmt.Println("CeilingEntry 100:", node.Key())
	}

	node = tree.CeilingEntry(-1)
	if node != nil {
		fmt.Println("CeilingEntry -1:", node.Key())
	}

	// Output:
	// CeilingEntry 1: 1
	// CeilingEntry 50: 51
	// CeilingEntry -1: 1
}

func ExampleTreeMap_LowerEntry() {
	tree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	for i := 1; i < 100; i++ {
		if i != 50 {
			tree.Put(i, i)
		}
	}

	node := tree.LowerEntry(95)
	if node != nil {
		fmt.Println("LowerEntry 95:", node.Key())
	}

	node = tree.LowerEntry(50)
	if node != nil {
		fmt.Println("LowerEntry 50:", node.Key())
	}

	node = tree.LowerEntry(100)
	if node != nil {
		fmt.Println("LowerEntry 100:", node.Key())
	}

	node = tree.LowerEntry(0)
	if node != nil {
		fmt.Println("LowerEntry 0:", node.Key())
	}

	// Output:
	// LowerEntry 95: 94
	// LowerEntry 50: 49
	// LowerEntry 100: 99
}

func ExampleTreeMap_HigherEntry() {
	tree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	for i := 1; i < 100; i++ {
		if i != 50 {
			tree.Put(i, i)
		}
	}

	node := tree.HigherEntry(1)
	if node != nil {
		fmt.Println("HigherEntry 1:", node.Key())
	}

	node = tree.HigherEntry(95)
	if node != nil {
		fmt.Println("HigherEntry 95:", node.Key())
	}

	node = tree.HigherEntry(50)
	if node != nil {
		fmt.Println("HigherEntry 50:", node.Key())
	}

	node = tree.HigherEntry(100)
	if node != nil {
		fmt.Println("HigherEntry 100:", node.Key())
	}

	node = tree.HigherEntry(-1)
	if node != nil {
		fmt.Println("HigherEntry -1:", node.Key())
	}

	// Output:
	// HigherEntry 1: 2
	// HigherEntry 95: 96
	// HigherEntry 50: 51
	// HigherEntry -1: 1
}

func ExampleTreeMap_ForEach() {
	tree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	for i := 0; i < 10; i++ {
		tree.Put(i, 10-i)
	}

	var totalKey, totalValue int
	tree.ForEach(func(key, value int) bool {
		totalKey += key
		totalValue += value

		return totalValue < 20
	})

	fmt.Println("totalKey:", totalKey)
	fmt.Println("totalValue:", totalValue)

	// Output:
	// totalKey: 3
	// totalValue: 27
}

func ExampleTreeMap_IteratorFrom() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := g.NewTreeMapFrom[int, int](comparators.ComparatorInt, m)

	tree.IteratorFrom(1, true, func(key, value int) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 1 , value: 10
	// key: 2 , value: 20
	// key: 3 , value: 30
	// key: 4 , value: 40
	// key: 5 , value: 50
}

func ExampleTreeMap_ForEachAsc() {
	tree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	for i := 0; i < 10; i++ {
		tree.Put(i, 10-i)
	}

	tree.ForEachAsc(func(key, value int) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 0 , value: 10
	// key: 1 , value: 9
	// key: 2 , value: 8
	// key: 3 , value: 7
	// key: 4 , value: 6
	// key: 5 , value: 5
	// key: 6 , value: 4
	// key: 7 , value: 3
	// key: 8 , value: 2
	// key: 9 , value: 1
}

func ExampleTreeMap_IteratorAscFrom_inclusive() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		if i == 3 {
			continue
		}
		m[i] = i * 10
	}
	tree := g.NewTreeMapFrom(comparators.ComparatorInt, m)

	tree.IteratorAscFrom(1, true, func(key, value int) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	tree.IteratorAscFrom(3, true, func(key, value int) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 1 , value: 10
	// key: 2 , value: 20
	// key: 4 , value: 40
	// key: 5 , value: 50
	// key: 4 , value: 40
	// key: 5 , value: 50
}

func ExampleTreeMap_IteratorAscFrom_nonInclusive() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		if i == 3 {
			continue
		}
		m[i] = i * 10
	}
	tree := g.NewTreeMapFrom(comparators.ComparatorInt, m)

	tree.IteratorAscFrom(1, false, func(key, value int) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	tree.IteratorAscFrom(3, false, func(key, value int) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})
	// Output:
	// key: 2 , value: 20
	// key: 4 , value: 40
	// key: 5 , value: 50
	// key: 4 , value: 40
	// key: 5 , value: 50
}

func ExampleTreeMap_ForEachDesc() {
	tree := g.NewTreeMap[int, int](comparators.ComparatorInt)
	for i := 0; i < 10; i++ {
		tree.Put(i, 10-i)
	}

	tree.ForEachDesc(func(key, value int) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 9 , value: 1
	// key: 8 , value: 2
	// key: 7 , value: 3
	// key: 6 , value: 4
	// key: 5 , value: 5
	// key: 4 , value: 6
	// key: 3 , value: 7
	// key: 2 , value: 8
	// key: 1 , value: 9
	// key: 0 , value: 10
}

func ExampleTreeMap_IteratorDescFrom_inclusive() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := g.NewTreeMapFrom(comparators.ComparatorInt, m)

	tree.IteratorDescFrom(5, true, func(key, value int) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 5 , value: 50
	// key: 4 , value: 40
	// key: 3 , value: 30
	// key: 2 , value: 20
	// key: 1 , value: 10
}

func ExampleTreeMap_IteratorDescFrom_nonInclusive() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := g.NewTreeMapFrom(comparators.ComparatorInt, m)

	tree.IteratorDescFrom(5, false, func(key, value int) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
	// key: 4 , value: 40
	// key: 3 , value: 30
	// key: 2 , value: 20
	// key: 1 , value: 10
}

func ExampleTreeMap_SubMap() {
	m := make(map[string]int)
	for i := 0; i < 10; i++ {
		m["key"+gconv.String(i)] = i * 10
	}
	tree := g.NewTreeMapFrom(comparators.ComparatorString, m)

	fmt.Println(tree.SubMap("key5", true, "key7", true).Values())
	fmt.Println(tree.SubMap("key5", false, "key7", true).Values())
	fmt.Println(tree.SubMap("key5", true, "key7", false).Values())
	fmt.Println(tree.SubMap("key5", false, "key7", false).Values())
	fmt.Println(tree.SubMap("key5.1", true, "key7", true).Values())
	fmt.Println(tree.SubMap("key5.1", false, "key7", true).Values())
	fmt.Println(tree.SubMap("key5.1", true, "key7", false).Values())
	fmt.Println(tree.SubMap("key5.1", false, "key7", false).Values())
	fmt.Println(tree.SubMap("key5.1", true, "key7.1", true).Values())
	fmt.Println(tree.SubMap("key5.1", false, "key7.1", true).Values())
	fmt.Println(tree.SubMap("key5.1", true, "key7.1", false).Values())
	fmt.Println(tree.SubMap("key5.1", false, "key7.1", false).Values())
	fmt.Println(tree.SubMap("key9.1", false, "key7.1", false).Values())
	fmt.Println(tree.SubMap("key9.1", false, "key9.1", false).Values())
	fmt.Println(tree.SubMap("aa", false, "key0.1", false).Values())
	fmt.Println(tree.SubMap("aa", false, "bb", false).Values())
	fmt.Println(tree.SubMap("bb", false, "aa", false).Values())
	fmt.Println(tree.SubMap("yy", false, "zz", false).Values())
	fmt.Println(tree.SubMap("zz", false, "yy", false).Values())
	fmt.Println(tree.SubMap("key9", true, "zz", false).Values())

	// Output:
	// [50 60 70]
	// [60 70]
	// [50 60]
	// [60]
	// [60 70]
	// [60 70]
	// [60]
	// [60]
	// [60 70]
	// [60 70]
	// [60 70]
	// [60 70]
	// []
	// []
	// [0]
	// []
	// []
	// []
	// []
	// [90]
}

func ExampleTreeMap_Clear() {
	tree := g.NewTreeMap[int, string](comparators.ComparatorInt)
	for i := 0; i < 6; i++ {
		tree.Put(1000+i, "val"+gconv.String(i))
	}
	fmt.Println(tree.Size())

	tree.Clear()
	fmt.Println(tree.Size())

	// Output:
	// 6
	// 0
}

func ExampleTreeMap_Replace() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())

	data := map[string]string{
		"newKey0": "newVal0",
		"newKey1": "newVal1",
		"newKey2": "newVal2",
	}

	tree.Replace(data)

	fmt.Println(tree.Map())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// map[newKey0:newVal0 newKey1:newVal1 newKey2:newVal2]
}

func ExampleTreeMap_String() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.String())

	// Output:
	// │           ┌── key5
	// │       ┌── key4
	// │   ┌── key3
	// │   │   └── key2
	// └── key1
	//     └── key0
}

func ExampleTreeMap_Print() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	tree.Print()

	// Output:
	// │           ┌── key5
	// │       ┌── key4
	// │   ┌── key3
	// │   │   └── key2
	// └── key1
	//     └── key0
}

func ExampleTreeMap_Search() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Search("key0"))
	fmt.Println(tree.Search("key6"))

	// Output:
	// val0 true
	//  false
}

func ExampleTreeMap_MarshalJSON() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	bytes, err := json.Marshal(tree)
	if err == nil {
		fmt.Println(gconv.String(bytes))
	}

	// Output:
	// {"key0":"val0","key1":"val1","key2":"val2","key3":"val3","key4":"val4","key5":"val5"}
}

func ExampleTreeMap_UnmarshalJSON() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}
	bytes, err := json.Marshal(tree)

	otherTree := g.NewTreeMap[string, string](comparators.ComparatorString)
	err = json.Unmarshal(bytes, &otherTree)
	if err == nil {
		fmt.Println(otherTree.Map())
	}

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleTreeMap_UnmarshalValue() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)

	type User struct {
		Uid   string
		Name  string
		Pass1 string
		Pass2 string
	}

	var (
		user = User{
			Uid:   "1",
			Name:  "john",
			Pass1: "123",
			Pass2: "456",
		}
	)
	if err := gconv.Scan(user, tree); err == nil {
		fmt.Printf("%#v", tree.Map())
	}

	// Output:
	// map[string]string{"Name":"john", "Pass1":"123", "Pass2":"456", "Uid":"1"}
}
