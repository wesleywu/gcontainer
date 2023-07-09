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

	// IsEmpty returns true if this collection contains no elements.
	IsEmpty() bool

	// ForEach iterates all elements in this collection readonly with custom callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	ForEach(func(v T) bool)

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
	At(index int) (value T)
	Get(index int) (value T, found bool)
	SetArray(array []T) Array[T]
	Sum() (sum int)
	Remove(index int) (value T, found bool)
	RemoveValue(value T) bool
	RemoveValues(values ...T)
	PopRand() (value T, found bool)
	PopRands(size int) []T
	PopLeft() (value T, found bool)
	PopRight() (value T, found bool)
	PopLefts(size int) []T
	PopRights(size int) []T
	Range(start int, end ...int) []T
	SubSlice(offset int, length ...int) []T
	Append(value ...T) Array[T]
	Len() int
	Slice() []T
	Interfaces() []T
	Clone() (newArray Array[T])
	Clear() Array[T]
	Contains(value T) bool
	ContainsI(value T) bool
	Search(value T) int
	Unique() Array[T]
	LockFunc(f func(array []T)) Array[T]
	RLockFunc(f func(array []T)) Array[T]
	Merge(array Array[T]) Array[T]
	MergeSlice(slice []T) Array[T]
	Chunk(size int) [][]T
	Rand() (value T, found bool)
	Rands(size int) []T
	Join(glue string) string
	CountValues() map[T]int
	Iterator(f func(k int, v T) bool)
	IteratorAsc(f func(k int, v T) bool)
	IteratorDesc(f func(k int, v T) bool)
	String() string
	Filter(filter func(index int, value T) bool) Array[T]
	FilterNil() Array[T]
	FilterEmpty() Array[T]
	Walk(f func(value T) T) Array[T]
	IsEmpty() bool
	DeepCopy() Array[T]
}
