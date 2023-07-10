// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package garray provides most commonly used array containers which also support concurrent-safe/unsafe switch feature.
package garray

type Collection[T comparable] interface {
	// Add adds all the elements in the specified slice to this collection.
	// Returns true if this collection changed as a result of the call
	Add(...T) bool

	// AddAll adds all the elements in the specified collection to this collection.
	// Returns true if this collection changed as a result of the call
	AddAll(Collection[T]) bool

	// Clear removes all the elements from this collection.
	Clear()

	// Contains returns true if this collection contains the specified element.
	Contains(T) bool

	// ContainsAll returns true if this collection contains all the elements in the specified collection.
	ContainsAll(Collection[T]) bool

	// DeepCopy implements interface for deep copy of current type.
	DeepCopy() Collection[T]

	// ForEach iterates all elements in this collection readonly with custom callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	ForEach(func(T) bool)

	// IsEmpty returns true if this collection contains no elements.
	IsEmpty() bool

	// Join joins array elements with a string `glue`.
	Join(glue string) string

	// Remove removes all of this collection's elements that are also contained in the specified slice
	// if it is present.
	// Returns true if this collection changed as a result of the call
	Remove(...T) bool

	// RemoveAll removes all of this collection's elements that are also contained in the specified collection
	// Returns true if this collection changed as a result of the call
	RemoveAll(Collection[T]) bool

	// Size returns the number of elements in this collection.
	Size() int

	// Slice returns an array containing shadow copy of all the elements in this collection.
	Slice() []T

	// String returns items as a string, which implements like json.Marshal does.
	String() string
}

type Array[T comparable] interface {
	Collection[T]

	// Chunk splits an array into multiple arrays,
	// the size of each array is determined by `size`.
	// The last chunk may contain less than size elements.
	Chunk(size int) [][]T

	// Clone returns a new array, which is a copy of current array.
	Clone() (newArray Array[T])

	// ContainsI checks whether a value exists in the array with case-insensitively, only applying to element type string
	// For element type other than string, ContainsI returns the same result as Contains
	// Note that it internally iterates the whole array to do the comparison with case-insensitively.
	ContainsI(value T) bool

	// CountValues counts the number of occurrences of all values in the array.
	CountValues() map[T]int

	// ForEachAsc iterates the array readonly in ascending order with given callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	ForEachAsc(f func(int, T) bool)

	// ForEachDesc iterates the array readonly in descending order with given callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	ForEachDesc(f func(int, T) bool)

	// Filter iterates array and filters elements using custom callback function.
	// It removes the element from array if callback function `filter` returns true,
	// it or else does nothing and continues iterating.
	Filter(filter func(index int, value T) bool) Array[T]

	// FilterNil removes all nil value of the array.
	FilterNil() Array[T]

	// FilterEmpty removes all empty value of the array.
	// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
	FilterEmpty() Array[T]

	// Get returns the element at the specified position in this list.
	// If given `index` is out of range, returns empty `value` for type T and bool value false as `found`.
	Get(index int) (value T, found bool)

	// LockFunc locks writing by callback function `f`.
	LockFunc(f func(array Array[T]))

	// MustGet returns the element at the specified position in this list.
	// If given `index` is out of range, returns empty `value` for type T.
	MustGet(index int) (value T)

	// PopLeft pops and returns an item from the beginning of array.
	// Note that if the array is empty, the `found` is false.
	PopLeft() (value T, found bool)

	// PopLefts pops and returns `size` items from the beginning of array.
	PopLefts(size int) []T

	// PopRand randomly pops and return an item out of array.
	// Note that if the array is empty, the `found` is false.
	PopRand() (value T, found bool)

	// PopRands randomly pops and returns `size` items out of array.
	PopRands(size int) []T

	// PopRight pops and returns an item from the end of array.
	// Note that if the array is empty, the `found` is false.
	PopRight() (value T, found bool)

	// PopRights pops and returns `size` items from the end of array.
	PopRights(size int) []T

	// Rand randomly returns one item from array(no deleting).
	Rand() (value T, found bool)

	// Rands randomly returns `size` items from array(no deleting).
	Rands(size int) []T

	// Range picks and returns items by range, like array[start:end].
	// Notice, if in concurrent-safe usage, it returns a copy of slice;
	// else a pointer to the underlying data.
	//
	// If `end` is negative, then the offset will start from the end of array.
	// If `end` is omitted, then the sequence will have everything from start up
	// until the end of the array.
	Range(start int, end ...int) []T

	// RemoveAt removes an item by index.
	// If the given `index` is out of range of the array, the `found` is false.
	RemoveAt(index int) (value T, found bool)

	// RLockFunc locks reading by callback function `f`.
	RLockFunc(f func(array Array[T]))

	// Search searches array by `value`, returns the index of `value`,
	// or returns -1 if not exists.
	Search(value T) int

	// Set replaces the element at the specified position in this list with the specified element.
	Set(index int, value T) error

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
	SubSlice(offset int, length ...int) []T

	// Sum returns the sum of converted integer of each value in an array.
	// Note: converting value into integer may result in unpredictable problems
	Sum() (sum int)

	// Unique uniques the array, clear repeated items.
	// Example: [1,1,2,3,2] -> [1,2,3]
	Unique() Array[T]

	// Walk applies a user supplied function `f` to every item of array.
	Walk(f func(value T) T) Array[T]
}
