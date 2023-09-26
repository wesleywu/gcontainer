// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g_test

import (
	"testing"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/internal/gtest"
	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

func checkListLen[T comparable](t *gtest.T, l *g.LinkedList[T], len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkListPointers[T comparable](t *gtest.T, l *g.LinkedList[T], es []*g.Element[T]) {
	if !checkListLen(t, l, len(es)) {
		return
	}
	for i, e := 0, l.Front(); i < l.Len(); i, e = i+1, e.Next() {
		if e.Prev() != es[i].Prev() {
			t.Errorf("list[%d].Prev = %p, want %p", i, e.Prev(), es[i].Prev())
		}
		if e.Next() != es[i].Next() {
			t.Errorf("list[%d].Next = %p, want %p", i, e.Next(), es[i].Next())
		}
	}
}

func TestVar(t *testing.T) {
	var l g.LinkedList[int]
	l.PushFront(1)
	l.PushFront(2)
	if v, ok := l.PopBack(); v != 1 || !ok {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		// fmt.Println(v)
	}
	if v, ok := l.PopBack(); v != 2 || !ok {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		// fmt.Println(v)
	}
	if v, ok := l.PopBack(); ok {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		// fmt.Println(v)
	}
	l.PushBack(1)
	l.PushBack(2)
	if v, ok := l.PopFront(); v != 1 || !ok {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		// fmt.Println(v)
	}
	if v, ok := l.PopFront(); v != 2 || !ok {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		// fmt.Println(v)
	}
	if v, ok := l.PopFront(); ok {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		// fmt.Println(v)
	}
}

func TestBasic(t *testing.T) {
	l := g.NewLinkedList[int]()
	l.PushFront(1)
	l.PushFront(2)
	if v, ok := l.PopBack(); v != 1 || !ok {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		// fmt.Println(v)
	}
	if v, ok := l.PopBack(); v != 2 || !ok {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		// fmt.Println(v)
	}
	if v, ok := l.PopBack(); ok {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		// fmt.Println(v)
	}
	l.PushBack(1)
	l.PushBack(2)
	if v, ok := l.PopFront(); v != 1 || !ok {
		t.Errorf("EXPECT %v, GOT %v", 1, v)
	} else {
		// fmt.Println(v)
	}
	if v, ok := l.PopFront(); v != 2 || !ok {
		t.Errorf("EXPECT %v, GOT %v", 2, v)
	} else {
		// fmt.Println(v)
	}
	if v, ok := l.PopFront(); ok {
		t.Errorf("EXPECT %v, GOT %v", nil, v)
	} else {
		// fmt.Println(v)
	}
}

func TestList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		checkListPointers(t, l, []*g.Element[int]{})

		// Single element list
		e := l.PushFront(9)
		checkListPointers(t, l, []*g.Element[int]{e})
		l.MoveToFront(e)
		checkListPointers(t, l, []*g.Element[int]{e})
		l.MoveToBack(e)
		checkListPointers(t, l, []*g.Element[int]{e})
		l.Remove(e.Value)
		checkListPointers(t, l, []*g.Element[int]{})

		// Bigger list
		e2 := l.PushFront(2)
		e1 := l.PushFront(1)
		e3 := l.PushBack(3)
		e4 := l.PushBack(88)
		checkListPointers(t, l, []*g.Element[int]{e1, e2, e3, e4})

		l.Remove(e2.Value)
		checkListPointers(t, l, []*g.Element[int]{e1, e3, e4})

		l.MoveToFront(e3) // move from middle
		checkListPointers(t, l, []*g.Element[int]{e3, e1, e4})

		l.MoveToFront(e1)
		l.MoveToBack(e3) // move from middle
		checkListPointers(t, l, []*g.Element[int]{e1, e4, e3})

		l.MoveToFront(e3) // move from back
		checkListPointers(t, l, []*g.Element[int]{e3, e1, e4})
		l.MoveToFront(e3) // should be no-op
		checkListPointers(t, l, []*g.Element[int]{e3, e1, e4})

		l.MoveToBack(e3) // move from front
		checkListPointers(t, l, []*g.Element[int]{e1, e4, e3})
		l.MoveToBack(e3) // should be no-op
		checkListPointers(t, l, []*g.Element[int]{e1, e4, e3})

		e2 = l.InsertBefore(e1, 2) // insert before front
		checkListPointers(t, l, []*g.Element[int]{e2, e1, e4, e3})
		l.Remove(e2.Value)
		e2 = l.InsertBefore(e4, 2) // insert before middle
		checkListPointers(t, l, []*g.Element[int]{e1, e2, e4, e3})
		l.Remove(e2.Value)
		e2 = l.InsertBefore(e3, 2) // insert before back
		checkListPointers(t, l, []*g.Element[int]{e1, e4, e2, e3})
		l.Remove(e2.Value)

		e2 = l.InsertAfter(e1, 2) // insert after front
		checkListPointers(t, l, []*g.Element[int]{e1, e2, e4, e3})
		l.Remove(e2.Value)
		e2 = l.InsertAfter(e4, 2) // insert after middle
		checkListPointers(t, l, []*g.Element[int]{e1, e4, e2, e3})
		l.Remove(e2.Value)
		e2 = l.InsertAfter(e3, 2) // insert after back
		checkListPointers(t, l, []*g.Element[int]{e1, e4, e3, e2})
		l.Remove(e2.Value)

		// Check standard iteration.
		sum := 0
		for e := l.Front(); e != nil; e = e.Next() {
			i := e.Value
			sum += i
		}
		if sum != 92 {
			t.Errorf("sum over l = %d, want 4", sum)
		}

		// Clear all elements by iterating
		var next *g.Element[int]
		for e := l.Front(); e != nil; e = next {
			next = e.Next()
			l.Remove(e.Value)
		}
		checkListPointers(t, l, []*g.Element[int]{})
	})
}

