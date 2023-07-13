// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package g_test

import (
	"strings"
	"testing"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/internal/gtest"
	"github.com/wesleywu/gcontainer/utils/comparators"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

func TestTreeSet_NewTreeSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewTreeSet[int](comparators.ComparatorInt)
		s.Add(1, 1, 2)
		s.Add([]int{3, 4}...)
		t.Assert(s.Size(), 4)
		t.AssertIN(1, s.Slice())
		t.AssertIN(2, s.Slice())
		t.AssertIN(3, s.Slice())
		t.AssertIN(4, s.Slice())
		t.AssertNI(0, s.Slice())
		t.Assert(s.Contains(4), true)
		t.Assert(s.Contains(5), false)
		s.Remove(1)
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestTreeSet_NewTreeSetDefault(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewTreeSetDefault[int]()
		s.Add(1, 1, 2)
		s.Add([]int{3, 4}...)
		t.Assert(s.Size(), 4)
		t.AssertIN(1, s.Slice())
		t.AssertIN(2, s.Slice())
		t.AssertIN(3, s.Slice())
		t.AssertIN(4, s.Slice())
		t.AssertNI(0, s.Slice())
		t.Assert(s.Contains(4), true)
		t.Assert(s.Contains(5), false)
		s.Remove(1)
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestTreeSet_NewTreeSetFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "f", "c"}
		a2 := []string{"h", "j", "i", "k"}
		func2 := func(v1, v2 string) int {
			return strings.Compare(v2, v1)
		}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorString)
		array2 := g.NewTreeSetFrom(a2, func2)

		t.Assert(array1.Size(), 3)
		t.Assert(array1.String(), []string{"a", "c", "f"})

		t.Assert(array2.Size(), 4)
		t.Assert(array2.String(), []string{"k", "j", "i", "h"})
	})
}

func TestTreeSet_Add(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewTreeSetDefault[int](true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Add(1), false)
		t.Assert(s.Add(2), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.Add(2), false)
		t.Assert(s.Contains(2), true)
	})
}

func TestTreeSet_AddAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewTreeSetDefault[int](true)
		s.AddAll(g.NewArrayListFrom([]int{3, 1, 2}))
		t.Assert(s.Add(1), false)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Add(4), true)
		t.Assert(s.Contains(4), true)
	})
}

func TestTreeSet_Ceiling(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := g.NewTreeSetFrom(
			[]string{"a", "d", "c", "b"},
			comparators.ComparatorString,
		)
		s1, ok := array1.Ceiling("z")
		t.Assert(ok, false)
		s1, ok = array1.Ceiling("d")
		t.Assert(ok, true)
		t.Assert(s1, "d")
		s1, ok = array1.Ceiling("c1")
		t.Assert(ok, true)
		t.Assert(s1, "d")
		s1, ok = array1.Ceiling("_")
		t.Assert(ok, true)
		t.Assert(s1, "a")
	})
}

func TestTreeSet_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorString)
		t.Assert(array1.Size(), 6)
		array1.Clear()
		t.Assert(array1.Size(), 0)
	})
}

func TestTreeSet_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e", "f"}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorString)
		array2 := array1.Clone()
		t.Assert(array1, array2)
		array1.Remove("a")
		t.AssertNE(array1, array2)
	})
}

func TestTreeSet_Contains(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewTreeSetDefault[int](true)
		s.AddAll(g.NewArrayListFrom([]int{3, 1, 2}))
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(4), false)
	})
}

func TestTreeSet_ContainsAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewTreeSetDefault[int](true)
		s.AddAll(g.NewArrayListFrom([]int{3, 1, 2}))
		t.Assert(s.ContainsAll(g.NewArrayListFrom([]int{1, 2, 3})), true)
		t.Assert(s.ContainsAll(g.NewArrayListFrom([]int{1, 2})), true)
		t.Assert(s.ContainsAll(g.NewArrayListFrom([]int{2, 3})), true)
		t.Assert(s.ContainsAll(g.NewArrayListFrom([]int{2, 3, 4})), false)
		t.Assert(s.ContainsAll(g.NewArrayListFrom([]int{4})), false)
	})
}

func TestTreeSet_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := g.NewTreeSetFrom([]int{1, 2, 3, 4, 5}, comparators.ComparatorInt)
		copyArray := array.DeepCopy().(*g.TreeSet[int])
		array.Add(6)
		copyArray.Add(7)
		t.AssertEQ(copyArray.Contains(5), true)
		t.AssertEQ(copyArray.Contains(6), false)
		t.AssertEQ(copyArray.Contains(7), true)
	})
}

func TestTreeSet_Equals(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewTreeSetDefault[int]()
		s2 := g.NewTreeSetDefault[int]()
		s3 := g.NewTreeSetDefault[int]()
		s4 := g.NewTreeSetDefault[int]()
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

func TestTreeSet_First(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := g.NewTreeSetFrom(
			[]string{"a", "d", "c", "b"},
			comparators.ComparatorString,
		)
		i1, ok := array1.First()
		t.Assert(ok, true)
		t.Assert(i1, "a")
	})
	gtest.C(t, func(t *gtest.T) {
		array := g.NewTreeSetFrom[int]([]int{3, 1, 2}, comparators.ComparatorInt)
		v, ok := array.First()
		t.Assert(v, 1)
		t.Assert(ok, true)
	})
}

