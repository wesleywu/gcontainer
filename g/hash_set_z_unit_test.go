// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package g_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/internal/gtest"
	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

func TestHashSet_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s g.HashSet[int]
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

func TestHashSet_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[int]()
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

func TestHashSet_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[int]()
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

func TestHashSet_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []int{0, 1, 2, 3}
		array1 := g.NewHashSetFrom(a1)
		array2 := array1.Clone().(*g.HashSet[int])

		t.Assert(array2.Size(), 4)
		t.Assert(array2.Sum(), 6)
	})
}

func TestHashSet_Equals(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[int]()
		s2 := g.NewHashSet[int]()
		s3 := g.NewHashSet[int]()
		s4 := g.NewHashSet[int]()
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

func TestHashSet_ForEach(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[int]()
		s.Add(1, 2, 3)
		t.Assert(s.Size(), 3)

		a1 := g.NewHashSet[int](true)
		a2 := g.NewHashSet[int](true)
		s.ForEach(func(v int) bool {
			a1.Add(v)
			return false
		})
		s.ForEach(func(v int) bool {
			a2.Add(v)
			return true
		})
		t.Assert(a1.Size(), 1)
		t.Assert(a2.Size(), 3)
	})
}

func TestHashSet_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[int]()
		s.Add(1, 2, 3)
		t.Assert(s.Size(), 3)
		s.LockFunc(func(m map[int]struct{}) {
			delete(m, 1)
		})
		t.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[int]struct{}) {
			t.Assert(m, map[int]struct{}{
				3: {},
				2: {},
			})
		})
	})
}

func TestHashSet_IsSubsetOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[int]()
		s2 := g.NewHashSet[int]()
		s3 := g.NewHashSet[int]()
		s1.Add(1, 2)
		s2.Add(1, 2, 3)
		s3.Add(1, 2, 3, 4)
		t.Assert(s1.IsSubsetOf(s2), true)
		t.Assert(s2.IsSubsetOf(s3), true)
		t.Assert(s1.IsSubsetOf(s3), true)
		t.Assert(s2.IsSubsetOf(s1), false)
		t.Assert(s3.IsSubsetOf(s2), false)

		s4 := s1
		t.Assert(s1.IsSubsetOf(s4), true)
	})
}

func TestHashSet_Union(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[int]()
		s2 := g.NewHashSet[int]()
		s1.Add(1, 2)
		s2.Add(3, 4)
		s3 := s1.Union(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), true)
		t.Assert(s3.Contains(4), true)
	})
}

func TestHashSet_Diff(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[int]()
		s2 := g.NewHashSet[int]()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Diff(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), false)
		t.Assert(s3.Contains(4), false)

		s4 := s1
		s5 := s1.Diff(s2, s4)
		t.Assert(s5.Contains(1), true)
		t.Assert(s5.Contains(2), true)
		t.Assert(s5.Contains(3), false)
		t.Assert(s5.Contains(4), false)
	})
}

func TestHashSet_Intersect(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[int]()
		s2 := g.NewHashSet[int]()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Intersect(s2)
		t.Assert(s3.Contains(1), false)
		t.Assert(s3.Contains(2), false)
		t.Assert(s3.Contains(3), true)
		t.Assert(s3.Contains(4), false)
	})
}

func TestHashSet_Complement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[int]()
		s2 := g.NewHashSet[int]()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Complement(s2)
		t.Assert(s3.Contains(1), false)
		t.Assert(s3.Contains(2), false)
		t.Assert(s3.Contains(4), true)
		t.Assert(s3.Contains(5), true)
	})
}

func TestNewFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSetFrom[string]([]string{"a"})
		s2 := g.NewHashSetFrom[string]([]string{"b"}, false)
		s3 := g.NewHashSetFrom[string]([]string{"3"}, true)
		s4 := g.NewHashSetFrom[string]([]string{"s1", "s2"}, true)
		t.Assert(s1.Contains("a"), true)
		t.Assert(s2.Contains("b"), true)
		t.Assert(s3.Contains("3"), true)
		t.Assert(s4.Contains("s1"), true)
		t.Assert(s4.Contains("s3"), false)

	})
}

func TestNew(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string]()
		s1.Add("a", "2")
		s2 := g.NewHashSet[string](true)
		s2.Add("b", "3")
		t.Assert(s1.Contains("a"), true)

	})
}

func TestHashSet_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string](true)
		s1.Add("a", "a1", "b", "c")
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, "a1"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string](true)
		s1.Add("a", `"b"`, `\c`)
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, `"b"`), true)
		t.Assert(strings.Contains(str1, `\c`), true)
		t.Assert(strings.Contains(str1, `a`), true)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := g.HashSet[int]{}
		t.Assert(s1.Join(","), "")
	})
}

