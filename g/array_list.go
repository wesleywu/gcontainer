// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
	"bytes"
	json2 "encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/wesleywu/gcontainer/internal/deepcopy"
	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/internal/rwmutex"
	"github.com/wesleywu/gcontainer/utils/empty"
	"github.com/wesleywu/gcontainer/utils/gconv"
	"github.com/wesleywu/gcontainer/utils/grand"
	"github.com/wesleywu/gcontainer/utils/gstr"
)

// ArrayList is a golang array with rich features.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type ArrayList[T comparable] struct {
	mu    rwmutex.RWMutex
	array []T
}

// NewArrayList creates and returns an empty array.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewArrayList[T comparable](safe ...bool) *ArrayList[T] {
	return NewArrayListSize[T](0, 0, safe...)
}

// NewArrayListSize create and returns an array with given size and cap.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewArrayListSize[T comparable](size int, cap int, safe ...bool) *ArrayList[T] {
	return &ArrayList[T]{
		mu:    rwmutex.Create(safe...),
		array: make([]T, size, cap),
	}
}

// NewArrayListRange creates and returns an array by a range from `start` to `end`
// with step value `step`.
func NewArrayListRange(start, end, step int, safe ...bool) *ArrayList[int] {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]int, 0)
	index := 0
	for i := start; i <= end; i += step {
		slice = append(slice, i)
		index++
	}
	return NewArrayListFrom[int](slice, safe...)
}

// NewArrayListFrom is alias of NewArrayListFrom.
// See NewArrayListFrom.
func NewArrayListFrom[T comparable](array []T, safe ...bool) *ArrayList[T] {
	return &ArrayList[T]{
		mu:    rwmutex.Create(safe...),
		array: array,
	}
}

// NewArrayListFromCopy is alias of NewArrayFromCopy.
// See NewArrayFromCopy.
func NewArrayListFromCopy[T comparable](array []T, safe ...bool) *ArrayList[T] {
	newArray := make([]T, len(array))
	copy(newArray, array)
	return &ArrayList[T]{
		mu:    rwmutex.Create(safe...),
		array: newArray,
	}
}

// MustGet returns the value by the specified index.
// If the given `index` is out of range of the array, it returns empty value of type T.
func (a *ArrayList[T]) MustGet(index int) (value T) {
	value, _ = a.Get(index)
	return
}

// Get returns the value by the specified index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *ArrayList[T]) Get(index int) (value T, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index < 0 || index >= len(a.array) {
		found = false
		return
	}
	return a.array[index], true
}

// Set sets value to specified index.
func (a *ArrayList[T]) Set(index int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(a.array)))
	}
	a.array[index] = value
	return nil
}

// Sum returns the sum of values in an array.
func (a *ArrayList[T]) Sum() (sum int) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		sum += gconv.Int(v)
	}
	return
}

// Sort sorts the array by custom function `less`.
func (a *ArrayList[T]) Sort(less func(v1, v2 T) bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	sort.Slice(a.array, func(i, j int) bool {
		return less(a.array[i], a.array[j])
	})
}

// InsertBefore inserts the `values` to the front of `index`.
func (a *ArrayList[T]) InsertBefore(index int, values ...T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(a.array)))
	}
	rear := append([]T{}, a.array[index:]...)
	a.array = append(a.array[0:index], values...)
	a.array = append(a.array, rear...)
	return nil
}

// InsertAfter inserts the `values` to the back of `index`.
func (a *ArrayList[T]) InsertAfter(index int, values ...T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(a.array)))
	}
	rear := append([]T{}, a.array[index+1:]...)
	a.array = append(a.array[0:index+1], values...)
	a.array = append(a.array, rear...)
	return nil
}

// RemoveAt removes an item by index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *ArrayList[T]) RemoveAt(index int) (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doRemoveWithoutLock(index)
}