func TestTreeSet_Floor(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := g.NewTreeSetFrom(
			[]string{"a", "d", "c", "b"},
			comparators.ComparatorString,
		)
		s1, ok := array1.Floor("_")
		t.Assert(ok, false)
		s1, ok = array1.Floor("a")
		t.Assert(ok, true)
		t.Assert(s1, "a")
		s1, ok = array1.Floor("a1")
		t.Assert(ok, true)
		t.Assert(s1, "a")
		s1, ok = array1.Floor("z")
		t.Assert(ok, true)
		t.Assert(s1, "d")
	})
}

func TestTreeSet_ForEach(t *testing.T) {
	slice := []string{"a", "b", "d", "c"}
	sliceSorted := []string{"a", "b", "c", "d"}
	treeSet := g.NewTreeSetFrom(slice, comparators.ComparatorString)
	gtest.C(t, func(t *gtest.T) {
		index := 0
		treeSet.ForEach(func(v string) bool {
			t.Assert(v, sliceSorted[index])
			index++
			return true
		})
		t.Assert(index, 4)
	})
	gtest.C(t, func(t *gtest.T) {
		index := 0
		treeSet.ForEach(func(v string) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
}

func TestTreeSet_ForEachDescending(t *testing.T) {
	slice := []string{"a", "b", "d", "c"}
	sliceSortedDescending := []string{"d", "c", "b", "a"}
	treeSet := g.NewTreeSetFrom(slice, comparators.ComparatorString)
	gtest.C(t, func(t *gtest.T) {
		index := 0
		treeSet.ForEachDescending(func(v string) bool {
			t.Assert(v, sliceSortedDescending[index])
			index++
			return true
		})
		t.Assert(index, 4)
	})
	gtest.C(t, func(t *gtest.T) {
		index := 0
		treeSet.ForEachDescending(func(v string) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
}

func TestTreeSet_HeadSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e"}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorString)

		var i1 g.SortedSet[string]
		i1 = array1.HeadSet("c", false)
		t.Assert(i1.Slice(), []string{"a", "b"})

		i1 = array1.HeadSet("c", true)
		t.Assert(i1.Slice(), []string{"a", "b", "c"})

		i1 = array1.HeadSet("c1", true)
		t.Assert(i1.Slice(), []string{"a", "b", "c"})

		i1 = array1.HeadSet("c1", false)
		t.Assert(i1.Slice(), []string{"a", "b", "c"})
	})
}

func TestTreeSet_Higher(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := g.NewTreeSetFrom(
			[]string{"a", "d", "c", "b"},
			comparators.ComparatorString,
		)
		s1, ok := array1.Higher("d")
		t.Assert(ok, false)
		s1, ok = array1.Higher("c")
		t.Assert(ok, true)
		t.Assert(s1, "d")
		s1, ok = array1.Higher("c1")
		t.Assert(ok, true)
		t.Assert(s1, "d")
		s1, ok = array1.Higher("_")
		t.Assert(ok, true)
		t.Assert(s1, "a")
	})
}

func TestTreeSet_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := g.NewTreeSetFrom([]string{}, comparators.ComparatorString)
		t.Assert(array.IsEmpty(), true)
	})
}

func TestTreeSet_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c"}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorString)
		t.Assert(array1.Join(","), `a,c,d`)
		t.Assert(array1.Join("."), `a.c.d`)
	})

	gtest.C(t, func(t *gtest.T) {
		a1 := []int{0, 1, 2, 3}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorInt)
		t.Assert(array1.Join("."), `0.1.2.3`)
	})

	gtest.C(t, func(t *gtest.T) {
		a1 := []string{}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorString)
		t.Assert(array1.Join("."), "")
	})
}

func TestTreeSet_Last(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := g.NewTreeSetFrom(
			[]string{"a", "d", "c", "b"},
			comparators.ComparatorString,
		)
		i1, ok := array1.Last()
		t.Assert(ok, true)
		t.Assert(i1, "d")
	})
	gtest.C(t, func(t *gtest.T) {
		array := g.NewTreeSetFrom[int]([]int{3, 1, 2}, comparators.ComparatorInt)
		v, ok := array.Last()
		t.Assert(v, 3)
		t.Assert(ok, true)
	})
}

func TestTreeSet_Lower(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := g.NewTreeSetFrom(
			[]string{"a", "d", "c", "b"},
			comparators.ComparatorString,
		)
		s1, ok := array1.Lower("a")
		t.Assert(ok, false)
		s1, ok = array1.Lower("b")
		t.Assert(ok, true)
		t.Assert(s1, "a")
		s1, ok = array1.Lower("b1")
		t.Assert(ok, true)
		t.Assert(s1, "b")
		s1, ok = array1.Lower("z")
		t.Assert(ok, true)
		t.Assert(s1, "d")
	})
}

func TestTreeSet_PollFirst(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := g.NewTreeSetFrom(
			[]string{"a", "d", "c", "b"},
			comparators.ComparatorString,
		)
		i1, ok := array1.PollFirst()
		t.Assert(ok, true)
		t.Assert(gconv.String(i1), "a")
		t.Assert(array1.Size(), 3)
		t.Assert(array1.Slice(), []string{"b", "c", "d"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := g.NewTreeSetFrom[int]([]int{1, 2, 3}, comparators.ComparatorInt)
		v, ok := array.PollFirst()
		t.Assert(v, 1)
		t.Assert(ok, true)
		t.Assert(array.Size(), 2)
		v, ok = array.PollFirst()
		t.Assert(v, 2)
		t.Assert(ok, true)
		t.Assert(array.Size(), 1)
		v, ok = array.PollFirst()
		t.Assert(v, 3)
		t.Assert(ok, true)
		t.Assert(array.Size(), 0)
	})
}

