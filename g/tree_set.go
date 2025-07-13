package g

import (
	"bytes"

	"github.com/wesleywu/gcontainer/internal/deepcopy"
	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/internal/rwmutex"
	"github.com/wesleywu/gcontainer/utils/comparators"
	"github.com/wesleywu/gcontainer/utils/gconv"
	"github.com/wesleywu/gcontainer/utils/gstr"
)

// TreeSet is a golang sorted set with rich features.
// It is using increasing order in default, which can be changed by
// setting it a custom comparators.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type TreeSet[T comparable] struct {
	mu   rwmutex.RWMutex
	tree *TreeMap[T, struct{}]
}

// NewTreeSet creates and returns an empty sorted set.
// The parameter `safe` is used to specify whether using array in concurrent-safety, which is false in default.
// The parameter `comparators` used to compare values to sort in array,
// if it returns value < 0, means `a` < `b`; the `a` will be inserted before `b`;
// if it returns value = 0, means `a` = `b`; the `a` will be replaced by     `b`;
// if it returns value > 0, means `a` > `b`; the `a` will be inserted after  `b`;
func NewTreeSet[T comparable](comparator comparators.Comparator[T], safe ...bool) *TreeSet[T] {
	return &TreeSet[T]{
		mu:   rwmutex.Create(safe...),
		tree: NewTreeMap[T, struct{}](comparator, false),
	}
}

// NewTreeSetDefault creates and returns an empty sorted set using default comparators.
// The parameter `safe` is used to specify whether using array in concurrent-safety, which is false in default.
// if it returns value < 0, means `a` < `b`; the `a` will be inserted before `b`;
// if it returns value = 0, means `a` = `b`; the `a` will be replaced by     `b`;
// if it returns value > 0, means `a` > `b`; the `a` will be inserted after  `b`;
func NewTreeSetDefault[T comparable](safe ...bool) *TreeSet[T] {
	return &TreeSet[T]{
		mu:   rwmutex.Create(safe...),
		tree: NewTreeMap[T, struct{}](comparators.ComparatorAny[T], false),
	}
}

// NewTreeSetFrom creates and returns an sorted array with given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewTreeSetFrom[T comparable](elements []T, comparator comparators.Comparator[T], safe ...bool) *TreeSet[T] {
	a := &TreeSet[T]{
		mu:   rwmutex.Create(safe...),
		tree: NewTreeMap[T, struct{}](comparator, false),
	}
	for _, value := range elements {
		a.tree.PutIfAbsent(value, struct{}{})
	}
	return a
}

