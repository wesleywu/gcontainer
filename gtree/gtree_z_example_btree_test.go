// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/Agogf/gf.

package gtree_test

import (
	"fmt"

	"github.com/wesleywu/gcontainer/gtree"
	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/utils/comparator"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

func ExampleBTree_Clone() {
	b := gtree.NewBTree[string, string](3, comparator.ComparatorString)
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

func ExampleBTree_Put() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())
	fmt.Println(tree.Size())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
	// 6
}

func ExampleBTree_Puts() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)

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

func ExampleBTree_Get() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Get("key1"))
	fmt.Println(tree.Get("key10"))

	// Output:
	// val1
	//
}

func ExampleBTree_GetOrPut() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.GetOrPut("key1", "newVal1"))
	fmt.Println(tree.GetOrPut("key6", "val6"))

	// Output:
	// val1
	// val6
}

func ExampleBTree_GetOrPutFunc() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
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

func ExampleBTree_PutIfAbsent() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.PutIfAbsent("key1", "newVal1"))
	fmt.Println(tree.PutIfAbsent("key6", "val6"))

	// Output:
	// false
	// true
}

func ExampleBTree_PutIfAbsentFunc() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
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

func ExampleBTree_ContainsKey() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.ContainsKey("key1"))
	fmt.Println(tree.ContainsKey("key6"))

	// Output:
	// true
	// false
}

func ExampleBTree_Remove() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
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

func ExampleBTree_Removes() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
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

func ExampleBTree_IsEmpty() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)

	fmt.Println(tree.IsEmpty())

	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.IsEmpty())

	// Output:
	// true
	// false
}

func ExampleBTree_Size() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)

	fmt.Println(tree.Size())

	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Size())

	// Output:
	// 0
	// 6
}

func ExampleBTree_Keys() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 6; i > 0; i-- {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Keys())

	// Output:
	// [key1 key2 key3 key4 key5 key6]
}

func ExampleBTree_Values() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 6; i > 0; i-- {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Values())

	// Output:
	// [val1 val2 val3 val4 val5 val6]
}

func ExampleBTree_Map() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Map())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleBTree_MapStrAny() {
	tree := gtree.NewBTree[int, string](3, comparator.ComparatorInt)
	for i := 0; i < 6; i++ {
		tree.Put(1000+i, "val"+gconv.String(i))
	}

	fmt.Println(tree.MapStrAny())

	// Output:
	// map[1000:val0 1001:val1 1002:val2 1003:val3 1004:val4 1005:val5]
}

func ExampleBTree_Clear() {
	tree := gtree.NewBTree[int, string](3, comparator.ComparatorInt)
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

func ExampleBTree_Replace() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
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

func ExampleBTree_Height() {
	tree := gtree.NewBTree[int, int](3, comparator.ComparatorInt)
	for i := 0; i < 100; i++ {
		tree.Put(i, i)
	}
	fmt.Println(tree.Height())

	// Output:
	// 6
}

func ExampleBTree_Left() {
	tree := gtree.NewBTree[int, int](3, comparator.ComparatorInt)
	for i := 1; i < 100; i++ {
		tree.Put(i, i)
	}
	fmt.Println(tree.Left().Key, tree.Left().Value)

	emptyTree := gtree.NewBTree[int, int](3, comparator.ComparatorInt)
	fmt.Println(emptyTree.Left())

	// Output:
	// 1 1
	// <nil>
}

func ExampleBTree_Right() {
	tree := gtree.NewBTree[int, int](3, comparator.ComparatorInt)
	for i := 1; i < 100; i++ {
		tree.Put(i, i)
	}
	fmt.Println(tree.Right().Key, tree.Right().Value)

	emptyTree := gtree.NewBTree[int, int](3, comparator.ComparatorInt)
	fmt.Println(emptyTree.Left())

	// Output:
	// 99 99
	// <nil>
}

func ExampleBTree_String() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree)

	// Output:
	// key0
	// key1
	//     key2
	// key3
	//     key4
	//     key5
}

func ExampleBTree_Search() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(tree.Search("key0"))
	fmt.Println(tree.Search("key6"))

	// Output:
	// val0 true
	//  false
}

func ExampleBTree_Print() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	tree.Print()

	// Output:
	// key0
	// key1
	//     key2
	// key3
	//     key4
	//     key5
}

func ExampleBTree_Iterator() {
	tree := gtree.NewBTree[int, int](3, comparator.ComparatorInt)
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

func ExampleBTree_IteratorFrom() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewBTreeFrom(3, comparator.ComparatorInt, m)

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

func ExampleBTree_IteratorAsc() {
	tree := gtree.NewBTree[int, int](3, comparator.ComparatorInt)
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

func ExampleBTree_IteratorAscFrom_Normal() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewBTreeFrom(3, comparator.ComparatorInt, m)

	tree.IteratorAscFrom(1, true, func(key, value int) bool {
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

func ExampleBTree_IteratorAscFrom_NoExistKey() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewBTreeFrom(3, comparator.ComparatorInt, m)

	tree.IteratorAscFrom(0, true, func(key, value int) bool {
		fmt.Println("key:", key, ", value:", value)
		return true
	})

	// Output:
}

func ExampleBTree_IteratorAscFrom_NoExistKeyAndMatchFalse() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewBTreeFrom(3, comparator.ComparatorInt, m)

	tree.IteratorAscFrom(0, false, func(key, value int) bool {
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

func ExampleBTree_IteratorDesc() {
	tree := gtree.NewBTree[int, int](3, comparator.ComparatorInt)
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

func ExampleBTree_IteratorDescFrom() {
	m := make(map[int]int)
	for i := 1; i <= 5; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewBTreeFrom(3, comparator.ComparatorInt, m)

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

func ExampleBTree_MarshalJSON() {
	tree := gtree.NewBTree[string, string](3, comparator.ComparatorString)
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
