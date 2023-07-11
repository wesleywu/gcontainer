// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gset provides kinds of concurrent-safe/unsafe sets.
package gset

import (
	"github.com/wesleywu/gcontainer/garray"
	"github.com/wesleywu/gcontainer/utils/comparator"
)

type Set[T comparable] interface {
	garray.Collection[T]
}

type SortedSet[T comparable] interface {
	Set[T]

	// Ceiling returns the least element in this set greater than or equal to the given `element` and true as `found`,
	// or empty of type T and false as `found` if there is no such element.
	Ceiling(element T) (ceiling T, found bool)

	// Comparator returns the comparator used to order the elements in this set,
	// or nil if this set uses the natural ordering of its elements.
	Comparator() comparator.Comparator[T]

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

	// PollLast retrieves and removes the last (highest) element and true as `found`,
	// or returns empty of type T and false as `found` if this set is empty.
	PollLast() (last T, found bool)

	// SubSet returns a view of the portion of this set whose elements range from fromElement to toElement.
	SubSet(fromElement T, fromInclusive bool, toElement T, toInclusive bool) SortedSet[T]

	// TailSet returns a view of the portion of this set whose elements are greater than (or equal to, if inclusive is true) fromElement.
	TailSet(fromElement T, inclusive bool) SortedSet[T]
}
