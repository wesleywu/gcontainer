// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
	"github.com/wesleywu/gcontainer/utils/comparators"
)

// Collection is the root interface in the collection hierarchy.
// A collection represents a group of objects, known as its elements.
// Some collections allow duplicate elements and others do not.
// Some are ordered and others unordered.
type Collection[T any] interface {
	// Add adds all the elements in the specified slice to this collection.
	// Returns true if this collection changed as a result of the call
	Add(...T) bool

	// AddAll adds all the elements in the specified collection to this collection.
	// Returns true if this collection changed as a result of the call
	AddAll(Collection[T]) bool

	// Clear removes all the elements from this collection.
	Clear()

	// Clone returns a new collection, which is a copy of current collection.
	Clone() Collection[T]

	// Contains returns true if this collection contains the specified element.
	Contains(T) bool

	// ContainsAll returns true if this collection contains all the elements in the specified collection.
	ContainsAll(Collection[T]) bool

	// DeepCopy implements interface for deep copy of current type.
	DeepCopy() Collection[T]

	// Equals compares the specified object with this collection for equality.
	Equals(another Collection[T]) bool

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

// Set is a collection that contains no duplicate elements. More formally,
// sets contain no pair of elements e1 and e2 such that e1.equals(e2), and at most one nil element.
// As implied by its name, this interface models the mathematical set abstraction.
type Set[T comparable] interface {
	Collection[T]
}

// SortedSet is a Set that further provides a total ordering on its elements.
// The elements are ordered using their natural ordering, or by a Comparator typically provided
// at sorted set creation time.
// The set's iterator will traverse the set in ascending element order.
// Several additional operations are provided to take advantage of the ordering.
// (This interface is the set analogue of SortedMap.)
//
// SortedSet also provides navigation methods lower, floor, ceiling, and higher return elements
// respectively less than, less than or equal, greater than or equal, and greater than a given element,
// returning empty if there is no such element.
type SortedSet[T comparable] interface {
	Set[T]

	// Ceiling returns the least element in this set greater than or equal to the given `element` and true as `found`,
	// or empty of type T and false as `found` if there is no such element.
	Ceiling(element T) (ceiling T, found bool)

	// Comparator returns the comparators used to order the elements in this set,
	// or nil if this set uses the natural ordering of its elements.
	Comparator() comparators.Comparator[T]

	// First returns the first (lowest) element currently in this set.
	First() (element T, found bool)

	// Floor returns the greatest element in this set less than or equal to the given `element` and true as `found`,
	// or empty of type T and false as `found` if there is no such element.
	Floor(element T) (floor T, found bool)

	// ForEachDescending iterates the tree readonly in descending order with given callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	ForEachDescending(func(T) bool)

	// HeadSet returns a view of the portion of this set whose elements are less than (or equal to, if inclusive is true) toElement.
	HeadSet(toElement T, inclusive bool) SortedSet[T]

	// Higher returns the least element in this set strictly greater than the given `element` and true as `found`,
	// or empty of type T and false as `found` if there is no such element.
	Higher(element T) (higher T, found bool)

	// Last Returns the last (highest) element currently in this set.
	Last() (element T, found bool)

	// Lower returns the greatest element in this set strictly less than the given `element` and true as `found`,
	// or empty of type T and false as `found` if there is no such element.
	Lower(element T) (lower T, found bool)

	// PollFirst retrieves and removes the first (lowest) element and true as `found`,
	// or returns empty of type T and false as `found` if this set is empty.
	PollFirst() (first T, found bool)

	// PollHeadSet retrieves and removes portion of this set whose elements are less than
	// (or equal to, if inclusive is true) toElement.
	PollHeadSet(toElement T, inclusive bool) SortedSet[T]

	// PollLast retrieves and removes the last (highest) element and true as `found`,
	// or returns empty of type T and false as `found` if this set is empty.
	PollLast() (last T, found bool)

	// PollTailSet retrieves and removes portion of this set whose elements are greater than
	// (or equal to, if inclusive is true) fromElement.
	PollTailSet(fromElement T, inclusive bool) SortedSet[T]

	// SubSet returns a view of the portion of this set whose elements range from fromElement to toElement.
	SubSet(fromElement T, fromInclusive bool, toElement T, toInclusive bool) SortedSet[T]

	// TailSet returns a view of the portion of this set whose elements are greater than (or equal to, if inclusive is true) fromElement.
	TailSet(fromElement T, inclusive bool) SortedSet[T]
}

// List is nn ordered collection (also known as a sequence). The user of this interface has precise control over
// where in the list each element is inserted. The user can access elements by their integer index (position in the list),
// and search for elements in the list.
type List[T any] interface {
	Collection[T]

	// Chunk splits an array into multiple arrays,
	// the size of each array is determined by `size`.
	// The last chunk may contain less than size elements.
	Chunk(size int) [][]T

	// ContainsI checks whether a value exists in the array with case-insensitively, only applying to element type string
	// For element type other than string, ContainsI returns the same result as Contains
	// Note that it internally iterates the whole array to do the comparison with case-insensitively.
	ContainsI(value T) bool

	// ForEachAsc iterates the array readonly in ascending order with given callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	ForEachAsc(f func(int, T) bool)

	// ForEachDesc iterates the array readonly in descending order with given callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	ForEachDesc(f func(int, T) bool)

	// Filter iterates array and filters elements using custom callback function.
	// It removes the element from array if callback function `filter` returns true,
	// it or else does nothing and continues iterating.
	Filter(filter func(index int, value T) bool) List[T]

	// FilterNil removes all nil value of the array.
	FilterNil() List[T]

	// FilterEmpty removes all empty value of the array.
	// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
	FilterEmpty() List[T]

	// Get returns the element at the specified position in this list.
	// If given `index` is out of range, returns empty `value` for type T and bool value false as `found`.
	Get(index int) (value T, found bool)

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

	// Search searches array by `value`, returns the index of `value`,
	// or returns -1 if not exists.
	Search(value T) int

	// Set replaces the element at the specified position in this list with the specified element.
	Set(index int, value T) error

	// Sort sorts the array by custom function `less`.
	Sort(less func(v1, v2 T) bool)

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
	Unique() List[T]

	// Walk applies a user supplied function `f` to every item of array.
	Walk(f func(value T) T) List[T]
}

// Map defines common functions of a `map` and provides more map features.
// The Map interface provides three collection views, which allow a map's contents to be viewed as a set of keys,
// collection of values, or set of key-value mappings.
// The order of a map is defined as the order in which the iterators on the map's collection views return their elements.
// Some map implementations, like the TreeMap struct, make specific guarantees as to their order;
// others, like the HashMap struct, do not.
type Map[K comparable, V any] interface {
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

	// ForEach iterates all entries in the map readonly with custom callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	ForEach(f func(key K, value V) bool)

	// ContainsKey checks whether `key` exists in the map.
	ContainsKey(key K) bool

	// Size returns the size of the map.
	Size() int

	// KeySet returns all keys of the map as a set.
	KeySet() Set[K]

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

	// Compute attempts to compute a mapping for the specified `key` and its current mapped value
	// (or empty if there is no current mapping).
	// For example, to either create or append a String msg to a value mapping:
	//
	// If the function `f` returns nil, the mapping is removed (or remains absent if initially absent).
	// If the function itself returns an error, the error is rethrown, and the current mapping is left unchanged.
	// todo implements me
	//Compute(key K, f func(key K, value V) (V, error)) error

	// String returns the map as a string.
	String() string
}

// SortedMap is a Map that further provides a total ordering on its keys. The map is ordered according to
// the natural ordering of its keys, or by a Comparator typically provided at sorted map creation time.
// This order is reflected when iterating over the sorted map's collection views (returned by the entrySet,
// keySet and values methods). Several additional operations are provided to take advantage of the ordering.
//
// (This interface is the map analogue of SortedSet.)
//
// All keys inserted into a sorted map must implement the Comparable interface (or be accepted by
// the specified comparator).
//
// SortedMap also provides navigation methods returning the closest matches for given search targets.
// Methods LowerEntry, FloorEntry, CeilingEntry, and HigherEntry return MapEntry associated with keys
// respectively less than, less than or equal, greater than or equal, and greater than a given key,
// returning empty if there is no such key.
// Similarly, methods LowerKey, FloorKey, CeilingKey, and HigherKey return only the associated keys.
// All of these methods are designed for locating, not traversing entries.
type SortedMap[K comparable, V any] interface {
	Map[K, V]

	// AscendingKeySet returns a view of the keys contained in this map, in its natural ascending order.
	AscendingKeySet() SortedSet[K]

	// CeilingEntry returns a key-value mapping associated with the least key greater than or equal to the given key, or nil if there is no such key.
	CeilingEntry(key K) MapEntry[K, V]

	// CeilingKey returns the least key greater than or equal to the given key, or empty of type K if there is no such key.
	// The parameter `ok` indicates whether a non-empty `ceilingKey` is returned.
	CeilingKey(key K) (ceilingKey K, ok bool)

	// DescendingKeySet returns a reversed order view of the keys contained in this map.
	DescendingKeySet() SortedSet[K]

	// FirstEntry returns a key-value mapping associated with the least key in this map, or nil if the map is empty.
	FirstEntry() MapEntry[K, V]

	// FloorEntry returns a key-value mapping associated with the greatest key less than or equal to the given key, or nil if there is no such key.
	FloorEntry(key K) MapEntry[K, V]

	// FloorKey returns the greatest key less than or equal to the given key, or empty of type K if there is no such key.
	// The parameter `ok` indicates whether a non-empty `floorKey` is returned.
	FloorKey(key K) (floorKey K, ok bool)

	// HeadMap returns a view of the portion of this map whose keys are less than (or equal to, if inclusive is true) toKey.
	HeadMap(toKey K, inclusive bool) SortedMap[K, V]

	// HigherEntry returns a key-value mapping associated with the least key strictly greater than the given key, or nil if there is no such key.
	HigherEntry(key K) MapEntry[K, V]

	// HigherKey returns the least key strictly greater than the given key, or empty of type K if there is no such key.
	// The parameter `ok` indicates whether a non-empty `higherKey` is returned.
	HigherKey(key K) (higherKey K, ok bool)

	// LastEntry returns a key-value mapping associated with the greatest key in this map, or nil if the map is empty.
	LastEntry() MapEntry[K, V]

	// LowerEntry returns a key-value mapping associated with the greatest key strictly less than the given key, or nil if there is no such key.
	LowerEntry(key K) MapEntry[K, V]

	// LowerKey returns the greatest key strictly less than the given key, or empty of type K if there is no such key.
	// The parameter `ok` indicates whether a non-empty `lowerKey` is returned.
	LowerKey(key K) (lowerKey K, ok bool)

	// PollFirstEntry removes and returns a key-value mapping associated with the least key in this map, or nil if the map is empty.
	PollFirstEntry() MapEntry[K, V]

	// PollLastEntry removes and returns a key-value mapping associated with the greatest key in this map, or nil if the map is empty.
	PollLastEntry() MapEntry[K, V]

	// Reverse returns a reverse order view of the mappings contained in this map.
	Reverse() SortedMap[K, V]

	// SubMap returns a view of the portion of this map whose keys range from fromKey to toKey.
	SubMap(fromKey K, fromInclusive bool, toKey K, toInclusive bool) SortedMap[K, V]

	// TailMap returns a view of the portion of this map whose keys are greater than (or equal to, if inclusive is true) fromKey.
	TailMap(fromKey K, inclusive bool) SortedMap[K, V]
}

// MapEntry is a key-value pair, usually representing the element entries in a Map
type MapEntry[K comparable, V any] interface {
	// Key returns the key corresponding to this entry.
	Key() K

	// Value returns the value corresponding to this entry.
	Value() V
}