// doRemoveWithoutLock removes an item by index without lock.
func (a *ArrayList[T]) doRemoveWithoutLock(index int) (value T, found bool) {
	if index < 0 || index >= len(a.array) {
		found = false
		return
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		value := a.array[0]
		a.array = a.array[1:]
		return value, true
	} else if index == len(a.array)-1 {
		value := a.array[index]
		a.array = a.array[:index]
		return value, true
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	value = a.array[index]
	a.array = append(a.array[:index], a.array[index+1:]...)
	return value, true
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (a *ArrayList[T]) RemoveValue(value T) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if i := a.doSearchWithoutLock(value); i != -1 {
		a.doRemoveWithoutLock(i)
		return true
	}
	return false
}

// Remove removes multiple items by `values`.
func (a *ArrayList[T]) Remove(values ...T) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	changed := false
	for _, value := range values {
		if i := a.doSearchWithoutLock(value); i != -1 {
			_, found := a.doRemoveWithoutLock(i)
			if found {
				changed = true
			}
		}
	}
	return changed
}

// RemoveAll removes multiple items by `values`.
func (a *ArrayList[T]) RemoveAll(values Collection[T]) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	changed := false
	values.ForEach(func(value T) bool {
		if i := a.doSearchWithoutLock(value); i != -1 {
			_, found := a.doRemoveWithoutLock(i)
			if found {
				changed = true
			}
		}
		return true
	})
	return changed
}

// PushLeft pushes one or multiple items to the beginning of array.
func (a *ArrayList[T]) PushLeft(value ...T) List[T] {
	a.mu.Lock()
	a.array = append(value, a.array...)
	a.mu.Unlock()
	return a
}

// PushRight pushes one or multiple items to the end of array.
// It equals to Append.
func (a *ArrayList[T]) PushRight(value ...T) List[T] {
	a.mu.Lock()
	a.array = append(a.array, value...)
	a.mu.Unlock()
	return a
}

// PopRand randomly pops and return an item out of array.
// Note that if the array is empty, the `found` is false.
func (a *ArrayList[T]) PopRand() (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doRemoveWithoutLock(grand.Intn(len(a.array)))
}

