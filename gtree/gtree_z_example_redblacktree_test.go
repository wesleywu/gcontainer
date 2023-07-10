// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gtree_test

import (
	"fmt"

	"github.com/wesleywu/gcontainer/gtree"
	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/utils/comparator"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

func ExampleRedBlackTree_SetComparator() {
	var tree gtree.RedBlackTree[string, string]
	tree.SetComparator(comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())
	fmt.Println(tree.Size())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// 6
}

func ExampleRedBlackTree_Clone() {
	b := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_Set() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())
	fmt.Println(tree.Size())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// 6
}

func ExampleRedBlackTree_Sets() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)

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

func ExampleRedBlackTree_Get() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Get("key1"))
	fmt.Println(tree.Get("key10"))

	// Output:
	// val1
	//
}

func ExampleRedBlackTree_GetOrSet() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetOrPut("key1", "newVal1"))
	fmt.Println(tree.GetOrPut("key6", "val6"))

	// Output:
	// val1
	// val6
}

func ExampleRedBlackTree_GetOrSetFunc() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_GetOrSetFuncLock() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_SetIfNotExist() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.PutIfAbsent("key1", "newVal1"))
	fmt.Println(tree.PutIfAbsent("key6", "val6"))

	// Output:
	// false
	// true
}

func ExampleRedBlackTree_SetIfNotExistFunc() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_SetIfNotExistFuncLock() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_Contains() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.ContainsKey("key1"))
	fmt.Println(tree.ContainsKey("key6"))

	// Output:
	// true
	// false
}

func ExampleRedBlackTree_Remove() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_Removes() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_IsEmpty() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)

	fmt.Println(tree.IsEmpty())

	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.IsEmpty())

	// Output:
	// true
	// false
}

func ExampleRedBlackTree_Size() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)

	fmt.Println(tree.Size())

	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Size())

	// Output:
	// 0
	// 6
}

func ExampleRedBlackTree_Keys() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 6; i > 0; i-- {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Keys())

	// Output:
	// [key1 key2 key3 key4 key5 key6]
}

func ExampleRedBlackTree_Values() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 6; i > 0; i-- {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Values())

	// Output:
	// [val1 val2 val3 val4 val5 val6]
}

func ExampleRedBlackTree_Map() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleRedBlackTree_MapStrAny() {
	tree := gtree.NewRedBlackTree[int, string](comparator.ComparatorInt)
	for i := 0; i < 6; i++ {
		tree.Put(1000+i, "val"+gconv.String(i))
	}

	fmt.Println(tree.MapStrAny())

	// Output:
	// map[1000:val0 1001:val1 1002:val2 1003:val3 1004:val4 1005:val5]
}

func ExampleRedBlackTree_Left() {
	tree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	for i := 1; i < 100; i++ {
		tree.Put(i, i)
	}
	fmt.Println(tree.Left().Key, tree.Left().Value)

	emptyTree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	fmt.Println(emptyTree.Left())

	// Output:
	// 1 1
	// <nil>
}

func ExampleRedBlackTree_Right() {
	tree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	for i := 1; i < 100; i++ {
		tree.Put(i, i)
	}
	fmt.Println(tree.Right().Key, tree.Right().Value)

	emptyTree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	fmt.Println(emptyTree.Left())

	// Output:
	// 99 99
	// <nil>
}

func ExampleRedBlackTree_Floor() {
	tree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	for i := 1; i < 100; i++ {
		if i != 50 {
			tree.Put(i, i)
		}
	}

	node, found := tree.Floor(95)
	if found {
		fmt.Println("Floor 95:", node.Key)
	}

	node, found = tree.Floor(50)
	if found {
		fmt.Println("Floor 50:", node.Key)
	}

	node, found = tree.Floor(100)
	if found {
		fmt.Println("Floor 100:", node.Key)
	}

	node, found = tree.Floor(0)
	if found {
		fmt.Println("Floor 0:", node.Key)
	}

	// Output:
	// Floor 95: 95
	// Floor 50: 49
	// Floor 100: 99
}

func ExampleRedBlackTree_Ceiling() {
	tree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	for i := 1; i < 100; i++ {
		if i != 50 {
			tree.Put(i, i)
		}
	}

	node, found := tree.Ceiling(1)
	if found {
		fmt.Println("Ceiling 1:", node.Key)
	}

	node, found = tree.Ceiling(50)
	if found {
		fmt.Println("Ceiling 50:", node.Key)
	}

	node, found = tree.Ceiling(100)
	if found {
		fmt.Println("Ceiling 100:", node.Key)
	}

	node, found = tree.Ceiling(-1)
	if found {
		fmt.Println("Ceiling -1:", node.Key)
	}

	// Output:
	// Ceiling 1: 1
	// Ceiling 50: 51
	// Ceiling -1: 1
}

