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

func TestStrSet_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s g.HashSet[string]
		s.Add("1", "1", "2")
		s.Add([]string{"3", "4"}...)
		t.Assert(s.Size(), 4)
		t.AssertIN("1", s.Slice())
		t.AssertIN("2", s.Slice())
		t.AssertIN("3", s.Slice())
		t.AssertIN("4", s.Slice())
		t.AssertNI("0", s.Slice())
		t.Assert(s.Contains("4"), true)
		t.Assert(s.Contains("5"), false)
		s.Remove("1")
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestStrSet_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[string]()
		s.Add("1", "1", "2")
		s.Add([]string{"3", "4"}...)
		t.Assert(s.Size(), 4)
		t.AssertIN("1", s.Slice())
		t.AssertIN("2", s.Slice())
		t.AssertIN("3", s.Slice())
		t.AssertIN("4", s.Slice())
		t.AssertNI("0", s.Slice())
		t.Assert(s.Contains("4"), true)
		t.Assert(s.Contains("5"), false)
		s.Remove("1")
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestStrSet_ContainsI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[string]()
		s.Add("a", "b", "C")
		t.Assert(s.Contains("A"), false)
		t.Assert(s.Contains("a"), true)
		t.Assert(s.ContainsI("A"), true)
		t.Assert(s.ContainsI("d"), false)
	})
}

func TestStrSet_Iterator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[string]()
		s.Add("1", "2", "3")
		t.Assert(s.Size(), 3)

		a1 := g.NewHashSet[string](true)
		a2 := g.NewHashSet[string](true)
		s.ForEach(func(v string) bool {
			a1.Add(v)
			return false
		})
		s.ForEach(func(v string) bool {
			a2.Add(v)
			return true
		})
		t.Assert(a1.Size(), 1)
		t.Assert(a2.Size(), 3)
	})
}

func TestStrSet_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[string]()
		s.Add("1", "2", "3")
		t.Assert(s.Size(), 3)
		s.LockFunc(func(m map[string]struct{}) {
			delete(m, "1")
		})
		t.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[string]struct{}) {
			t.Assert(m, map[string]struct{}{
				"3": {},
				"2": {},
			})
		})
	})
}

func TestStrSet_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string]()
		s2 := g.NewHashSet[string]()
		s3 := g.NewHashSet[string]()
		s4 := g.NewHashSet[string]()
		s1.Add("1", "2", "3")
		s2.Add("1", "2", "3")
		s3.Add("1", "2", "3", "4")
		s4.Add("4", "5", "6")
		t.Assert(s1.Equals(s2), true)
		t.Assert(s1.Equals(s3), false)
		t.Assert(s1.Equals(s4), false)
		s5 := s1
		t.Assert(s1.Equals(s5), true)
	})
}

func TestStrSet_IsSubsetOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string]()
		s2 := g.NewHashSet[string]()
		s3 := g.NewHashSet[string]()
		s1.Add("1", "2")
		s2.Add("1", "2", "3")
		s3.Add("1", "2", "3", "4")
		t.Assert(s1.IsSubsetOf(s2), true)
		t.Assert(s2.IsSubsetOf(s3), true)
		t.Assert(s1.IsSubsetOf(s3), true)
		t.Assert(s2.IsSubsetOf(s1), false)
		t.Assert(s3.IsSubsetOf(s2), false)

		s4 := s1
		t.Assert(s1.IsSubsetOf(s4), true)
	})
}

func TestStrSet_Union(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string]()
		s2 := g.NewHashSet[string]()
		s1.Add("1", "2")
		s2.Add("3", "4")
		s3 := s1.Union(s2)
		t.Assert(s3.Contains("1"), true)
		t.Assert(s3.Contains("2"), true)
		t.Assert(s3.Contains("3"), true)
		t.Assert(s3.Contains("4"), true)
	})
}