func checkList[T comparable](t *gtest.T, l *g.LinkedList[T], es []T) {
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = e.Next() {

		if le := e.Value; le != es[i] {
			t.Errorf("elt[%d].value() = %v, want %v", i, le, es[i])
		}
		i++
	}
}

func TestExtending(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l1 := g.NewLinkedList[int]()
		l2 := g.NewLinkedList[int]()

		l1.PushBack(1)
		l1.PushBack(2)
		l1.PushBack(3)

		l2.PushBack(4)
		l2.PushBack(5)

		l3 := g.NewLinkedList[int]()
		l3.PushBackList(l1)
		checkList(t, l3, []int{1, 2, 3})
		l3.PushBackList(l2)
		checkList(t, l3, []int{1, 2, 3, 4, 5})

		l3 = g.NewLinkedList[int]()
		l3.PushFrontList(l2)
		checkList(t, l3, []int{4, 5})
		l3.PushFrontList(l1)
		checkList(t, l3, []int{1, 2, 3, 4, 5})

		checkList(t, l1, []int{1, 2, 3})
		checkList(t, l2, []int{4, 5})

		l3 = g.NewLinkedList[int]()
		l3.PushBackList(l1)
		checkList(t, l3, []int{1, 2, 3})
		l3.PushBackList(l3)
		checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

		l3 = g.NewLinkedList[int]()
		l3.PushFrontList(l1)
		checkList(t, l3, []int{1, 2, 3})
		l3.PushFrontList(l3)
		checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

		l3 = g.NewLinkedList[int]()
		l1.PushBackList(l3)
		checkList(t, l1, []int{1, 2, 3})
		l1.PushFrontList(l3)
		checkList(t, l1, []int{1, 2, 3})
	})
}

func TestRemove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		e1 := l.PushBack(1)
		e2 := l.PushBack(2)
		checkListPointers(t, l, []*g.Element[int]{e1, e2})
		// e := l.Front()
		// l.Remove(e)
		// checkListPointers(t, l, []*Element{e2})
		// l.Remove(e)
		// checkListPointers(t, l, []*Element{e2})
	})
}

