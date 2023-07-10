package gset

import (
	"bytes"

	"github.com/wesleywu/gcontainer/garray"
	"github.com/wesleywu/gcontainer/gtree"
	"github.com/wesleywu/gcontainer/internal/deepcopy"
	"github.com/wesleywu/gcontainer/internal/rwmutex"
	"github.com/wesleywu/gcontainer/utils/comparator"
	"github.com/wesleywu/gcontainer/utils/gconv"
	"github.com/wesleywu/gcontainer/utils/gstr"
)

// TreeSet is a golang sorted set with rich features.
// It is using increasing order in default, which can be changed by
// setting it a custom comparator.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type TreeSet[T comparable] struct {
	mu   rwmutex.RWMutex
	tree *gtree.RedBlackTree[T, struct{}]
}

// NewTreeSet creates and returns an empty sorted set.
// The parameter `safe` is used to specify whether using array in concurrent-safety, which is false in default.
// The parameter `comparator` used to compare values to sort in array,
// if it returns value < 0, means `a` < `b`; the `a` will be inserted before `b`;
// if it returns value = 0, means `a` = `b`; the `a` will be replaced by     `b`;
// if it returns value > 0, means `a` > `b`; the `a` will be inserted after  `b`;
func NewTreeSet[T comparable](comparator comparator.Comparator[T], safe ...bool) *TreeSet[T] {
	return &TreeSet[T]{
		mu:   rwmutex.Create(safe...),
		tree: gtree.NewRedBlackTree[T, struct{}](comparator, false),
	}
}

// NewTreeSetDefault creates and returns an empty sorted set using default comparator.
// The parameter `safe` is used to specify whether using array in concurrent-safety, which is false in default.
// if it returns value < 0, means `a` < `b`; the `a` will be inserted before `b`;
// if it returns value = 0, means `a` = `b`; the `a` will be replaced by     `b`;
// if it returns value > 0, means `a` > `b`; the `a` will be inserted after  `b`;
func NewTreeSetDefault[T comparable](safe ...bool) *TreeSet[T] {
	return &TreeSet[T]{
		mu:   rwmutex.Create(safe...),
		tree: gtree.NewRedBlackTree[T, struct{}](comparator.ComparatorAny[T], false),
	}
}

// NewTreeSetFrom creates and returns an sorted array with given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewTreeSetFrom[T comparable](elements []T, comparator comparator.Comparator[T], safe ...bool) *TreeSet[T] {
	a := &TreeSet[T]{
		mu:   rwmutex.Create(safe...),
		tree: gtree.NewRedBlackTree[T, struct{}](comparator, false),
	}
	for _, value := range elements {
		a.tree.PutIfAbsent(value, struct{}{})
	}
	return a
}

func (t *TreeSet[T]) lazyInit() {
	if t.tree == nil {
		t.tree = gtree.NewRedBlackTree[T, struct{}](comparator.ComparatorAny[T], false)
	}
}

func (t *TreeSet[T]) Add(elements ...T) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	changed := false
	for _, value := range elements {
		putOk := t.tree.PutIfAbsent(value, struct{}{})
		if putOk {
			changed = putOk
		}
	}
	return changed
}

func (t *TreeSet[T]) AddAll(elements garray.Collection[T]) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	changed := false
	elements.ForEach(func(value T) bool {
		putOk := t.tree.PutIfAbsent(value, struct{}{})
		if putOk {
			changed = putOk
		}
		return true
	})
	return changed
}

func (t *TreeSet[T]) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	t.tree.Clear()
}

func (t *TreeSet[T]) Contains(element T) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	return t.tree.ContainsKey(element)
}

func (t *TreeSet[T]) ContainsAll(elements garray.Collection[T]) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	allFound := true
	elements.ForEach(func(value T) bool {
		if found := t.tree.ContainsKey(value); !found {
			allFound = false
			return false
		}
		return true
	})
	return allFound
}