func TestStrSet_Diff(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string]()
		s2 := g.NewHashSet[string]()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Diff(s2)
		t.Assert(s3.Contains("1"), true)
		t.Assert(s3.Contains("2"), true)
		t.Assert(s3.Contains("3"), false)
		t.Assert(s3.Contains("4"), false)

		s4 := s1
		s5 := s1.Diff(s2, s4)
		t.Assert(s5.Contains("1"), false)
		t.Assert(s5.Contains("2"), false)
		t.Assert(s5.Contains("3"), false)
		t.Assert(s5.Contains("4"), false)
	})
}

func TestStrSet_Intersect(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string]()
		s2 := g.NewHashSet[string]()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Intersect(s2)
		t.Assert(s3.Contains("1"), false)
		t.Assert(s3.Contains("2"), false)
		t.Assert(s3.Contains("3"), true)
		t.Assert(s3.Contains("4"), false)
	})
}

func TestStrSet_Complement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string]()
		s2 := g.NewHashSet[string]()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Complement(s2)
		t.Assert(s3.Contains("1"), false)
		t.Assert(s3.Contains("2"), false)
		t.Assert(s3.Contains("4"), true)
		t.Assert(s3.Contains("5"), true)
	})
}

func TestNewIntSetFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewArrayListFrom[int]([]int{1, 2, 3, 4})
		s2 := g.NewArrayListFrom[int]([]int{5, 6, 7, 8})
		t.Assert(s1.Contains(3), true)
		t.Assert(s1.Contains(5), false)
		t.Assert(s2.Contains(3), false)
		t.Assert(s2.Contains(5), true)
	})
}

func TestStrSet_Merge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string]()
		s2 := g.NewHashSet[string]()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Merge(s2)
		t.Assert(s3.Contains("1"), true)
		t.Assert(s3.Contains("6"), false)
		t.Assert(s3.Contains("4"), true)
		t.Assert(s3.Contains("5"), true)
	})
}

func TestNewStrSetFrom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewArrayListFrom[string]([]string{"a", "b", "c"}, true)
		t.Assert(s1.Contains("b"), true)
		t.Assert(s1.Contains("d"), false)
	})
}

func TestStrSet_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewArrayListFrom[string]([]string{"a", "b", "c"}, true)
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, "b"), true)
		t.Assert(strings.Contains(str1, "d"), false)
	})

	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string]()
		t.Assert(s1.Join(","), "")
		s1.Add("a", `"b"`, `\c`)
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, `"b"`), true)
		t.Assert(strings.Contains(str1, `\c`), true)
		t.Assert(strings.Contains(str1, `a`), true)
	})
}

func TestStrSet_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewArrayListFrom[string]([]string{"a", "b", "c"}, true)
		str1 := s1.String()
		t.Assert(strings.Contains(str1, "b"), true)
		t.Assert(strings.Contains(str1, "d"), false)
		s1 = nil
		t.Assert(s1.String(), "")
	})

	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewHashSet[string](true)
		s1.Add("a", "a2", "b", "c")
		str1 := s1.String()
		t.Assert(strings.Contains(str1, "["), true)
		t.Assert(strings.Contains(str1, "]"), true)
		t.Assert(strings.Contains(str1, "a2"), true)
	})
}

func TestStrSet_Sum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewArrayListFrom[string]([]string{"a", "b", "c"}, true)
		s2 := g.NewArrayListFrom[int]([]int{2, 3, 4}, true)
		t.Assert(s1.Sum(), 0)
		t.Assert(s2.Sum(), 9)
	})
}

func TestStrSet_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewArrayListFrom[string]([]string{"a", "b", "c"}, true)
		t.Assert(s1.Size(), 3)

	})
}

func TestStrSet_Remove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := g.NewArrayListFrom[string]([]string{"a", "b", "c"}, true)
		s1.Remove("b")
		t.Assert(s1.Contains("b"), false)
		t.Assert(s1.Contains("c"), true)
	})
}