func TestTreeSet_PollHeadSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e"}
		setOrigin := g.NewTreeSetFrom(a1, comparators.ComparatorString)

		s1 := setOrigin.Clone().(*g.TreeSet[string])
		r1 := s1.PollHeadSet("c", false)
		t.Assert(r1.Slice(), []string{"a", "b"})
		t.Assert(s1.Slice(), []string{"c", "d", "e"})

		s2 := setOrigin.Clone().(*g.TreeSet[string])
		r2 := s2.PollHeadSet("c", true)
		t.Assert(r2.Slice(), []string{"a", "b", "c"})
		t.Assert(s2.Slice(), []string{"d", "e"})

		s3 := setOrigin.Clone().(*g.TreeSet[string])
		r3 := s3.PollHeadSet("c1", true)
		t.Assert(r3.Slice(), []string{"a", "b", "c"})
		t.Assert(s3.Slice(), []string{"d", "e"})

		s4 := setOrigin.Clone().(*g.TreeSet[string])
		r4 := s4.PollHeadSet("c1", false)
		t.Assert(r4.Slice(), []string{"a", "b", "c"})
		t.Assert(s4.Slice(), []string{"d", "e"})

		s5 := setOrigin.Clone().(*g.TreeSet[string])
		r5 := s5.PollHeadSet("z", true)
		t.Assert(r5.Slice(), []string{"a", "b", "c", "d", "e"})
		t.Assert(s5.Slice(), []string{})

		s6 := setOrigin.Clone().(*g.TreeSet[string])
		r6 := s6.PollHeadSet("_", true)
		t.Assert(r6.Slice(), []string{})
		t.Assert(s6.Slice(), []string{"a", "b", "c", "d", "e"})
	})
}

func TestTreeSet_PollLast(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array1 := g.NewTreeSetFrom(
			[]string{"a", "d", "c", "b"},
			comparators.ComparatorString,
		)
		i1, ok := array1.PollLast()
		t.Assert(ok, true)
		t.Assert(gconv.String(i1), "d")
		t.Assert(array1.Size(), 3)
		t.Assert(array1.Slice(), []string{"a", "b", "c"})
	})
	gtest.C(t, func(t *gtest.T) {
		array := g.NewTreeSetFrom[int]([]int{1, 2, 3}, comparators.ComparatorInt)
		v, ok := array.PollLast()
		t.Assert(v, 3)
		t.Assert(ok, true)
		t.Assert(array.Size(), 2)

		v, ok = array.PollLast()
		t.Assert(v, 2)
		t.Assert(ok, true)
		t.Assert(array.Size(), 1)

		v, ok = array.PollLast()
		t.Assert(v, 1)
		t.Assert(ok, true)
		t.Assert(array.Size(), 0)
	})
}

func TestTreeSet_PollTailSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e"}
		setOrigin := g.NewTreeSetFrom(a1, comparators.ComparatorString)

		s1 := setOrigin.Clone().(*g.TreeSet[string])
		r1 := s1.PollTailSet("c", false)
		t.Assert(r1.Slice(), []string{"d", "e"})
		t.Assert(s1.Slice(), []string{"a", "b", "c"})

		s2 := setOrigin.Clone().(*g.TreeSet[string])
		r2 := s2.PollTailSet("c", true)
		t.Assert(r2.Slice(), []string{"c", "d", "e"})
		t.Assert(s2.Slice(), []string{"a", "b"})

		s3 := setOrigin.Clone().(*g.TreeSet[string])
		r3 := s3.PollTailSet("c1", true)
		t.Assert(r3.Slice(), []string{"d", "e"})
		t.Assert(s3.Slice(), []string{"a", "b", "c"})

		s4 := setOrigin.Clone().(*g.TreeSet[string])
		r4 := s4.PollTailSet("c1", false)
		t.Assert(r4.Slice(), []string{"d", "e"})
		t.Assert(s4.Slice(), []string{"a", "b", "c"})

		s5 := setOrigin.Clone().(*g.TreeSet[string])
		r5 := s5.PollTailSet("z", true)
		t.Assert(r5.Slice(), []string{})
		t.Assert(s5.Slice(), []string{"a", "b", "c", "d", "e"})

		s6 := setOrigin.Clone().(*g.TreeSet[string])
		r6 := s6.PollTailSet("_", true)
		t.Assert(r6.Slice(), []string{"a", "b", "c", "d", "e"})
		t.Assert(s6.Slice(), []string{})
	})
}

func TestTreeSet_Remove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorString)
		changed := array1.Remove("b")
		t.Assert(changed, true)
		t.Assert(array1.Size(), 3)
		t.Assert(array1.Contains("b"), false)

		changed = array1.Remove("e")
		t.Assert(changed, false)

		changed = array1.Remove("a", "d")
		t.Assert(changed, true)
		t.Assert(array1.Size(), 1)
		t.Assert(array1.Contains("a"), false)
		t.Assert(array1.Contains("d"), false)

		changed = array1.Remove("a", "d")
		t.Assert(changed, false)
		t.Assert(array1.Size(), 1)
		t.Assert(array1.Contains("a"), false)
		t.Assert(array1.Contains("d"), false)

		changed = array1.Remove("a", "d", "c", "b")
		t.Assert(changed, true)
		t.Assert(array1.Size(), 0)
		t.Assert(array1.Contains("c"), false)
	})
}