// PopRands randomly pops and returns `size` items out of array.
func (a *ArrayList[T]) PopRands(size int) []T {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	if size >= len(a.array) {
		size = len(a.array)
	}
	array := make([]T, size)
	for i := 0; i < size; i++ {
		array[i], _ = a.doRemoveWithoutLock(grand.Intn(len(a.array)))
	}
	return array
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the `found` is false.
func (a *ArrayList[T]) PopLeft() (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) == 0 {
		found = false
		return
	}
	value = a.array[0]
	a.array = a.array[1:]
	return value, true
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the `found` is false.
func (a *ArrayList[T]) PopRight() (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	index := len(a.array) - 1
	if index < 0 {
		found = false
		return
	}
	value = a.array[index]
	a.array = a.array[:index]
	return value, true
}

// PopLefts pops and returns `size` items from the beginning of array.
func (a *ArrayList[T]) PopLefts(size int) []T {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	if size >= len(a.array) {
		array := a.array
		a.array = a.array[:0]
		return array
	}
	value := a.array[0:size]
	a.array = a.array[size:]
	return value
}

// PopRights pops and returns `size` items from the end of array.
func (a *ArrayList[T]) PopRights(size int) []T {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	index := len(a.array) - size
	if index <= 0 {
		array := a.array
		a.array = a.array[:0]
		return array
	}
	value := a.array[index:]
	a.array = a.array[:index]
	return value
}

// Range picks and returns items by range, like array[start:end].
// Notice, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// If `end` is negative, then the offset will start from the end of array.
// If `end` is omitted, then the sequence will have everything from start up
// until the end of the array.
func (a *ArrayList[T]) Range(start int, end ...int) []T {
	a.mu.RLock()
	defer a.mu.RUnlock()
	offsetEnd := len(a.array)
	if len(end) > 0 && end[0] < offsetEnd {
		offsetEnd = end[0]
	}
	if start > offsetEnd {
		return nil
	}
	if start < 0 {
		start = 0
	}
	array := ([]T)(nil)
	if a.mu.IsSafe() {
		array = make([]T, offsetEnd-start)
		copy(array, a.array[start:offsetEnd])
	} else {
		array = a.array[start:offsetEnd]
	}
	return array
}

// SubSlice returns a slice of elements from the array as specified
// by the `offset` and `size` parameters.
// If in concurrent safe usage, it returns a copy of the slice; else a pointer.
//
// If offset is non-negative, the sequence will start at that offset in the array.
// If offset is negative, the sequence will start that far from the end of the array.
//
// If length is given and is positive, then the sequence will have up to that many elements in it.
// If the array is shorter than the length, then only the available array elements will be present.
// If length is given and is negative then the sequence will stop that many elements from the end of the array.
// If it is omitted, then the sequence will have everything from offset up until the end of the array.
//
// Any possibility crossing the left border of array, it will fail.
func (a *ArrayList[T]) SubSlice(offset int, length ...int) []T {
	a.mu.RLock()
	defer a.mu.RUnlock()
	size := len(a.array)
	if len(length) > 0 {
		size = length[0]
	}
	if offset > len(a.array) {
		return nil
	}
	if offset < 0 {
		offset = len(a.array) + offset
		if offset < 0 {
			return nil
		}
	}
	if size < 0 {
		offset += size
		size = -size
		if offset < 0 {
			return nil
		}
	}
	end := offset + size
	if end > len(a.array) {
		end = len(a.array)
		size = len(a.array) - offset
	}
	if a.mu.IsSafe() {
		s := make([]T, size)
		copy(s, a.array[offset:])
		return s
	} else {
		return a.array[offset:end]
	}
}

// Add is alias of PushRight, please See PushRight.
func (a *ArrayList[T]) Add(values ...T) bool {
	a.PushRight(values...)
	return true
}

// AddAll adds all the elements in the specified collection to this collection.
// Returns true if this collection changed as a result of the call
func (a *ArrayList[T]) AddAll(values Collection[T]) bool {
	a.PushRight(values.Slice()...)
	return true
}

// Len returns the length of array.
func (a *ArrayList[T]) Len() int {
	return a.Size()
}

// Size returns the length of array.
func (a *ArrayList[T]) Size() int {
	a.mu.RLock()
	length := len(a.array)
	a.mu.RUnlock()
	return length
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *ArrayList[T]) Slice() []T {
	if a.mu.IsSafe() {
		a.mu.RLock()
		defer a.mu.RUnlock()
		array := make([]T, len(a.array))
		copy(array, a.array)
		return array
	} else {
		return a.array
	}
}

// Interfaces returns current array as []T.
func (a *ArrayList[T]) Interfaces() []T {
	return a.Slice()
}

// Clone returns a new array, which is a copy of current array.
func (a *ArrayList[T]) Clone() (newArray Collection[T]) {
	a.mu.RLock()
	array := make([]T, len(a.array))
	copy(array, a.array)
	a.mu.RUnlock()
	return NewArrayListFrom[T](array, a.mu.IsSafe())
}

// Clear deletes all items of current array.
func (a *ArrayList[T]) Clear() {
	a.mu.Lock()
	if len(a.array) > 0 {
		a.array = make([]T, 0)
	}
	a.mu.Unlock()
}

// Contains checks whether a value exists in the array.
func (a *ArrayList[T]) Contains(value T) bool {
	return a.Search(value) != -1
}

// ContainsAll checks whether a value exists in the array.
func (a *ArrayList[T]) ContainsAll(values Collection[T]) bool {
	values.ForEach(func(value T) bool {
		found := a.Search(value) != -1
		if !found {
			return false
		}
		return true
	})
	return true
}

func (a *ArrayList[T]) Equals(another Collection[T]) bool {
	if a == another {
		return true
	}
	var (
		ano *ArrayList[T]
		ok  bool
	)
	if ano, ok = another.(*ArrayList[T]); !ok {
		return false
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	ano.mu.RLock()
	defer ano.mu.RUnlock()
	if len(a.array) != len(ano.array) {
		return false
	}
	for index, value := range a.array {
		if value != ano.array[index] {
			return false
		}
	}
	return true
}

// ContainsI checks whether a value exists in the array with case-insensitively.
// Note that it internally iterates the whole array to do the comparison with case-insensitively.
func (a *ArrayList[T]) ContainsI(value T) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return false
	}
	if s, ok := any(value).(string); ok {
		for _, v := range a.array {
			if strings.EqualFold(any(v).(string), s) {
				return true
			}
		}
		return false
	}
	return a.Contains(value)
}