func (t *TreeSet[T]) lazyInit() {
	if t.tree == nil {
		t.tree = NewTreeMap[T, struct{}](comparators.ComparatorAny[T], false)
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

func (t *TreeSet[T]) AddAll(elements Collection[T]) bool {
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

func (t *TreeSet[T]) Ceiling(element T) (ceiling T, found bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	if ceilingNode := t.tree.CeilingEntry(element); ceilingNode != nil {
		return ceilingNode.Key(), true
	}
	return ceiling, false
}

func (t *TreeSet[T]) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	t.tree.Clear()
}

func (t *TreeSet[T]) Clone() Collection[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()
	newTree := t.tree.Clone(false)
	return &TreeSet[T]{
		mu:   rwmutex.Create(t.mu.IsSafe()),
		tree: newTree.(*TreeMap[T, struct{}]),
	}
}

func (t *TreeSet[T]) Comparator() comparators.Comparator[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	return t.tree.Comparator()
}

func (t *TreeSet[T]) Contains(element T) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	return t.tree.ContainsKey(element)
}

func (t *TreeSet[T]) ContainsAll(elements Collection[T]) bool {
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

func (t *TreeSet[T]) DeepCopy() Collection[T] {
	if t == nil {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	data := make([]T, 0)
	t.tree.ForEach(func(k T, _ struct{}) bool {
		data = append(data, deepcopy.Copy(k).(T))
		return true
	})
	return NewTreeSetFrom[T](data, t.Comparator(), t.mu.IsSafe())
}

// Equals checks whether the two sets equal.
func (t *TreeSet[T]) Equals(another Collection[T]) bool {
	if t == another {
		return true
	}
	var (
		ano *TreeSet[T]
		ok  bool
	)
	if ano, ok = another.(*TreeSet[T]); !ok {
		return false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	ano.mu.RLock()
	defer ano.mu.RUnlock()
	if t.tree.Size() != ano.tree.Size() {
		return false
	}
	values := t.tree.Map()
	valuesAno := ano.tree.Map()
	for k, v := range values {
		vAno, vOk := valuesAno[k]
		if !vOk {
			return false
		}
		if v != vAno {
			return false
		}
	}
	return true
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
	return first.Key(), true
}

func (t *TreeSet[T]) Floor(element T) (floor T, found bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	if floorNode := t.tree.FloorEntry(element); floorNode != nil {
		return floorNode.Key(), true
	}
	return floor, false
}

func (t *TreeSet[T]) ForEach(f func(T) bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	t.tree.ForEach(func(key T, _ struct{}) bool {
		return f(key)
	})
}

func (t *TreeSet[T]) ForEachDescending(f func(T) bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	t.tree.ForEachDesc(func(key T, _ struct{}) bool {
		return f(key)
	})
}

func (t *TreeSet[T]) HeadSet(toElement T, inclusive bool) SortedSet[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	result := NewTreeSet(t.tree.Comparator(), t.mu.IsSafe())

	t.tree.IteratorDescFrom(toElement, inclusive, func(key T, _ struct{}) bool {
		result.Add(key)
		return true
	})
	return result
}

func (t *TreeSet[T]) Higher(element T) (higher T, found bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	if higherNode := t.tree.HigherEntry(element); higherNode != nil {
		return higherNode.Key(), true
	}
	return higher, false
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
	t.tree.ForEach(func(key T, value struct{}) bool {
		buffer.WriteString(gconv.String(key))
		if i != size-1 {
			buffer.WriteString(glue)
		}
		i++
		return true
	})
	return buffer.String()
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
	return last.Key(), true
}

func (t *TreeSet[T]) Lower(element T) (lower T, found bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	if lowerNode := t.tree.LowerEntry(element); lowerNode != nil {
		return lowerNode.Key(), true
	}
	return lower, false
}

func (t *TreeSet[T]) PollFirst() (first T, found bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	firstNode := t.tree.PollFirstEntry()
	if firstNode != nil {
		return firstNode.Key(), true
	}
	return first, false
}

func (t *TreeSet[T]) PollHeadSet(toElement T, inclusive bool) SortedSet[T] {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	result := NewTreeSet(t.Comparator(), t.mu.IsSafe())
	headKeys := t.tree.HeadMap(toElement, inclusive).Keys()
	for _, key := range headKeys {
		t.tree.Remove(key)
		result.Add(key)
	}
	return result
}

func (t *TreeSet[T]) PollLast() (last T, found bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	lastNode := t.tree.PollLastEntry()
	if lastNode != nil {
		return lastNode.Key(), true
	}
	return last, false
}

func (t *TreeSet[T]) PollTailSet(fromElement T, inclusive bool) SortedSet[T] {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	result := NewTreeSet(t.Comparator(), t.mu.IsSafe())
	tailKeys := t.tree.TailMap(fromElement, inclusive).Keys()
	for _, key := range tailKeys {
		t.tree.Remove(key)
		result.Add(key)
	}
	return result
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

func (t *TreeSet[T]) RemoveAll(elements Collection[T]) bool {
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
	t.tree.ForEach(func(key T, _ struct{}) bool {
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

func (t *TreeSet[T]) SubSet(fromElement T, fromInclusive bool, toElement T, toInclusive bool) SortedSet[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	subKeys := t.tree.SubMap(fromElement, fromInclusive, toElement, toInclusive).Keys()
	return NewTreeSetFrom[T](subKeys, t.tree.Comparator(), t.mu.IsSafe())
}

func (t *TreeSet[T]) TailSet(fromElement T, inclusive bool) SortedSet[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	result := NewTreeSet(t.tree.Comparator(), t.mu.IsSafe())

	t.tree.IteratorAscFrom(fromElement, inclusive, func(key T, _ struct{}) bool {
		result.Add(key)
		return true
	})
	return result
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (t TreeSet[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Slice())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (t *TreeSet[T]) UnmarshalJSON(b []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	var array []T
	if err := json.UnmarshalUseNumber(b, &array); err != nil {
		return err
	}
	for _, v := range array {
		t.tree.Put(v, struct{}{})
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for set.
func (t *TreeSet[T]) UnmarshalValue(value interface{}) (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	var array []T
	switch value.(type) {
	case string, []byte:
		err = json.UnmarshalUseNumber(gconv.Bytes(value), &array)
	default:
		array = gconv.SliceAny[T](value)
	}
	for _, v := range array {
		t.tree.Put(v, struct{}{})
	}
	return
}

// Diff returns a new set which is the difference set from `t` to `others`.
// Which means, all the items in `newSet` are in `t` but not in any of the `others`.
func (t *TreeSet[T]) Diff(others ...Set[T]) Set[T] {
	newSet := NewTreeSet(t.Comparator(), t.mu.IsSafe())
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	for _, other := range others {
		if t == other {
			return newSet // 差集为空
		}
	}
	for _, k := range t.tree.Keys() {
		found := false
		for _, other := range others {
			if other.Contains(k) {
				found = true
				break
			}
		}
		if !found {
			newSet.Add(k)
		}
	}
	return newSet
}

// Complement returns a new set which is the complement from `t` to `full`.
// Which means, all the items in `newSet` are in `full` and not in `t`.
func (t *TreeSet[T]) Complement(full Set[T]) Set[T] {
	newSet := NewTreeSet(t.Comparator(), t.mu.IsSafe())
	if t == full {
		return newSet // 补集为空，避免死锁
	}
	full.ForEach(func(v T) bool {
		if !t.Contains(v) {
			newSet.Add(v)
		}
		return true
	})
	return newSet
}

// Merge merges the given sets into the current set (in-place).
// This method modifies the current set by adding all elements from the given sets.
// Returns the current set after modification.
func (t *TreeSet[T]) Merge(others ...Set[T]) Set[T] {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lazyInit()
	for _, other := range others {
		if t == other {
			continue
		}
		other.ForEach(func(v T) bool {
			t.tree.PutIfAbsent(v, struct{}{})
			return true
		})
	}
	return t
}

// Intersect returns a new set which is the intersection from `t` to `others`.
// Which means, all the items in `newSet` are in `t` and also in all the `others`.
func (t *TreeSet[T]) Intersect(others ...Set[T]) Set[T] {
	newSet := NewTreeSet(t.Comparator(), t.mu.IsSafe())
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.lazyInit()
	for _, k := range t.tree.Keys() {
		foundInAll := true
		for _, other := range others {
			if t == other {
				continue // 避免死锁
			}
			if !other.Contains(k) {
				foundInAll = false
				break
			}
		}
		if foundInAll {
			newSet.Add(k)
		}
	}
	return newSet
}

// Union returns a new set which is the union of `t` and `others`.
// Which means, all the items in `newSet` are in `t` or in `others`.
func (t *TreeSet[T]) Union(others ...Set[T]) Set[T] {
	newSet := NewTreeSet(t.Comparator(), t.mu.IsSafe())
	t.ForEach(func(v T) bool {
		newSet.Add(v)
		return true
	})
	for _, other := range others {
		if t == other {
			continue // 跳过自身
		}
		other.ForEach(func(v T) bool {
			newSet.Add(v)
			return true
		})
	}
	return newSet
}

// IsSubsetOf returns true if the current set is a subset of the given set.
func (t *TreeSet[T]) IsSubsetOf(other Set[T]) bool {
	if t == other {
		return true
	}
	isSubset := true
	t.ForEach(func(v T) bool {
		if !other.Contains(v) {
			isSubset = false
			return false // 终止遍历
		}
		return true
	})
	return isSubset
}

// IsSupersetOf returns true if the current set is a superset of the given set.
func (t *TreeSet[T]) IsSupersetOf(other Set[T]) bool {
	if t == other {
		return true
	}
	return other.IsSubsetOf(t)
}

// IsDisjointWith returns true if the current set has no elements in common with the given set.
func (t *TreeSet[T]) IsDisjointWith(other Set[T]) bool {
	if t == other {
		return t.IsEmpty()
	}
	isDisjoint := true
	t.ForEach(func(v T) bool {
		if other.Contains(v) {
			isDisjoint = false
			return false // 终止遍历
		}
		return true
	})
	return isDisjoint
}