func TestIssue4103(t *testing.T) {
	l1 := g.NewLinkedList[int]()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := g.NewLinkedList[int]()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.Front()
	l2.Remove(e.Value) // l2 should not change because e is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(e, 8)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestIssue6349(t *testing.T) {
	l := g.NewLinkedList[int]()
	l.PushBack(1)
	l.PushBack(2)

	e := l.Front()
	l.Remove(e.Value)
	if e.Value != 1 {
		t.Errorf("e.value = %d, want 1", e.Value)
	}
	// if e.Next() != nil {
	//    t.Errorf("e.Next() != nil")
	// }
	// if e.Prev() != nil {
	//    t.Errorf("e.Prev() != nil")
	// }
}

func TestMove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		e1 := l.PushBack(1)
		e2 := l.PushBack(2)
		e3 := l.PushBack(3)
		e4 := l.PushBack(4)

		l.MoveAfter(e3, e3)
		checkListPointers(t, l, []*g.Element[int]{e1, e2, e3, e4})
		l.MoveBefore(e2, e2)
		checkListPointers(t, l, []*g.Element[int]{e1, e2, e3, e4})

		l.MoveAfter(e3, e2)
		checkListPointers(t, l, []*g.Element[int]{e1, e2, e3, e4})
		l.MoveBefore(e2, e3)
		checkListPointers(t, l, []*g.Element[int]{e1, e2, e3, e4})

		l.MoveBefore(e2, e4)
		checkListPointers(t, l, []*g.Element[int]{e1, e3, e2, e4})
		e2, e3 = e3, e2

		l.MoveBefore(e4, e1)
		checkListPointers(t, l, []*g.Element[int]{e4, e1, e2, e3})
		e1, e2, e3, e4 = e4, e1, e2, e3

		l.MoveAfter(e4, e1)
		checkListPointers(t, l, []*g.Element[int]{e1, e4, e2, e3})
		e2, e3, e4 = e4, e2, e3

		l.MoveAfter(e2, e3)
		checkListPointers(t, l, []*g.Element[int]{e1, e3, e2, e4})
		e2, e3 = e3, e2
	})
}

// Test PushFront, PushBack, PushFrontList, PushBackList with uninitialized LinkedList
func TestZeroList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var l1 = g.NewLinkedList[int]()
		l1.PushFront(1)
		checkList(t, l1, []int{1})

		var l2 = g.NewLinkedList[int]()
		l2.PushBack(1)
		checkList(t, l2, []int{1})

		var l3 = g.NewLinkedList[int]()
		l3.PushFrontList(l1)
		checkList(t, l3, []int{1})

		var l4 = g.NewLinkedList[int]()
		l4.PushBackList(l2)
		checkList(t, l4, []int{1})
	})
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an element of l.
func TestInsertBeforeUnknownMark(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)
		l.InsertBefore(new(g.Element[int]), 1)
		checkList(t, l, []int{1, 2, 3})
	})
}

// Test that a list l is not modified when calling InsertAfter with a mark that is not an element of l.
func TestInsertAfterUnknownMark(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)
		l.InsertAfter(new(g.Element[int]), 1)
		checkList(t, l, []int{1, 2, 3})
	})
}

// Test that a list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an element of l.
func TestMoveUnknownMark(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l1 := g.NewLinkedList[int]()
		e1 := l1.PushBack(1)

		l2 := g.NewLinkedList[int]()
		e2 := l2.PushBack(2)

		l1.MoveAfter(e1, e2)
		checkList(t, l1, []int{1})
		checkList(t, l2, []int{2})

		l1.MoveBefore(e1, e2)
		checkList(t, l1, []int{1})
		checkList(t, l2, []int{2})
	})
}