// Search searches array by `value`, returns the index of `value`,
// or returns -1 if not exists.
func (a *ArrayList[T]) Search(value T) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.doSearchWithoutLock(value)
}

func (a *ArrayList[T]) doSearchWithoutLock(value T) int {
	if len(a.array) == 0 {
		return -1
	}
	result := -1
	for index, v := range a.array {
		if v == value {
			result = index
			break
		}
	}
	return result
}

// Unique uniques the array, clear repeated items.
// Example: [1,1,2,3,2] -> [1,2,3]
func (a *ArrayList[T]) Unique() List[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) == 0 {
		return a
	}
	var (
		ok          bool
		temp        T
		uniqueSet   = make(map[T]struct{})
		uniqueArray = make([]T, 0, len(a.array))
	)
	for i := 0; i < len(a.array); i++ {
		temp = a.array[i]
		if _, ok = uniqueSet[temp]; ok {
			continue
		}
		uniqueSet[temp] = struct{}{}
		uniqueArray = append(uniqueArray, temp)
	}
	a.array = uniqueArray
	return a
}

// LockFunc locks writing by callback function `f`.
func (a *ArrayList[T]) LockFunc(f func(array []T)) {
	a.mu.Lock()
	defer a.mu.Unlock()
	f(a.array)
}

// RLockFunc locks reading by callback function `f`.
func (a *ArrayList[T]) RLockFunc(f func(array []T)) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	f(a.array)
}

// Fill fills an array with num entries of the value `value`,
// keys starting at the `startIndex` parameter.
func (a *ArrayList[T]) Fill(startIndex int, num int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if startIndex < 0 || startIndex > len(a.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", startIndex, len(a.array)))
	}
	for i := startIndex; i < startIndex+num; i++ {
		if i > len(a.array)-1 {
			a.array = append(a.array, value)
		} else {
			a.array[i] = value
		}
	}
	return nil
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by `size`.
// The last chunk may contain less than size elements.
func (a *ArrayList[T]) Chunk(size int) [][]T {
	if size < 1 {
		return nil
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	length := len(a.array)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]T
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, a.array[i*size:end])
		i++
	}
	return n
}

// Pad pads array to the specified length with `value`.
// If size is positive then the array is padded on the right, or negative on the left.
// If the absolute value of `size` is less than or equal to the length of the array
// then no padding takes place.
func (a *ArrayList[T]) Pad(size int, val T) List[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size == 0 || (size > 0 && size < len(a.array)) || (size < 0 && size > -len(a.array)) {
		return a
	}
	n := size
	if size < 0 {
		n = -size
	}
	n -= len(a.array)
	tmp := make([]T, n)
	for i := 0; i < n; i++ {
		tmp[i] = val
	}
	if size > 0 {
		a.array = append(a.array, tmp...)
	} else {
		a.array = append(tmp, a.array...)
	}
	return a
}

// Rand randomly returns one item from array(no deleting).
func (a *ArrayList[T]) Rand() (value T, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		found = false
		return
	}
	return a.array[grand.Intn(len(a.array))], true
}

// Rands randomly returns `size` items from array(no deleting).
func (a *ArrayList[T]) Rands(size int) []T {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	array := make([]T, size)
	for i := 0; i < size; i++ {
		array[i] = a.array[grand.Intn(len(a.array))]
	}
	return array
}

// Shuffle randomly shuffles the array.
func (a *ArrayList[T]) Shuffle() List[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, v := range grand.Perm(len(a.array)) {
		a.array[i], a.array[v] = a.array[v], a.array[i]
	}
	return a
}

