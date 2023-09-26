// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g_test

import (
	"fmt"

	"github.com/wesleywu/gcontainer/g"
)

func ExampleNewArrayList() {
	n := 10
	l := g.NewLinkedList[int]()
	for i := 0; i < n; i++ {
		l.PushBack(i)
	}

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.FrontAll())
	fmt.Println(l.BackAll())

	for i := 0; i < n; i++ {
		v, _ := l.PopFront()
		fmt.Print(v)
	}

	fmt.Println()
	fmt.Println(l.Len())
	v, _ := l.PopFront()
	fmt.Println(v)
	fmt.Println(l.Len())

	// Output:
	// 10
	// [0,1,2,3,4,5,6,7,8,9]
	// [0 1 2 3 4 5 6 7 8 9]
	// [9 8 7 6 5 4 3 2 1 0]
	// 0123456789
	// 0
	// 0
	// 0
}

func ExampleNewArrayListFrom() {
	n := 10
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 10, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.FrontAll())
	fmt.Println(l.BackAll())

	for i := 0; i < n; i++ {
		v, _ := l.PopFront()
		fmt.Print(v)
	}

	fmt.Println()
	fmt.Println(l.Len())

	// Output:
	// 10
	// [1,2,3,4,5,6,7,8,9,10]
	// [1 2 3 4 5 6 7 8 9 10]
	// [10 9 8 7 6 5 4 3 2 1]
	// 12345678910
	// 0
}

func ExampleLinkedList_PushFront() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.PushFront(0)

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 6
	// [0,1,2,3,4,5]
}

func ExampleLinkedList_PushBack() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.PushBack(6)

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 6
	// [1,2,3,4,5,6]
}

func ExampleLinkedList_PushFronts() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.PushFronts([]int{0, -1, -2, -3, -4})

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 10
	// [-4,-3,-2,-1,0,1,2,3,4,5]
}

func ExampleLinkedList_PushBacks() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.PushBacks([]int{6, 7, 8, 9, 10})

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 10
	// [1,2,3,4,5,6,7,8,9,10]
}

func ExampleLinkedList_PopBack() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	v, _ := l.PopBack()
	fmt.Println(v)
	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 5
	// 4
	// [1,2,3,4]
}

func ExampleLinkedList_PopFront() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	v, _ := l.PopFront()
	fmt.Println(v)
	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 1
	// 4
	// [2,3,4,5]
}

func ExampleLinkedList_PopBacks() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.PopBacks(2))
	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// [5 4]
	// 3
	// [1,2,3]
}

func ExampleLinkedList_PopFronts() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.PopFronts(2))
	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// [1 2]
	// 3
	// [3,4,5]
}

func ExampleLinkedList_PopBackAll() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.PopBackAll())
	fmt.Println(l.Len())

	// Output:
	// 5
	// [1,2,3,4,5]
	// [5 4 3 2 1]
	// 0
}

func ExampleLinkedList_PopFrontAll() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.PopFrontAll())
	fmt.Println(l.Len())

	// Output:
	// 5
	// [1,2,3,4,5]
	// [1 2 3 4 5]
	// 0
}

func ExampleLinkedList_FrontAll() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.FrontAll())

	// Output:
	// [1,2,3,4,5]
	// [1 2 3 4 5]
}

func ExampleLinkedList_BackAll() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.BackAll())

	// Output:
	// [1,2,3,4,5]
	// [5 4 3 2 1]
}

func ExampleLinkedList_FrontValue() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.FrontValue())

	// Output:
	// [1,2,3,4,5]
	// 1
}

func ExampleLinkedList_BackValue() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.BackValue())

	// Output:
	// [1,2,3,4,5]
	// 5
}

func ExampleLinkedList_Front() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Front().Value)
	fmt.Println(l)

	e := l.Front()
	l.InsertBefore(e, 0)
	l.InsertAfter(e, 9)

	fmt.Println(l)

	// Output:
	// 1
	// [1,2,3,4,5]
	// [0,1,9,2,3,4,5]
}

func ExampleLinkedList_Back() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Back().Value)
	fmt.Println(l)

	e := l.Back()
	l.InsertBefore(e, 9)
	l.InsertAfter(e, 6)

	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// [1,2,3,4,9,5,6]
}

func ExampleLinkedList_Len() {
	l := g.NewLinkedListFrom[int]([]int{1, 2, 3, 4, 5})

	fmt.Println(l.Len())

	// Output:
	// 5
}

func ExampleLinkedList_Size() {
	l := g.NewLinkedListFrom[int]([]int{1, 2, 3, 4, 5})

	fmt.Println(l.Size())

	// Output:
	// 5
}