func TestTreeSet_RemoveAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b"}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorString)
		changed := array1.RemoveAll(g.NewArrayListFrom([]string{"b"}))
		t.Assert(changed, true)
		t.Assert(array1.Size(), 3)
		t.Assert(array1.Contains("b"), false)

		changed = array1.RemoveAll(g.NewArrayListFrom([]string{"e"}))
		t.Assert(changed, false)

		changed = array1.RemoveAll(g.NewArrayListFrom([]string{"a", "d"}))
		t.Assert(changed, true)
		t.Assert(array1.Size(), 1)
		t.Assert(array1.Contains("a"), false)
		t.Assert(array1.Contains("d"), false)

		changed = array1.RemoveAll(g.NewArrayListFrom([]string{"a", "d"}))
		t.Assert(changed, false)
		t.Assert(array1.Size(), 1)
		t.Assert(array1.Contains("a"), false)
		t.Assert(array1.Contains("d"), false)

		changed = array1.RemoveAll(g.NewArrayListFrom([]string{"a", "d", "c", "b"}))
		t.Assert(changed, true)
		t.Assert(array1.Size(), 0)
		t.Assert(array1.Contains("c"), false)
	})
}

func TestTreeSet_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{0, 1, 2, 3}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorInt)
		t.Assert(array1.String(), `[0,1,2,3]`)

		array1 = nil
		t.Assert(array1.String(), "")

		array1 = g.NewTreeSetDefault[int]()
		t.Assert(array1.String(), "[]")

		array2 := g.NewTreeSetDefault[string]()
		t.Assert(array2.String(), "[]")
	})
}

func TestTreeSet_TailSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []string{"a", "d", "c", "b", "e"}
		array1 := g.NewTreeSetFrom(a1, comparators.ComparatorString)

		var i1 g.SortedSet[string]
		i1 = array1.TailSet("c", false)
		t.Assert(i1.Slice(), []string{"d", "e"})

		i1 = array1.TailSet("c", true)
		t.Assert(i1.Slice(), []string{"c", "d", "e"})

		i1 = array1.TailSet("c1", true)
		t.Assert(i1.Slice(), []string{"d", "e"})

		i1 = array1.TailSet("c1", false)
		t.Assert(i1.Slice(), []string{"d", "e"})
	})
}

func TestTreeSet_SubSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := make([]string, 10)
		for i := 0; i < 10; i++ {
			m[i] = "key" + gconv.String(i)
		}
		treeSet := g.NewTreeSetFrom(m, comparators.ComparatorString)
		// both key exists in map
		t.Assert(treeSet.SubSet("key5", true, "key7", true).Slice(), []string{"key5", "key6", "key7"})
		t.Assert(treeSet.SubSet("key5", false, "key7", true).Slice(), []string{"key6", "key7"})
		t.Assert(treeSet.SubSet("key5", true, "key7", false).Slice(), []string{"key5", "key6"})
		t.Assert(treeSet.SubSet("key5", false, "key7", false).Slice(), []string{"key6"})
		// only fromKey exists in map
		t.Assert(treeSet.SubSet("key5.1", true, "key7", true).Slice(), []string{"key6", "key7"})
		t.Assert(treeSet.SubSet("key5.1", false, "key7", true).Slice(), []string{"key6", "key7"})
		t.Assert(treeSet.SubSet("key5.1", true, "key7", false).Slice(), []string{"key6"})
		t.Assert(treeSet.SubSet("key5.1", false, "key7", false).Slice(), []string{"key6"})
		// both key do not exist in map
		t.Assert(treeSet.SubSet("key5.1", true, "key7.1", true).Slice(), []string{"key6", "key7"})
		t.Assert(treeSet.SubSet("key5.1", false, "key7.1", true).Slice(), []string{"key6", "key7"})
		t.Assert(treeSet.SubSet("key5.1", true, "key7.1", false).Slice(), []string{"key6", "key7"})
		t.Assert(treeSet.SubSet("key5.1", false, "key7.1", false).Slice(), []string{"key6", "key7"})
		// fromKey out of upper bound
		t.Assert(treeSet.SubSet("zz", false, "key7.1", false).Slice(), []string{})
		// fromKey out of lower bound
		t.Assert(treeSet.SubSet("aa", false, "key0.1", false).Slice(), []string{"key0"})
		// both key out of lower bound
		t.Assert(treeSet.SubSet("aa", false, "bb", false).Slice(), []string{})
		// both key out of lower bound
		t.Assert(treeSet.SubSet("bb", false, "aa", false).Slice(), []string{})
		// both key out of upper bound
		t.Assert(treeSet.SubSet("yy", false, "zz", false).Slice(), []string{})
		// both key out of upper bound
		t.Assert(treeSet.SubSet("zz", false, "yy", false).Slice(), []string{})
		// toKey out of upper bound
		t.Assert(treeSet.SubSet("key9", true, "zz", false).Slice(), []string{"key9"})
		// fromKey out of lower bound and toKey out of upper bound
		t.Assert(treeSet.SubSet("aa", true, "zz", false).Slice(),
			[]string{"key0", "key1", "key2", "key3", "key4", "key5", "key6", "key7", "key8", "key9"})
	})
}

