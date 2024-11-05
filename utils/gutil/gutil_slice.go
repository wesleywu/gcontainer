// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

// SliceCopy does a shallow copy of slice `data` for most commonly used slice type
// []T.
func SliceCopy[T any](slice []T) []T {
	newSlice := make([]T, len(slice))
	copy(newSlice, slice)
	return newSlice
}

// SliceInsertBefore inserts the `values` to the front of `index` and returns a new slice.
func SliceInsertBefore[T any](slice []T, index int, values ...T) (newSlice []T) {
	if index < 0 || index >= len(slice) {
		return slice
	}
	newSlice = make([]T, len(slice)+len(values))
	copy(newSlice, slice[0:index])
	copy(newSlice[index:], values)
	copy(newSlice[index+len(values):], slice[index:])
	return
}

// SliceInsertAfter inserts the `values` to the back of `index` and returns a new slice.
func SliceInsertAfter[T any](slice []T, index int, values ...T) (newSlice []T) {
	if index < 0 || index >= len(slice) {
		return slice
	}
	newSlice = make([]T, len(slice)+len(values))
	copy(newSlice, slice[0:index+1])
	copy(newSlice[index+1:], values)
	copy(newSlice[index+1+len(values):], slice[index+1:])
	return
}

// SliceDelete deletes an element at `index` and returns the new slice.
// It does nothing if the given `index` is invalid.
func SliceDelete[T any](slice []T, index int) (newSlice []T) {
	if index < 0 || index >= len(slice) {
		return slice
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		return slice[1:]
	} else if index == len(slice)-1 {
		return slice[:index]
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	return append(slice[:index], slice[index+1:]...)
}