// Reverse makes array with elements in reverse order.
func (a *ArrayList[T]) Reverse() List[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, j := 0, len(a.array)-1; i < j; i, j = i+1, j-1 {
		a.array[i], a.array[j] = a.array[j], a.array[i]
	}
	return a
}

// Join joins array elements with a string `glue`.
func (a *ArrayList[T]) Join(glue string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	for k, v := range a.array {
		buffer.WriteString(gconv.String(v))
		if k != len(a.array)-1 {
			buffer.WriteString(glue)
		}
	}
	return buffer.String()
}

// CountValues counts the number of occurrences of all values in the array.
func (a *ArrayList[T]) CountValues() map[T]int {
	m := make(map[T]int)
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		m[v]++
	}
	return m
}

// ForEach iterates all elements in this collection readonly with custom callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *ArrayList[T]) ForEach(f func(value T) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		if !f(v) {
			break
		}
	}
}

// ForEachAsc iterates the array readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *ArrayList[T]) ForEachAsc(f func(index int, value T) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for k, v := range a.array {
		if !f(k, v) {
			break
		}
	}
}

// ForEachDesc iterates the array readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (a *ArrayList[T]) ForEachDesc(f func(k int, v T) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for i := len(a.array) - 1; i >= 0; i-- {
		if !f(i, a.array[i]) {
			break
		}
	}
}

// String returns current array as a string, which implements like json.Marshal does.
func (a *ArrayList[T]) String() string {
	if a == nil {
		return ""
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('[')
	s := ""
	for k, v := range a.array {
		s = gconv.String(v)
		if gstr.IsNumeric(s) {
			buffer.WriteString(s)
		} else {
			buffer.WriteString(`"` + gstr.QuoteMeta(s, `"\`) + `"`)
		}
		if k != len(a.array)-1 {
			buffer.WriteByte(',')
		}
	}
	buffer.WriteByte(']')
	return buffer.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (a ArrayList[T]) MarshalJSON() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return json.Marshal(a.array)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (a *ArrayList[T]) UnmarshalJSON(b []byte) error {
	if a.array == nil {
		a.array = make([]T, 0)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := json.UnmarshalUseNumber(b, &a.array); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (a *ArrayList[T]) UnmarshalValue(value interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	switch value.(type) {
	case string, []byte, json2.Number:
		return json.UnmarshalUseNumber(gconv.Bytes(value), &a.array)
	default:
		a.array = gconv.SliceAny[T](value)
	}
	return nil
}

// Filter iterates array and filters elements using custom callback function.
// It removes the element from array if callback function `filter` returns true,
// it or else does nothing and continues iterating.
func (a *ArrayList[T]) Filter(filter func(index int, value T) bool) List[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := 0; i < len(a.array); {
		if filter(i, a.array[i]) {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			i++
		}
	}
	return a
}

// FilterNil removes all nil value of the array.
func (a *ArrayList[T]) FilterNil() List[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := 0; i < len(a.array); {
		if empty.IsNil(a.array[i]) {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			i++
		}
	}
	return a
}

// FilterEmpty removes all empty value of the array.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (a *ArrayList[T]) FilterEmpty() List[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := 0; i < len(a.array); {
		if empty.IsEmpty(a.array[i]) {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			i++
		}
	}
	return a
}

// Walk applies a user supplied function `f` to every item of array.
func (a *ArrayList[T]) Walk(f func(value T) T) List[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, v := range a.array {
		a.array[i] = f(v)
	}
	return a
}

// IsEmpty checks whether the array is empty.
func (a *ArrayList[T]) IsEmpty() bool {
	return a.Len() == 0
}

// DeepCopy implements interface for deep copy of current type.
func (a *ArrayList[T]) DeepCopy() Collection[T] {
	if a == nil {
		return nil
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	newSlice := make([]T, len(a.array))
	for i, v := range a.array {
		newSlice[i] = deepcopy.Copy(v).(T)
	}
	return NewArrayListFrom[T](newSlice, a.mu.IsSafe())
}