//func TestTreeSet_PopRand(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []string{"a", "d", "c", "b"}
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		array1 := gset.NewTreeSetFrom(a1, func1)
//		i1, ok := array1.PopRand()
//		t.Assert(ok, true)
//		t.AssertIN(i1, []string{"a", "d", "c", "b"})
//		t.Assert(array1.Size(), 3)
//
//	})
//}
//
//func TestTreeSet_PopRands(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []string{"a", "d", "c", "b"}
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		array1 := gset.NewTreeSetFrom(a1, func1)
//		i1 := array1.PopRands(2)
//		t.Assert(len(i1), 2)
//		t.AssertIN(i1, []string{"a", "d", "c", "b"})
//		t.Assert(array1.Size(), 2)
//
//		i2 := array1.PopRands(3)
//		t.Assert(len(i1), 2)
//		t.AssertIN(i2, []string{"a", "d", "c", "b"})
//		t.Assert(array1.Size(), 0)
//
//	})
//}
//
//func TestTreeSet_Empty(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		array := gset.NewTreeSet[int](comparators.ComparatorInt)
//		v, ok := array.PopLeft()
//		t.Assert(v, 0)
//		t.Assert(ok, false)
//		t.Assert(array.PopLefts(10), nil)
//
//		v, ok = array.PopRight()
//		t.Assert(v, 0)
//		t.Assert(ok, false)
//		t.Assert(array.PopRights(10), nil)
//
//		v, ok = array.PopRand()
//		t.Assert(v, 0)
//		t.Assert(ok, false)
//		t.Assert(array.PopRands(10), nil)
//	})
//}
//
//func TestTreeSet_PopLefts(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []string{"a", "d", "c", "b", "e", "f"}
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		array1 := gset.NewTreeSetFrom(a1, func1)
//		i1 := array1.PopLefts(2)
//		t.Assert(len(i1), 2)
//		t.AssertIN(i1, []string{"a", "d", "c", "b", "e", "f"})
//		t.Assert(array1.Size(), 4)
//
//		i2 := array1.PopLefts(5)
//		t.Assert(len(i2), 4)
//		t.AssertIN(i1, []string{"a", "d", "c", "b", "e", "f"})
//		t.Assert(array1.Size(), 0)
//	})
//}
//
//func TestTreeSet_PopRights(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []string{"a", "d", "c", "b", "e", "f"}
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		array1 := gset.NewTreeSetFrom(a1, func1)
//		i1 := array1.PopRights(2)
//		t.Assert(len(i1), 2)
//		t.Assert(i1, []string{"e", "f"})
//		t.Assert(array1.Size(), 4)
//
//		i2 := array1.PopRights(10)
//		t.Assert(len(i2), 4)
//		t.Assert(array1.Size(), 0)
//	})
//}
//

//
//func TestTreeSet_Sum(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []string{"a", "d", "c", "b", "e", "f"}
//		a2 := []string{"1", "2", "3", "b", "e", "f"}
//		a3 := []string{"4", "5", "6"}
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		array1 := gset.NewTreeSetFrom(a1, func1)
//		array2 := gset.NewTreeSetFrom(a2, func1)
//		array3 := gset.NewTreeSetFrom(a3, func1)
//		t.Assert(array1.Sum(), 0)
//		t.Assert(array2.Sum(), 6)
//		t.Assert(array3.Sum(), 15)
//
//	})
//}

//
//func TestTreeSet_Chunk(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []string{"a", "d", "c", "b", "e"}
//
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		array1 := gset.NewTreeSetFrom(a1, func1)
//		i1 := array1.Chunk(2)
//		t.Assert(len(i1), 3)
//		t.Assert(i1[0], []string{"a", "b"})
//		t.Assert(i1[2], []string{"e"})
//
//		i1 = array1.Chunk(0)
//		t.Assert(len(i1), 0)
//	})
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []int{1, 2, 3, 4, 5}
//		array1 := gset.NewTreeSetFrom[int](a1, comparators.ComparatorInt)
//		chunks := array1.Chunk(3)
//		t.Assert(len(chunks), 2)
//		t.Assert(chunks[0], []int{1, 2, 3})
//		t.Assert(chunks[1], []int{4, 5})
//		t.Assert(array1.Chunk(0), nil)
//	})
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []int{1, 2, 3, 4, 5, 6}
//		array1 := gset.NewTreeSetFrom(a1, comparators.ComparatorInt)
//		chunks := array1.Chunk(2)
//		t.Assert(len(chunks), 3)
//		t.Assert(chunks[0], []int{1, 2})
//		t.Assert(chunks[1], []int{3, 4})
//		t.Assert(chunks[2], []int{5, 6})
//		t.Assert(array1.Chunk(0), nil)
//	})
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []int{1, 2, 3, 4, 5, 6}
//		array1 := gset.NewTreeSetFrom(a1, comparators.ComparatorInt)
//		chunks := array1.Chunk(3)
//		t.Assert(len(chunks), 2)
//		t.Assert(chunks[0], []int{1, 2, 3})
//		t.Assert(chunks[1], []int{4, 5, 6})
//		t.Assert(array1.Chunk(0), nil)
//	})
//}
//
//func TestTreeSet_SubSlice(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []string{"a", "d", "c", "b", "e"}
//
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		array1 := gset.NewTreeSetFrom(a1, func1)
//		array2 := gset.NewTreeSetFrom(a1, func1, true)
//		i1 := array1.SubSlice(2, 3)
//		t.Assert(len(i1), 3)
//		t.Assert(i1, []string{"c", "d", "e"})
//
//		i1 = array1.SubSlice(2, 6)
//		t.Assert(len(i1), 3)
//		t.Assert(i1, []string{"c", "d", "e"})
//
//		i1 = array1.SubSlice(7, 2)
//		t.Assert(len(i1), 0)
//
//		s1 := array1.SubSlice(1, -2)
//		t.Assert(s1, nil)
//
//		s1 = array1.SubSlice(-9, 2)
//		t.Assert(s1, nil)
//		t.Assert(array2.SubSlice(1, 3), []string{"b", "c", "d"})
//
//	})
//}
//
//func TestTreeSet_Rand(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []string{"a", "d", "c"}
//
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		array1 := gset.NewTreeSetFrom(a1, func1)
//		i1, ok := array1.Rand()
//		t.Assert(ok, true)
//		t.AssertIN(i1, []string{"a", "d", "c"})
//		t.Assert(array1.Size(), 3)
//
//		array2 := gset.NewTreeSetFrom([]string{}, func1)
//		v, ok := array2.Rand()
//		t.Assert(ok, false)
//		t.Assert(v, nil)
//	})
//}
//
//func TestTreeSet_Rands(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []string{"a", "d", "c"}
//
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		array1 := gset.NewTreeSetFrom(a1, func1)
//		i1 := array1.Rands(2)
//		t.AssertIN(i1, []string{"a", "d", "c"})
//		t.Assert(len(i1), 2)
//		t.Assert(array1.Size(), 3)
//
//		i1 = array1.Rands(4)
//		t.Assert(len(i1), 4)
//
//		array2 := gset.NewTreeSetFrom([]string{}, func1)
//		v := array2.Rands(1)
//		t.Assert(v, nil)
//	})
//}
//

