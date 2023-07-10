// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmap provides most commonly used map container which also support concurrent-safe/unsafe switch feature.
package gmap

// Map defines common functions of a `map` and provides more map features.
type Map[K comparable, V comparable] interface {
	// Put sets key-value to the map.
	Put(key K, value V)

	// Puts batch sets key-values to the map.
	Puts(data map[K]V)

	// PutIfAbsent sets `value` to the map if the `key` does not exist, and then returns true.
	// It returns false if `key` exists, and `value` would be ignored.
	PutIfAbsent(key K, value V) bool

	// PutIfAbsentFunc sets value with return value of callback function `f`, and then returns true.
	// It returns false if `key` exists, and `value` would be ignored.
	PutIfAbsentFunc(key K, f func() V) bool

	// Search searches the map with given `key`.
	// Second return parameter `found` is true if key was found, otherwise false.
	Search(key K) (value V, found bool)

	// Get returns the value by given `key`, or empty value of type V if key is not found in the map.
	Get(key K) (value V)

	// GetOrPut returns the value for the given key.
	// If the key is not found in the map, sets its value with given `value` and returns it.
	GetOrPut(key K, value V) V

	// GetOrPutFunc returns the value by key,
	// If the key is not found in the map, calls the f function, puts its result into the map under the given key and returns it.
	GetOrPutFunc(key K, f func() V) V

	// Remove removes the node from the tree by `key`.
	Remove(key K) (value V, removed bool)

	// Removes batch deletes values of the tree by `keys`.
	Removes(keys []K)

	// Iterator iterates all entries in the map readonly with custom callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	Iterator(f func(key K, value V) bool)

	// ContainsKey checks whether `key` exists in the map.
	ContainsKey(key K) bool

	// Size returns the size of the map.
	Size() int

	// Keys returns all keys of the map as a slice, maintaining the order of belonging entries in the map.
	Keys() []K

	// Values returns all values of the map as a slice, maintaining the order of belonging entries in the map.
	Values() []V

	// Map returns a shallow copy of the underlying data of the hash map.
	Map() map[K]V

	// MapStrAny returns a copy of the underlying data of the map as map[string]any.
	MapStrAny() map[string]V

	// IsEmpty checks whether the map is empty.
	// It returns true if map is empty, or else false.
	IsEmpty() bool

	// Clear deletes all data of the map, it will remake a new underlying data map.
	Clear()

	// Replace the data of the map with given `data`.
	Replace(data map[K]V)

	// Clone returns a new hash map with copy of current map data.
	Clone(safe ...bool) Map[K, V]

	// String returns the map as a string.
	String() string
}