func ExampleLinkedList_MoveBefore() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	// element of `l`
	e := l.PushBack(6)
	fmt.Println(l.Size())
	fmt.Println(l)

	l.MoveBefore(e, l.Front())

	fmt.Println(l.Size())
	fmt.Println(l)

	// not element of `l`
	e = &g.Element[int]{Value: 7}
	l.MoveBefore(e, l.Front())

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 6
	// [1,2,3,4,5,6]
	// 6
	// [6,1,2,3,4,5]
	// 6
	// [6,1,2,3,4,5]
}

func ExampleLinkedList_MoveAfter() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	// element of `l`
	e := l.PushFront(0)
	fmt.Println(l.Size())
	fmt.Println(l)

	l.MoveAfter(e, l.Back())

	fmt.Println(l.Size())
	fmt.Println(l)

	// not element of `l`
	e = &g.Element[int]{Value: -1}
	l.MoveAfter(e, l.Back())

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 6
	// [0,1,2,3,4,5]
	// 6
	// [1,2,3,4,5,0]
	// 6
	// [1,2,3,4,5,0]
}

func ExampleLinkedList_MoveToFront() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	// element of `l`
	l.MoveToFront(l.Back())

	fmt.Println(l.Size())
	fmt.Println(l)

	// not element of `l`
	e := &g.Element[int]{Value: 6}
	l.MoveToFront(e)

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 5
	// [5,1,2,3,4]
	// 5
	// [5,1,2,3,4]
}

func ExampleLinkedList_MoveToBack() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	// element of `l`
	l.MoveToBack(l.Front())

	fmt.Println(l.Size())
	fmt.Println(l)

	// not element of `l`
	e := &g.Element[int]{Value: 0}
	l.MoveToBack(e)

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 5
	// [2,3,4,5,1]
	// 5
	// [2,3,4,5,1]
}

func ExampleLinkedList_PushBackList() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	other := g.NewLinkedListFrom[int]([]int{6, 7, 8, 9, 10})

	fmt.Println(other.Size())
	fmt.Println(other)

	l.PushBackList(other)

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 5
	// [6,7,8,9,10]
	// 10
	// [1,2,3,4,5,6,7,8,9,10]
}

func ExampleLinkedList_PushFrontList() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	other := g.NewLinkedListFrom[int]([]int{-4, -3, -2, -1, 0})

	fmt.Println(other.Size())
	fmt.Println(other)

	l.PushFrontList(other)

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 5
	// [-4,-3,-2,-1,0]
	// 10
	// [-4,-3,-2,-1,0,1,2,3,4,5]
}

func ExampleLinkedList_InsertAfter() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.InsertAfter(l.Front(), 8)
	l.InsertAfter(l.Back(), 9)

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 7
	// [1,8,2,3,4,5,9]
}

func ExampleLinkedList_InsertBefore() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.InsertBefore(l.Front(), 8)
	l.InsertBefore(l.Back(), 9)

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 7
	// [8,1,2,3,4,9,5]
}

func ExampleLinkedList_Remove() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	fmt.Println(l.Remove(l.Front().Value))
	fmt.Println(l.Remove(l.Back().Value))

	fmt.Println(l.Len())
	fmt.Println(l)

	fmt.Println(l.Remove(l.Front().Value, l.Back().Value))

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// true
	// true
	// 3
	// [2,3,4]
	// true
	// 1
	// [3]
}

func ExampleLinkedList_RemoveAll() {
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.Clear()

	fmt.Println(l.Len())

	// Output:
	// 5
	// [1,2,3,4,5]
	// 0
}

func ExampleLinkedList_ForEachAsc() {
	// concurrent-safe list.
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 10, 1).Slice(), true)
	// iterate reading from head using ForEachAsc.
	l.ForEachAsc(func(e int) bool {
		fmt.Print(e)
		return true
	})

	// Output:
	// 12345678910
}

func ExampleLinkedList_ForEachDesc() {
	// concurrent-safe list.
	l := g.NewLinkedListFrom[int](g.NewArrayListRange(1, 10, 1).Slice(), true)
	// iterate reading from tail using ForEachDesc.
	l.ForEachDesc(func(e int) bool {
		fmt.Print(e)
		return true
	})
	// Output:
	// 10987654321
}

func ExampleLinkedList_Join() {
	var l g.LinkedList[string]
	l.PushBacks([]string{"a", "b", "c", "d"})

	fmt.Println(l.Join(","))

	// Output:
	// a,b,c,d
}