//func TestTreeSet_CountValues(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []string{"a", "d", "c", "c"}
//
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		array1 := gset.NewTreeSetFrom(a1, func1)
//		m1 := array1.CountValues()
//		t.Assert(len(m1), 3)
//		t.Assert(m1["c"], 2)
//		t.Assert(m1["a"], 1)
//
//	})
//}
//
//func TestTreeSet_SetUnique(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []int{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
//		array1 := gset.NewTreeSetFrom(a1, comparators.ComparatorInt)
//		array1.SetUnique(true)
//		t.Assert(array1.Size(), 5)
//		t.Assert(array1, []int{1, 2, 3, 4, 5})
//	})
//}
//
//func TestTreeSet_Unique(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		a1 := []int{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
//		array1 := gset.NewTreeSetFrom(a1, comparators.ComparatorInt)
//		array1.Unique()
//		t.Assert(array1.Size(), 5)
//		t.Assert(array1, []int{1, 2, 3, 4, 5})
//
//		array2 := gset.NewTreeSetFrom([]int{}, comparators.ComparatorInt)
//		array2.Unique()
//		t.Assert(array2.Size(), 0)
//		t.Assert(array2, []int{})
//	})
//}
//
//func TestTreeSet_LockFunc(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		s1 := []string{"a", "b", "c", "d"}
//		a1 := gset.NewTreeSetFrom(s1, func1, true)
//
//		ch1 := make(chan int64, 3)
//		ch2 := make(chan int64, 3)
//		// go1
//		go a1.LockFunc(func(n1 []string) { // 读写锁
//			time.Sleep(2 * time.Second) // 暂停2秒
//			n1[2] = "g"
//			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
//		})
//
//		// go2
//		go func() {
//			time.Sleep(100 * time.Millisecond) // 故意暂停0.01秒,等go1执行锁后，再开始执行.
//			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
//			a1.Size()
//			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
//		}()
//
//		t1 := <-ch1
//		t2 := <-ch1
//		<-ch2 // 等待go1完成
//
//		// 防止ci抖动,以豪秒为单位
//		t.AssertGT(t2-t1, 20) // go1加的读写互斥锁，所go2读的时候被阻塞。
//		t.Assert(a1.Contains("g"), true)
//	})
//}
//
//func TestTreeSet_RLockFunc(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		func1 := func(v1, v2 string) int {
//			return strings.Compare(gconv.String(v1), gconv.String(v2))
//		}
//		s1 := []string{"a", "b", "c", "d"}
//		a1 := gset.NewTreeSetFrom(s1, func1, true)
//
//		ch1 := make(chan int64, 3)
//		ch2 := make(chan int64, 3)
//		// go1
//		go a1.RLockFunc(func(n1 []string) { // 读写锁
//			time.Sleep(2 * time.Second) // 暂停2秒
//			n1[2] = "g"
//			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
//		})
//
//		// go2
//		go func() {
//			time.Sleep(100 * time.Millisecond) // 故意暂停0.01秒,等go1执行锁后，再开始执行.
//			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
//			a1.Size()
//			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
//		}()
//
//		t1 := <-ch1
//		t2 := <-ch1
//		<-ch2 // 等待go1完成
//
//		// 防止ci抖动,以豪秒为单位
//		t.AssertLT(t2-t1, 20) // go1加的读锁，所go2读的时候不会被阻塞。
//		t.Assert(a1.Contains("g"), true)
//	})
//}
//
//func TestTreeSet_Merge(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		s1 := []string{"a", "b", "c", "d"}
//		s2 := []string{"e", "f"}
//		i2 := garray.NewArrayFrom([]string{"3"})
//		s3 := garray.NewArrayFrom([]string{"g", "h"})
//		s4 := gset.NewTreeSetFrom([]string{"4", "5"}, comparators.ComparatorString)
//		s5 := gset.NewTreeSetFrom(s2, comparators.ComparatorString)
//		s6 := gset.NewTreeSetFrom([]string{"1", "2", "3"}, comparators.ComparatorString)
//
//		a1 := gset.NewTreeSetFrom(s1, comparators.ComparatorString)
//
//		t.Assert(a1.MergeSlice(s2).Size(), 6)
//		t.Assert(a1.Merge(s3).Size(), 8)
//		t.Assert(a1.Merge(i2).Size(), 9)
//		t.Assert(a1.Merge(s3).Size(), 11)
//		t.Assert(a1.Merge(s4).Size(), 13)
//		t.Assert(a1.Merge(s5).Size(), 15)
//		t.Assert(a1.Merge(s6).Size(), 18)
//	})
//}
//
//func TestTreeSet_Json(t *testing.T) {
//	// array pointer
//	gtest.C(t, func(t *gtest.T) {
//		s1 := []string{"a", "b", "d", "c"}
//		s2 := []string{"a", "b", "c", "d"}
//		a1 := gset.NewTreeSetFrom(s1, comparators.ComparatorString)
//		b1, err1 := json.Marshal(a1)
//		b2, err2 := json.Marshal(s1)
//		t.Assert(b1, b2)
//		t.Assert(err1, err2)
//
//		a2 := gset.NewTreeSet(comparators.ComparatorString)
//		err1 = json.UnmarshalUseNumber(b2, &a2)
//		t.AssertNil(err1)
//		t.Assert(a2.Slice(), s2)
//
//		var a3 gset.TreeSet[string]
//		err := json.UnmarshalUseNumber(b2, &a3)
//		t.AssertNil(err)
//		t.Assert(a3.Slice(), s1)
//		t.Assert(a3.Interfaces(), s1)
//	})
//	// array value
//	gtest.C(t, func(t *gtest.T) {
//		s1 := []string{"a", "b", "d", "c"}
//		s2 := []string{"a", "b", "c", "d"}
//		a1 := *gset.NewTreeSetFrom(s1, comparators.ComparatorString)
//		b1, err1 := json.Marshal(a1)
//		b2, err2 := json.Marshal(s1)
//		t.Assert(b1, b2)
//		t.Assert(err1, err2)
//
//		a2 := gset.NewTreeSet(comparators.ComparatorString)
//		err1 = json.UnmarshalUseNumber(b2, &a2)
//		t.AssertNil(err1)
//		t.Assert(a2.Slice(), s2)
//
//		var a3 gset.TreeSet[string]
//		err := json.UnmarshalUseNumber(b2, &a3)
//		t.AssertNil(err)
//		t.Assert(a3.Slice(), s1)
//		t.Assert(a3.Interfaces(), s1)
//	})
//	// array pointer
//	gtest.C(t, func(t *gtest.T) {
//		type User struct {
//			Name   string
//			Scores *gset.TreeSet[int]
//		}
//		data := map[string]interface{}{
//			"Name":   "john",
//			"Scores": []int{99, 100, 98},
//		}
//		b, err := json.Marshal(data)
//		t.AssertNil(err)
//
//		user := new(User)
//		err = json.UnmarshalUseNumber(b, user)
//		t.AssertNil(err)
//		t.Assert(user.Name, data["Name"])
//		t.AssertNE(user.Scores, nil)
//		t.Assert(user.Scores.Size(), 3)
//
//		v, ok := user.Scores.PopLeft()
//		t.AssertIN(v, data["Scores"])
//		t.Assert(ok, true)
//
//		v, ok = user.Scores.PopLeft()
//		t.AssertIN(v, data["Scores"])
//		t.Assert(ok, true)
//
//		v, ok = user.Scores.PopLeft()
//		t.AssertIN(v, data["Scores"])
//		t.Assert(ok, true)
//
//		v, ok = user.Scores.PopLeft()
//		t.Assert(v, 0)
//		t.Assert(ok, false)
//	})
//	// array value
//	gtest.C(t, func(t *gtest.T) {
//		type User struct {
//			Name   string
//			Scores gset.TreeSet[int]
//		}
//		data := map[string]interface{}{
//			"Name":   "john",
//			"Scores": []int{99, 100, 98},
//		}
//		b, err := json.Marshal(data)
//		t.AssertNil(err)
//
//		user := new(User)
//		err = json.UnmarshalUseNumber(b, user)
//		t.AssertNil(err)
//		t.Assert(user.Name, data["Name"])
//		t.AssertNE(user.Scores, nil)
//		t.Assert(user.Scores.Size(), 3)
//
//		v, ok := user.Scores.PopLeft()
//		t.AssertIN(v, data["Scores"])
//		t.Assert(ok, true)
//
//		v, ok = user.Scores.PopLeft()
//		t.AssertIN(v, data["Scores"])
//		t.Assert(ok, true)
//
//		v, ok = user.Scores.PopLeft()
//		t.AssertIN(v, data["Scores"])
//		t.Assert(ok, true)
//
//		v, ok = user.Scores.PopLeft()
//		t.Assert(v, 0)
//		t.Assert(ok, false)
//	})
//}
//