func TestHashSet_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string](true)
		s1.Add("a", "a2", "b", "c")
		str1 := s1.String()
		t.Assert(strings.Contains(str1, "["), true)
		t.Assert(strings.Contains(str1, "]"), true)
		t.Assert(strings.Contains(str1, "a2"), true)

		s1 = nil
		t.Assert(s1.String(), "")

		s2 := g.NewHashSet[int]()
		s2.Add(1)
		t.Assert(s2.String(), "[1]")
	})
}

func TestHashSet_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string](true)
		s2 := g.NewHashSet[string](true)
		s1.Add("a", "a2", "b", "c")
		s2.Add("b", "b1", "e", "f")
		ss := s1.Merge(s2)
		t.Assert(ss.Contains("a2"), true)
		t.Assert(ss.Contains("b1"), true)

	})
}

func TestHashSet_Sum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[int](true)
		s1.Add(1, 2, 3, 4)
		t.Assert(s1.Sum(), int(10))

	})
}

func TestHashSet_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[int](true)
		t.Assert(s.Pop(), 0)
		s.Add(1, 2, 3, 4)
		t.Assert(s.Size(), 4)
		t.AssertIN(s.Pop(), []int{1, 2, 3, 4})
		t.Assert(s.Size(), 3)
	})
}

func TestHashSet_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[int](true)
		s.Add(1, 2, 3, 4)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(0), nil)
		t.AssertIN(s.Pops(1), []int{1, 2, 3, 4})
		t.Assert(s.Size(), 3)
		a := s.Pops(6)
		t.Assert(len(a), 3)
		t.AssertIN(a, []int{1, 2, 3, 4})
		t.Assert(s.Size(), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[int](true)
		a := []int{1, 2, 3, 4}
		s.Add(a...)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(-2), nil)
		t.AssertIN(s.Pops(-1), a)
	})
}

func TestHashSet_Json(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		a1 := g.NewHashSetFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(len(b1), len(b2))
		t.Assert(err1, err2)

		a2 := g.NewHashSet[string]()
		err2 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(err2, nil)
		t.Assert(a2.Contains("a"), true)
		t.Assert(a2.Contains("b"), true)
		t.Assert(a2.Contains("c"), true)
		t.Assert(a2.Contains("d"), true)
		t.Assert(a2.Contains("e"), false)

		var a3 g.HashSet[string]
		err := json.UnmarshalUseNumber(b2, &a3)
		t.AssertNil(err)
		t.Assert(a3.Contains("a"), true)
		t.Assert(a3.Contains("b"), true)
		t.Assert(a3.Contains("c"), true)
		t.Assert(a3.Contains("d"), true)
		t.Assert(a3.Contains("e"), false)
	})
}

func TestHashSet_Add(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[int](true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Add(1), false)
		t.Assert(s.Add(2), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.Add(2), false)
		t.Assert(s.Contains(2), true)
	})
}

func TestHashSet_Walk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var set g.HashSet[int]
		set.Add([]int{1, 2}...)
		set.Walk(func(item int) int {
			return gconv.Int(item) + 10
		})
		t.Assert(set.Size(), 2)
		t.Assert(set.Contains(11), true)
		t.Assert(set.Contains(12), true)
	})
}

func TestHashSet_AddIfNotExistFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[int](true)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			r := s.Add(1)
			t.Assert(r, true)
		}()
		time.Sleep(100 * time.Millisecond)
		go func() {
			defer wg.Done()
			r := s.Add(1)
			t.Assert(r, false)
		}()
		wg.Wait()
	})
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[*exampleElement](true)
		t.Assert(s.Add(nil), false)
		s1 := g.HashSet[int]{}
		t.Assert(s1.Add(1), true)
	})
}

func TestHashSet_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Set  *g.HashSet[string]
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"set":  []byte(`["k1","k2","k3"]`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains("k1"), true)
		t.Assert(v.Set.Contains("k2"), true)
		t.Assert(v.Set.Contains("k3"), true)
		t.Assert(v.Set.Contains("k4"), false)
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"set":  []string{"k1", "k2", "k3"},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains("k1"), true)
		t.Assert(v.Set.Contains("k2"), true)
		t.Assert(v.Set.Contains("k3"), true)
		t.Assert(v.Set.Contains("k4"), false)
	})
}

func TestHashSet_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		set := g.NewHashSet[int]()
		set.Add(1, 2, 3)

		copySet := set.DeepCopy().(*g.HashSet[int])
		copySet.Add(4)
		t.AssertNE(set.Size(), copySet.Size())
		t.AssertNE(set.String(), copySet.String())

		set = nil
		t.AssertNil(set.DeepCopy())
	})
}