func (t *TreeSet[T]) DeepCopy() garray.Collection[T] {
	if t == nil {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	data := make([]T, 0)
	t.tree.Iterator(func(k T, _ struct{}) bool {
		data = append(data, deepcopy.Copy(k).(T))
		return true
	})
	return NewTreeSetFrom[T](data, t.Comparator(), t.mu.IsSafe())
}

func (t *TreeSet[T]) ForEach(f func(T) bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	t.tree.Iterator(func(key T, _ struct{}) bool {
		return f(key)
	})
}

func (t *TreeSet[T]) ForEachDescending(f func(T) bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	t.tree.IteratorDesc(func(key T, _ struct{}) bool {
		return f(key)
	})
}

func (t *TreeSet[T]) IsEmpty() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	return t.tree.IsEmpty()
}

func (t *TreeSet[T]) Join(glue string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	if t.tree.Size() == 0 {
		return ""
	}
	size := t.tree.Size()
	if size == 0 {
		return ""
	}
	var (
		i      = 0
		buffer = bytes.NewBuffer(nil)
	)
	t.tree.Iterator(func(key T, value struct{}) bool {
		buffer.WriteString(gconv.String(key))
		if i != size-1 {
			buffer.WriteString(glue)
		}
		i++
		return true
	})
	return buffer.String()
}

func (t *TreeSet[T]) Remove(elements ...T) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	changed := false
	for _, value := range elements {
		if _, removed := t.tree.Remove(value); removed {
			changed = true
		}
	}
	return changed
}

func (t *TreeSet[T]) RemoveAll(elements garray.Collection[T]) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	changed := false
	elements.ForEach(func(value T) bool {
		if _, removed := t.tree.Remove(value); removed {
			changed = true
		}
		return true
	})
	return changed
}

func (t *TreeSet[T]) Size() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	return t.tree.Size()
}

func (t *TreeSet[T]) Slice() []T {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	return t.tree.Keys()
}

func (t *TreeSet[T]) String() string {
	if t == nil {
		return ""
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	size := t.tree.Size()
	if size == 0 {
		return "[]"
	}
	var (
		i      = 0
		buffer = bytes.NewBuffer(nil)
	)
	buffer.WriteByte('[')
	s := ""
	t.tree.Iterator(func(key T, _ struct{}) bool {
		s = gconv.String(key)
		if gstr.IsNumeric(s) {
			buffer.WriteString(s)
		} else {
			buffer.WriteString(`"` + gstr.QuoteMeta(s, `"\`) + `"`)
		}
		if i != size-1 {
			buffer.WriteByte(',')
		}
		i++
		return true
	})
	buffer.WriteByte(']')
	return buffer.String()
}

func (t *TreeSet[T]) Comparator() comparator.Comparator[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	return t.tree.Comparator()
}

func (t *TreeSet[T]) First() (element T, found bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	first := t.tree.Left()
	if first == nil {
		found = false
		return
	}
	return first.Key, true
}

func (t *TreeSet[T]) HeadSet(toElement T, inclusive bool) SortedSet[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	result := NewTreeSet[T](t.tree.Comparator(), t.mu.IsSafe())

	t.tree.IteratorDescFrom(toElement, inclusive, func(key T, _ struct{}) bool {
		result.Add(key)
		return true
	})
	return result
}

func (t *TreeSet[T]) Last() (element T, found bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	last := t.tree.Right()
	if last == nil {
		found = false
		return
	}
	return last.Key, true
}

func (t *TreeSet[T]) SubSet(fromElement T, fromInclusive bool, toElement T, toInclusive bool) SortedSet[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	subKeys := t.tree.SubMap(fromElement, fromInclusive, toElement, toInclusive).Keys()
	return NewTreeSetFrom(subKeys, t.tree.Comparator(), t.mu.IsSafe())
}

func (t *TreeSet[T]) TailSet(fromElement T, inclusive bool) SortedSet[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	result := NewTreeSet[T](t.tree.Comparator(), t.mu.IsSafe())

	t.tree.IteratorAscFrom(fromElement, inclusive, func(key T, _ struct{}) bool {
		result.Add(key)
		return true
	})
	return result
}