func ExampleRedBlackTree_Lower() {
	tree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	for i := 1; i < 100; i++ {
		if i != 50 {
			tree.Put(i, i)
		}
	}

	node, found := tree.Lower(95)
	if found {
		fmt.Println("Lower 95:", node.Key)
	}

	node, found = tree.Lower(50)
	if found {
		fmt.Println("Lower 50:", node.Key)
	}

	node, found = tree.Lower(100)
	if found {
		fmt.Println("Lower 100:", node.Key)
	}

	node, found = tree.Lower(0)
	if found {
		fmt.Println("Lower 0:", node.Key)
	}

	// Output:
	// Lower 95: 94
	// Lower 50: 49
	// Lower 100: 99
}

func ExampleRedBlackTree_Higher() {
	tree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	for i := 1; i < 100; i++ {
		if i != 50 {
			tree.Put(i, i)
		}
	}

	node, found := tree.Higher(1)
	if found {
		fmt.Println("Higher 1:", node.Key)
	}

	node, found = tree.Higher(95)
	if found {
		fmt.Println("Higher 95:", node.Key)
	}

	node, found = tree.Higher(50)
	if found {
		fmt.Println("Higher 50:", node.Key)
	}

	node, found = tree.Higher(100)
	if found {
		fmt.Println("Higher 100:", node.Key)
	}

	node, found = tree.Higher(-1)
	if found {
		fmt.Println("Higher -1:", node.Key)
	}

	// Output:
	// Higher 1: 2
	// Higher 95: 96
	// Higher 50: 51
	// Higher -1: 1
}

func ExampleRedBlackTree_Iterator() {
	tree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	for i := 0; i < 10; i++ {
		tree.Put(i, 10-i)
	}

	var totalKey, totalValue int
	tree.Iterator(func(key, value int) bool {
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

func ExampleRedBlackTree_IteratorFrom() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewRedBlackTreeFrom[int, int](comparator.ComparatorInt, m)

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

func ExampleRedBlackTree_IteratorAsc() {
	tree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	for i := 0; i < 10; i++ {
		tree.Put(i, 10-i)
	}

	tree.IteratorAsc(func(key, value int) bool {
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

func ExampleRedBlackTree_IteratorAscFrom_inclusive() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		if i == 3 {
			continue
		}
		m[i] = i * 10
	}
	tree := gtree.NewRedBlackTreeFrom(comparator.ComparatorInt, m)

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

func ExampleRedBlackTree_IteratorAscFrom_nonInclusive() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		if i == 3 {
			continue
		}
		m[i] = i * 10
	}
	tree := gtree.NewRedBlackTreeFrom(comparator.ComparatorInt, m)

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

func ExampleRedBlackTree_IteratorDesc() {
	tree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	for i := 0; i < 10; i++ {
		tree.Put(i, 10-i)
	}

	tree.IteratorDesc(func(key, value int) bool {
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

func ExampleRedBlackTree_IteratorDescFrom_inclusive() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewRedBlackTreeFrom(comparator.ComparatorInt, m)

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

func ExampleRedBlackTree_IteratorDescFrom_nonInclusive() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewRedBlackTreeFrom(comparator.ComparatorInt, m)

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

func ExampleRedBlackTree_SubMap() {
	m := make(map[string]int)
	for i := 0; i < 10; i++ {
		m["key"+gconv.String(i)] = i * 10
	}
	tree := gtree.NewRedBlackTreeFrom(comparator.ComparatorString, m)

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

func ExampleRedBlackTree_Clear() {
	tree := gtree.NewRedBlackTree[int, string](comparator.ComparatorInt)
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

func ExampleRedBlackTree_Replace() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_String() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_Print() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_Search() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Search("key0"))
	fmt.Println(tree.Search("key6"))

	// Output:
	// val0 true
	//  false
}

func ExampleRedBlackTree_Flip() {
	tree := gtree.NewRedBlackTree[int, int](comparator.ComparatorInt)
	for i := 1; i < 6; i++ {
		tree.Put(i, i*10)
	}

	fmt.Println("Before Flip", tree.Map())

	treeFlipped := tree.Flip(comparator.ComparatorInt)

	fmt.Println("After Flip", treeFlipped.Map())

	// Output:
	// Before Flip map[1:10 2:20 3:30 4:40 5:50]
	// After Flip map[10:1 20:2 30:3 40:4 50:5]
}

func ExampleRedBlackTree_MarshalJSON() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
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

func ExampleRedBlackTree_UnmarshalJSON() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}
	bytes, err := json.Marshal(tree)

	otherTree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)
	err = json.Unmarshal(bytes, &otherTree)
	if err == nil {
		fmt.Println(otherTree.Map())
	}

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleRedBlackTree_UnmarshalValue() {
	tree := gtree.NewRedBlackTree[string, string](comparator.ComparatorString)

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