func TestLinkedList_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		l.PushBack(1)
		l.Clear()
		checkList(t, l, []int{})
		l.PushBack(2)
		checkList(t, l, []int{2})
	})
}

func TestLinkedList_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{0, 1, 2, 3}
		array1 := g.NewLinkedListFrom(a1)
		array2 := array1.Clone().(*g.LinkedList[int])

		t.Assert(array2.Len(), 4)
		t.Assert(array2.Sum(), 6)
		t.AssertEQ(array1, array2)
	})
}

func TestLinkedList_Equals(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewLinkedList[int]()
		s2 := g.NewLinkedList[int]()
		s3 := g.NewLinkedList[int]()
		s4 := g.NewLinkedList[int]()
		s1.Add(1, 2, 3)
		s2.Add(1, 2, 3)
		s3.Add(1, 2, 3, 4)
		s4.Add(4, 5, 6)
		t.Assert(s1.Equals(s2), true)
		t.Assert(s1.Equals(s3), false)
		t.Assert(s1.Equals(s4), false)
		s5 := s1
		t.Assert(s1.Equals(s5), true)
	})
}

func TestLinkedList_PushFronts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2}
		l.PushFronts(a1)
		checkList(t, l, []int{2, 1})
		a1 = []int{3, 4, 5}
		l.PushFronts(a1)
		checkList(t, l, []int{5, 4, 3, 2, 1})
	})
}

func TestLinkedList_PushBacks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2}
		l.PushBacks(a1)
		checkList(t, l, []int{1, 2})
		a1 = []int{3, 4, 5}
		l.PushBacks(a1)
		checkList(t, l, []int{1, 2, 3, 4, 5})
	})
}

func TestLinkedList_PopBacks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		a2 := []int{7, 8, 9, 10}
		l.PushFronts(a1)
		i1 := l.PopBacks(2)
		t.Assert(i1, []interface{}{1, 2})

		l.PushBacks(a2) // 4.3,7,8,9,10
		i1 = l.PopBacks(3)
		t.Assert(i1, []interface{}{10, 9, 8})
	})
}

func TestLinkedList_PopFronts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.PopFronts(2)
		t.Assert(i1, []int{4, 3})
		t.Assert(l.Len(), 2)
	})
}

func TestLinkedList_PopBackAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.PopBackAll()
		t.Assert(i1, []int{1, 2, 3, 4})
		t.Assert(l.Len(), 0)
	})
}

func TestLinkedList_PopFrontAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.PopFrontAll()
		t.Assert(i1, []int{4, 3, 2, 1})
		t.Assert(l.Len(), 0)
	})
}

func TestLinkedList_FrontAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.FrontAll()
		t.Assert(i1, []int{4, 3, 2, 1})
		t.Assert(l.Len(), 4)
	})
}

func TestLinkedList_BackAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.BackAll()
		t.Assert(i1, []int{1, 2, 3, 4})
		t.Assert(l.Len(), 4)
	})
}

func TestLinkedList_FrontValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		l2 := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.FrontValue()
		t.Assert(gconv.Int(i1), 4)
		t.Assert(l.Len(), 4)

		i1 = l2.FrontValue()
		t.Assert(i1, 0)
	})
}

func TestLinkedList_BackValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		l2 := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		i1 := l.BackValue()
		t.Assert(gconv.Int(i1), 1)
		t.Assert(l.Len(), 4)

		i1 = l2.FrontValue()
		t.Assert(i1, 0)
	})
}

func TestLinkedList_Back(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		t.Assert(e1.Value, 1)
		t.Assert(l.Len(), 4)
	})
}

func TestLinkedList_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		t.Assert(l.Size(), 4)
		l.PopFront()
		t.Assert(l.Size(), 3)
	})
}

func TestLinkedList_Removes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		l.Remove(e1.Value)
		t.Assert(l.Len(), 3)

		e2 := l.Back()
		l.Remove(e2.Value)
		t.Assert(l.Len(), 2)
		checkList(t, l, []int{4, 3})
	})
}