//	func TestTreeSet_RemoveValue(t *testing.T) {
//		slice := []string{"a", "b", "d", "c"}
//		array := gset.NewTreeSetFrom(slice, comparators.ComparatorString)
//		gtest.C(t, func(t *gtest.T) {
//			t.Assert(array.RemoveValue("e"), false)
//			t.Assert(array.RemoveValue("b"), true)
//			t.Assert(array.RemoveValue("a"), true)
//			t.Assert(array.RemoveValue("c"), true)
//			t.Assert(array.RemoveValue("f"), false)
//		})
//	}
//
//	func TestTreeSet_RemoveValues(t *testing.T) {
//		slice := []string{"a", "b", "d", "c"}
//		array := gset.NewTreeSetFrom(slice, comparators.ComparatorString)
//		gtest.C(t, func(t *gtest.T) {
//			array.RemoveValues("a", "b", "c")
//			t.Assert(array.Slice(), []string{"d"})
//		})
//	}
func TestTreeSet_UnmarshalValue(t *testing.T) {
	type V struct {
		Name  string
		Array *g.TreeSet[byte]
	}
	type VInt struct {
		Name  string
		Array *g.TreeSet[int]
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name":  "john",
			"array": []byte(`[2,3,1]`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), []byte{1, 2, 3})
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *VInt
		err := gconv.Struct(map[string]interface{}{
			"name":  "john",
			"array": []int{2, 3, 1},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), []int{1, 2, 3})
	})
}

