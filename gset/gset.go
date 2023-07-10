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
	garray.Collection[T]

	// Comparator returns the comparator used to order the elements in this set,
	// or nil if this set uses the natural ordering of its elements.
	Comparator() comparator.Comparator[T]

	// First returns the first (lowest) element currently in this set.
	First() (element T, found bool)

	// ForEachDescending iterates the tree readonly in descending order with given callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	ForEachDescending(func(T) bool)

	// HeadSet returns a view of the portion of this set whose elements are less than (or equal to, if inclusive is true) toElement.
	HeadSet(toElement T, inclusive bool) SortedSet[T]

	// Last Returns the last (highest) element currently in this set.
	Last() (element T, found bool)

	// SubSet returns a view of the portion of this set whose elements range from fromElement to toElement.
	SubSet(fromElement T, fromInclusive bool, toElement T, toInclusive bool) SortedSet[T]

	// TailSet returns a view of the portion of this set whose elements are greater than (or equal to, if inclusive is true) fromElement.
	TailSet(fromElement T, inclusive bool) SortedSet[T]
}