func TestLinkedList_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedListFrom([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})

		back, _ := l.PopBack()
		t.Assert(back, 9)
		t.Assert(l.PopBacks(2), []int{8, 7})
		front, _ := l.PopFront()
		t.Assert(front, 1)
		t.Assert(l.PopFronts(2), []int{2, 3})
	})
}

func TestLinkedList_ForEachAsc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 5, 6, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		fun1 := func(e int) bool {
			return gconv.Int(e1.Value) > 2
		}
		checkList(t, l, []int{4, 3, 6, 5, 2, 1})
		l.ForEachAsc(fun1)
		checkList(t, l, []int{4, 3, 6, 5, 2, 1})
	})
}

func TestLinkedList_ForEachDesc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[int]()
		a1 := []int{1, 2, 3, 4}
		l.PushFronts(a1)
		e1 := l.Back()
		fun1 := func(e int) bool {
			return gconv.Int(e1.Value) > 6
		}
		l.ForEachDesc(fun1)
		t.Assert(l.Len(), 4)
		checkList(t, l, []int{4, 3, 2, 1})
	})
}

func TestLinkedList_ForEach(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedList[string]()
		a1 := []string{"a", "b", "c", "d", "e"}
		l.PushFronts(a1)
		e1 := l.Back()
		fun1 := func(e string) bool {
			return gconv.String(e1.Value) > "c"
		}
		checkList(t, l, []string{"e", "d", "c", "b", "a"})
		l.ForEach(fun1)
		checkList(t, l, []string{"e", "d", "c", "b", "a"})
	})
}

func TestLinkedList_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedListFrom[any]([]any{1, 2, "a", `"b"`, `\c`})
		t.Assert(l.Join(","), `1,2,a,"b",\c`)
		t.Assert(l.Join("."), `1.2.a."b".\c`)
	})
}

func TestLinkedList_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedListFrom[any]([]any{1, 2, "a", `"b"`, `\c`})
		t.Assert(l.String(), `[1,2,a,"b",\c]`)
	})
}

func TestLinkedList_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		a := []string{"a", "b", "c"}
		l := g.NewLinkedList[string]()
		l.PushBacks(a)
		b1, err1 := json.Marshal(l)
		b2, err2 := json.Marshal(a)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})
	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		a := []string{"a", "b", "c"}
		l := g.NewLinkedList[string]()
		b, err := json.Marshal(a)
		t.AssertNil(err)

		err = json.UnmarshalUseNumber(b, l)
		t.AssertNil(err)
		t.Assert(l.FrontAll(), a)
	})
	gtest.C(t, func(t *gtest.T) {
		var l g.LinkedList[string]
		a := []string{"a", "b", "c"}
		b, err := json.Marshal(a)
		t.AssertNil(err)

		err = json.UnmarshalUseNumber(b, &l)
		t.AssertNil(err)
		t.Assert(l.FrontAll(), a)
	})
}

func TestLinkedList_UnmarshalValue(t *testing.T) {
	type TList struct {
		Name string
		List *g.LinkedList[byte]
	}
	type TListInt struct {
		Name string
		List *g.LinkedList[int]
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var tlist *TList
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"list": []byte(`[1,2,3]`),
		}, &tlist)
		t.AssertNil(err)
		t.Assert(tlist.Name, "john")
		t.Assert(tlist.List.FrontAll(), []byte{1, 2, 3})
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var tlist *TListInt
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"list": []int{1, 2, 3},
		}, &tlist)
		t.AssertNil(err)
		t.Assert(tlist.Name, "john")
		t.Assert(tlist.List.FrontAll(), []int{1, 2, 3})
	})
}

func TestLinkedList_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := g.NewLinkedListFrom([]int{1, 2, 3, 4, 5})
		copyList := l.DeepCopy().(*g.LinkedList[int])
		copyList.PopBack()
		t.AssertNE(l.Size(), copyList.Size())
	})
}