//func comparatorExampleElement(a, b *garray.exampleElement) int {
//	if a == nil && b == nil {
//		return 0
//	}
//	if a == nil && b != nil {
//		return -1
//	}
//	if a != nil && b == nil {
//		return 1
//	}
//	return a.code - b.code
//}
//
//func TestTreeSet_Filter(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		values := []*garray.exampleElement{
//			{code: 2},
//			{code: 0},
//			{code: 1},
//		}
//		array := gset.NewTreeSetFromCopy[*garray.exampleElement](values, comparatorExampleElement)
//		t.Assert(array.Filter(func(index int, value *garray.exampleElement) bool {
//			return empty.IsNil(value)
//		}).Slice(), []*garray.exampleElement{
//			{code: 0},
//			{code: 1},
//			{code: 2},
//		})
//	})
//	gtest.C(t, func(t *gtest.T) {
//		values := []*garray.exampleElement{
//			nil,
//			{code: 2},
//			{code: 0},
//			{code: 1},
//			nil,
//		}
//		array := gset.NewTreeSetFromCopy[*garray.exampleElement](values, comparatorExampleElement)
//		t.Assert(array.Filter(func(index int, value *garray.exampleElement) bool {
//			return empty.IsNil(value)
//		}).Slice(), []*garray.exampleElement{
//			{code: 0},
//			{code: 1},
//			{code: 2},
//		})
//	})
//	gtest.C(t, func(t *gtest.T) {
//		values := []*garray.exampleElement{
//			{code: 2},
//			{},
//			{code: 0},
//			{code: 1},
//		}
//		array := garray.NewArrayFromCopy(values)
//		t.Assert(array.Filter(func(index int, value *garray.exampleElement) bool {
//			return empty.IsEmpty(value)
//		}).Slice(), []*garray.exampleElement{
//			{code: 1},
//			{code: 2},
//		})
//	})
//	gtest.C(t, func(t *gtest.T) {
//		values := []*garray.exampleElement{
//			{code: 2},
//			{code: 3},
//			{code: 1},
//		}
//		array := gset.NewTreeSetFromCopy[*garray.exampleElement](values, comparatorExampleElement)
//		t.Assert(array.Filter(func(index int, value *garray.exampleElement) bool {
//			return empty.IsEmpty(value)
//		}).Slice(), []*garray.exampleElement{
//			{code: 1},
//			{code: 2},
//			{code: 3},
//		})
//	})
//}
//
//func TestTreeSet_FilterNil(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		values := []*garray.exampleElement{
//			{code: 1},
//			{code: 0},
//			{code: 2},
//		}
//		array := gset.NewTreeSetFromCopy[*garray.exampleElement](values, comparatorExampleElement)
//		t.Assert(array.FilterNil().Slice(), []*garray.exampleElement{
//			{code: 0},
//			{code: 1},
//			{code: 2},
//		})
//	})
//	gtest.C(t, func(t *gtest.T) {
//		values := []*garray.exampleElement{
//			nil,
//			{code: 1},
//			{code: 0},
//			{code: 2},
//			nil,
//		}
//		array := gset.NewTreeSetFromCopy[*garray.exampleElement](values, comparatorExampleElement)
//		t.Assert(array.FilterNil().Slice(), []*garray.exampleElement{
//			{code: 0},
//			{code: 1},
//			{code: 2},
//		})
//	})
//}
//
//func TestTreeSet_FilterEmpty(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		values := []*garray.exampleElement{
//			{code: 2},
//			{},
//			{code: 0},
//			{code: 1},
//		}
//		array := garray.NewArrayFromCopy(values)
//		t.Assert(array.FilterEmpty().Slice(), []*garray.exampleElement{
//			{code: 1},
//			{code: 2},
//		})
//	})
//	gtest.C(t, func(t *gtest.T) {
//		values := []*garray.exampleElement{
//			{code: 2},
//			{code: 3},
//			{code: 1},
//		}
//		array := garray.NewArrayFromCopy(values)
//		t.Assert(array.FilterEmpty().Slice(), []*garray.exampleElement{
//			{code: 1},
//			{code: 2},
//			{code: 3},
//		})
//	})
//}
//
//func TestTreeSet_Walk(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		array := gset.NewTreeSetFrom([]string{"1", "2"}, comparators.ComparatorString)
//		t.Assert(array.Walk(func(value string) string {
//			return "key-" + gconv.String(value)
//		}), []string{"key-1", "key-2"})
//	})
//}
//