func TestStrSet_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := []string{"a", "b", "c", "d"}
		s := g.NewHashSetFrom[string](a, true)
		t.Assert(s.Size(), 4)
		t.AssertIN(s.Pop(), a)
		t.Assert(s.Size(), 3)
		t.AssertIN(s.Pop(), a)
		t.Assert(s.Size(), 2)

		s1 := g.HashSet[string]{}
		t.Assert(s1.Pop(), "")
	})
}

func TestStrSet_Pops(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := []string{"a", "b", "c", "d"}
		s := g.NewHashSetFrom[string](a, true)
		array := s.Pops(2)
		t.Assert(len(array), 2)
		t.Assert(s.Size(), 2)
		t.AssertIN(array, a)
		t.Assert(s.Pops(0), nil)
		t.AssertIN(s.Pops(2), a)
		t.Assert(s.Size(), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[string](true)
		a := []string{"1", "2", "3", "4"}
		s.Add(a...)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(-2), nil)
		t.AssertIN(s.Pops(-1), a)
	})
}

func TestStrSet_AddIfNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[string](true)
		s.Add("1")
		t.Assert(s.Contains("1"), true)
		t.Assert(s.Add("1"), false)
		t.Assert(s.Add("2"), true)
		t.Assert(s.Contains("2"), true)
		t.Assert(s.Add("2"), false)
		t.Assert(s.Contains("2"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		s := g.HashSet[string]{}
		t.Assert(s.Add("1"), true)
	})
}

func TestStrSet_AddIfNotExistFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[string](true)
		s.Add("1")
		t.Assert(s.Contains("1"), true)
		t.Assert(s.Contains("2"), false)
		t.Assert(s.Contains("2"), false)
		t.Assert(s.Add("2"), true)
		t.Assert(s.Contains("2"), true)
		t.Assert(s.Add("2"), false)
		t.Assert(s.Contains("2"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[string](true)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := s.Add("1")
			t.Assert(r, false)
		}()
		s.Add("1")
		wg.Wait()
	})
	gtest.C(t, func(t *gtest.T) {
		s := g.HashSet[string]{}
		t.Assert(s.Add("1"), true)
	})
}

func TestStrSet_AddIfNotExistFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.NewHashSet[string](true)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			r := s.Add("1")
			t.Assert(r, true)
		}()
		time.Sleep(100 * time.Millisecond)
		go func() {
			defer wg.Done()
			r := s.Add("1")
			t.Assert(r, false)
		}()
		wg.Wait()
	})
	gtest.C(t, func(t *gtest.T) {
		s := g.HashSet[string]{}
		t.Assert(s.Add("1"), true)
	})
}

func TestStrSet_Json(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		a1 := g.NewArrayListFrom[string](s1)
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

func TestStrSet_Walk(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			set    g.HashSet[string]
			names  = []string{"user", "user_detail"}
			prefix = "gf_"
		)
		set.Add(names...)
		// Add prefix for given table names.
		set.Walk(func(item string) string {
			return prefix + item
		})
		t.Assert(set.Size(), 2)
		t.Assert(set.Contains("gf_user"), true)
		t.Assert(set.Contains("gf_user_detail"), true)
	})
}

func TestStrSet_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Set  *g.HashSet[string]
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"set":  []byte(`["1","2","3"]`),
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains("1"), true)
		t.Assert(v.Set.Contains("2"), true)
		t.Assert(v.Set.Contains("3"), true)
		t.Assert(v.Set.Contains("4"), false)
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"set":  []string{"1", "2", "3"},
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains("1"), true)
		t.Assert(v.Set.Contains("2"), true)
		t.Assert(v.Set.Contains("3"), true)
		t.Assert(v.Set.Contains("4"), false)
	})
}

func TestStrSet_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		set := g.NewHashSet[string]()
		set.Add("1", "2", "3")

		copySet := set.DeepCopy().(*g.HashSet[string])
		copySet.Add("4")
		t.AssertNE(set.Size(), copySet.Size())
		t.AssertNE(set.String(), copySet.String())

		set = nil
		t.AssertNil(set.DeepCopy())
	})
}
