// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
	"bytes"
	json2 "encoding/json"
	"fmt"

	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/internal/rwmutex"
	"github.com/wesleywu/gcontainer/utils/comparators"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

type color bool

const (
	black, red color = true, false
)

// TreeMap implements the red-black tree.
type TreeMap[K comparable, V any] struct {
	mu         rwmutex.RWMutex
	root       *RedBlackTreeNode[K, V]
	size       int
	comparator comparators.Comparator[K]
}

// RedBlackTreeNode is a single element within the tree.
type RedBlackTreeNode[K comparable, V any] struct {
	key    K
	value  V
	color  color
	left   *RedBlackTreeNode[K, V]
	right  *RedBlackTreeNode[K, V]
	parent *RedBlackTreeNode[K, V]
}

// NewTreeMap instantiates a red-black tree with the custom key comparators.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewTreeMap[K comparable, V any](comparator comparators.Comparator[K], safe ...bool) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
	}
}

// NewTreeMapDefault instantiates a red-black tree with default key comparators.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewTreeMapDefault[K comparable, V any](safe ...bool) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		mu:         rwmutex.Create(safe...),
		comparator: comparators.ComparatorAny[K],
	}
}

// NewTreeMapFrom instantiates a red-black tree with the custom key comparators and `data` map.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewTreeMapFrom[K comparable, V any](comparator func(v1, v2 K) int, data map[K]V, safe ...bool) *TreeMap[K, V] {
	tree := NewTreeMap[K, V](comparator, safe...)
	for k, v := range data {
		tree.insertEntry(k, v)
	}
	return tree
}

func (n *RedBlackTreeNode[K, V]) Key() K {
	return n.key
}

func (n *RedBlackTreeNode[K, V]) Value() V {
	return n.value
}

// AscendingKeySet returns a ascending order view of the keys contained in this map.
func (tree *TreeMap[K, V]) AscendingKeySet() SortedSet[K] {
	var (
		keySet = NewTreeSet[K](tree.Comparator(), tree.mu.IsSafe())
		index  = 0
	)
	tree.ForEachAsc(func(key K, value V) bool {
		keySet.Add(key)
		index++
		return true
	})
	return keySet
}

func (tree *TreeMap[K, V]) Comparator() comparators.Comparator[K] {
	if tree.comparator == nil {
		tree.comparator = comparators.ComparatorAny[K]
	}
	return tree.comparator
}

// SetComparator sets/changes the comparators for sorting.
func (tree *TreeMap[K, V]) SetComparator(comparator comparators.Comparator[K]) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.comparator = comparator
	if tree.size > 0 {
		data := make(map[K]V, tree.size)
		tree.doIteratorAsc(tree.leftNode(), func(key K, value V) bool {
			data[key] = value
			return true
		})
		// Resort the tree if comparators is changed.
		tree.root = nil
		tree.size = 0
		for k, v := range data {
			tree.insertEntry(k, v)
		}
	}
}

// Clone returns a new tree with a copy of current tree.
func (tree *TreeMap[K, V]) Clone(safe ...bool) Map[K, V] {
	newTree := NewTreeMap[K, V](tree.comparator, safe...)
	newTree.Puts(tree.Map())
	return newTree
}

// DescendingKeySet returns a reversed order view of the keys contained in this map.
func (tree *TreeMap[K, V]) DescendingKeySet() SortedSet[K] {
	var (
		keySet = NewTreeSet[K](comparators.Reverse(tree.Comparator()), tree.mu.IsSafe())
		index  = 0
	)
	tree.ForEachDesc(func(key K, value V) bool {
		keySet.Add(key)
		index++
		return true
	})
	return keySet
}

func (tree *TreeMap[K, V]) FirstEntry() MapEntry[K, V] {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.leftNode()
}

// Put inserts key-value item into the tree.
func (tree *TreeMap[K, V]) Put(key K, value V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.insertEntry(key, value)
}

// Puts batch sets key-values to the tree.
func (tree *TreeMap[K, V]) Puts(data map[K]V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for k, v := range data {
		tree.insertEntry(k, v)
	}
}

func (tree *TreeMap[K, V]) insertEntry(key K, value V) (putValue V) {
	t := tree.root
	if t == nil {
		tree.Comparator()(key, key) // type (and possibly nil) check

		tree.root = &RedBlackTreeNode[K, V]{key: key, value: value, color: black}
		tree.size = 1
		//modCount++
		return
	}
	var cmp int
	var parent *RedBlackTreeNode[K, V]
	// split comparator and comparable paths
	var cpr = tree.Comparator()
	for {
		parent = t
		cmp = cpr(key, t.key)
		if cmp < 0 {
			t = t.left
		} else if cmp > 0 {
			t = t.right
		} else {
			t.value = value
			return value
		}
		if t == nil {
			break
		}
	}
	e := &RedBlackTreeNode[K, V]{key: key, value: value, parent: parent, color: black}
	if cmp < 0 {
		parent.left = e
	} else {
		parent.right = e
	}
	tree.fixAfterInsertion(e)
	tree.size++
	//modCount++
	return
}

// Get returns the value by given `key`, or empty value of type K if the key is not found in the map.
func (tree *TreeMap[K, V]) Get(key K) (value V) {
	value, _ = tree.Search(key)
	return
}

// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given `key`,
// or else just return the existing value.
//
// When setting value, if `value` is type of <func() interface {}>,
// it will be executed with mutex.Lock of the hash map,
// and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (tree *TreeMap[K, V]) doSetWithLockCheck(key K, value V) V {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if node := tree.getEntry(key); node != nil {
		return node.value
	}
	if any(value) != nil {
		tree.insertEntry(key, value)
	}
	return value
}

// doSetWithLockCheckFunc checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given `key`,
// or else just return the existing value.
//
// When setting value, if `value` is type of <func() interface {}>,
// it will be executed with mutex.Lock of the hash map,
// and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (tree *TreeMap[K, V]) doSetWithLockCheckFunc(key K, f func() V) V {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if node := tree.getEntry(key); node != nil {
		return node.value
	}
	value := f()
	if any(value) != nil {
		tree.insertEntry(key, value)
	}
	return value
}

// GetOrPut returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (tree *TreeMap[K, V]) GetOrPut(key K, value V) V {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrPutFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
// with mutex.Lock of the hash map.
func (tree *TreeMap[K, V]) GetOrPutFunc(key K, f func() V) V {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheckFunc(key, f)
	} else {
		return v
	}
}

// PutIfAbsent sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (tree *TreeMap[K, V]) PutIfAbsent(key K, value V) bool {
	if !tree.ContainsKey(key) {
		tree.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// PutIfAbsentFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (tree *TreeMap[K, V]) PutIfAbsentFunc(key K, f func() V) bool {
	if !tree.ContainsKey(key) {
		tree.doSetWithLockCheckFunc(key, f)
		return true
	}
	return false
}

// ContainsKey checks whether `key` exists in the tree.
func (tree *TreeMap[K, V]) ContainsKey(key K) bool {
	_, ok := tree.Search(key)
	return ok
}

func (tree *TreeMap[K, V]) deleteEntry(p *RedBlackTreeNode[K, V]) {
	tree.size--

	// If strictly internal, copy successor's element to p and then make p
	// point to successor.
	if p.left != nil && p.right != nil {
		s := successor(p)
		p.key = s.key
		p.value = s.value
		p = s
	} // p has 2 children

	// Start fixup at replacement node, if it exists.
	var replacement *RedBlackTreeNode[K, V]
	if p.left != nil {
		replacement = p.left
	} else {
		replacement = p.right
	}

	if replacement != nil {
		// Link replacement to parent
		replacement.parent = p.parent
		if p.parent == nil {
			tree.root = replacement
		} else if p == p.parent.left {
			p.parent.left = replacement
		} else {
			p.parent.right = replacement
		}

		// Nil out links so they are OK to use by fixAfterDeletion.
		p.left = nil
		p.right = nil
		p.parent = nil

		// Fix replacement
		if p.color == black {
			tree.fixAfterDeletion(replacement)
		}
	} else if p.parent == nil { // return if we are the only node.
		tree.root = nil
	} else { //  No children. Use self as phantom replacement and unlink.
		if p.color == black {
			tree.fixAfterDeletion(p)
		}
		if p.parent != nil {
			if p == p.parent.left {
				p.parent.left = nil
			} else if p == p.parent.right {
				p.parent.right = nil
			}
			p.parent = nil
		}
	}
}

func (tree *TreeMap[K, V]) fixAfterDeletion(x *RedBlackTreeNode[K, V]) {
	for x != tree.root && x.color == black {
		if x == leftOf(parentOf(x)) {
			sib := rightOf(parentOf(x))

			if colorOf(sib) == red {
				setColor(sib, black)
				setColor(parentOf(x), red)
				tree.rotateLeft(parentOf(x))
				sib = rightOf(parentOf(x))
			}

			if colorOf(leftOf(sib)) == black &&
				colorOf(rightOf(sib)) == black {
				setColor(sib, red)
				x = parentOf(x)
			} else {
				if colorOf(rightOf(sib)) == black {
					setColor(leftOf(sib), black)
					setColor(sib, red)
					tree.rotateRight(sib)
					sib = rightOf(parentOf(x))
				}
				setColor(sib, colorOf(parentOf(x)))
				setColor(parentOf(x), black)
				setColor(rightOf(sib), black)
				tree.rotateLeft(parentOf(x))
				x = tree.root
			}
		} else { // symmetric
			sib := leftOf(parentOf(x))

			if colorOf(sib) == red {
				setColor(sib, black)
				setColor(parentOf(x), red)
				tree.rotateRight(parentOf(x))
				sib = leftOf(parentOf(x))
			}

			if colorOf(rightOf(sib)) == black &&
				colorOf(leftOf(sib)) == black {
				setColor(sib, red)
				x = parentOf(x)
			} else {
				if colorOf(leftOf(sib)) == black {
					setColor(rightOf(sib), black)
					setColor(sib, red)
					tree.rotateLeft(sib)
					sib = leftOf(parentOf(x))
				}
				setColor(sib, colorOf(parentOf(x)))
				setColor(parentOf(x), black)
				setColor(leftOf(sib), black)
				tree.rotateRight(parentOf(x))
				x = tree.root
			}
		}
	}
	x.color = black
}

func leftOf[K comparable, V any](p *RedBlackTreeNode[K, V]) *RedBlackTreeNode[K, V] {
	if p == nil {
		return nil
	}
	return p.left
}

func rightOf[K comparable, V any](p *RedBlackTreeNode[K, V]) *RedBlackTreeNode[K, V] {
	if p == nil {
		return nil
	}
	return p.right
}

func parentOf[K comparable, V any](p *RedBlackTreeNode[K, V]) *RedBlackTreeNode[K, V] {
	if p == nil {
		return nil
	}
	return p.parent
}

func colorOf[K comparable, V any](p *RedBlackTreeNode[K, V]) color {
	if p == nil {
		return black
	}
	return p.color
}

func setColor[K comparable, V any](p *RedBlackTreeNode[K, V], c color) {
	if p != nil {
		p.color = c
	}
}

// successor returns the successor of the specified Entry, or nil if no such.
func successor[K comparable, V any](t *RedBlackTreeNode[K, V]) *RedBlackTreeNode[K, V] {
	if t == nil {
		return nil
	} else if t.right != nil {
		p := t.right
		for p.left != nil {
			p = p.left
		}
		return p
	} else {
		p := t.parent
		ch := t
		for p != nil && ch == p.right {
			ch = p
			p = p.parent
		}
		return p
	}
}

// predecessor returns the predecessor of the specified Entry, or nil if no such.
func predecessor[K comparable, V any](t *RedBlackTreeNode[K, V]) *RedBlackTreeNode[K, V] {
	if t == nil {
		return nil
	} else if t.left != nil {
		p := t.left
		for p.right != nil {
			p = p.right
		}
		return p
	} else {
		p := t.parent
		ch := t
		for p != nil && ch == p.left {
			ch = p
			p = p.parent
		}
		return p
	}
}

func (tree *TreeMap[K, V]) PollFirstEntry() MapEntry[K, V] {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	node := tree.leftNode()
	if node == nil {
		return nil
	}
	tree.deleteEntry(node)
	if tree.mu.IsSafe() {
		return &RedBlackTreeNode[K, V]{
			key:   node.key,
			value: node.value,
		}
	}
	return node
}

func (tree *TreeMap[K, V]) PollLastEntry() MapEntry[K, V] {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	node := tree.rightNode()
	if node == nil {
		return nil
	}
	tree.deleteEntry(node)
	if tree.mu.IsSafe() {
		return &RedBlackTreeNode[K, V]{
			key:   node.key,
			value: node.value,
		}
	}
	return node
}

// Remove removes the node from the tree by `key`.
func (tree *TreeMap[K, V]) Remove(key K) (value V, removed bool) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	node := tree.getEntry(key)
	if node == nil {
		return
	}
	value = node.value
	tree.deleteEntry(node)
	return value, true
}

// Removes batch deletes values of the tree by `keys`.
func (tree *TreeMap[K, V]) Removes(keys []K) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		node := tree.getEntry(key)
		if node == nil {
			continue
		}
		tree.deleteEntry(node)
	}
}

// Reverse returns a reverse order view of the mappings contained in this map.
func (tree *TreeMap[K, V]) Reverse() SortedMap[K, V] {
	newTree := NewTreeMap[K, V](comparators.Reverse(tree.comparator), tree.mu.IsSafe())
	newTree.Puts(tree.Map())
	return newTree
}

// IsEmpty returns true if tree does not contain any nodes.
func (tree *TreeMap[K, V]) IsEmpty() bool {
	return tree.Size() == 0
}

// Size returns number of nodes in the tree.
func (tree *TreeMap[K, V]) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.size
}

// Keys returns all keys in asc order.
func (tree *TreeMap[K, V]) Keys() []K {
	var (
		keys  = make([]K, tree.Size())
		index = 0
	)
	tree.ForEachAsc(func(key K, value V) bool {
		keys[index] = key
		index++
		return true
	})
	return keys
}

// Values returns all values in asc order based on the key.
func (tree *TreeMap[K, V]) Values() []V {
	var (
		values = make([]V, tree.Size())
		index  = 0
	)
	tree.ForEachAsc(func(key K, value V) bool {
		values[index] = value
		index++
		return true
	})
	return values
}

// Map returns all key-value items as map.
func (tree *TreeMap[K, V]) Map() map[K]V {
	m := make(map[K]V, tree.Size())
	tree.ForEachAsc(func(key K, value V) bool {
		m[key] = value
		return true
	})
	return m
}

// MapStrAny returns all key-value items as map[string]V.
func (tree *TreeMap[K, V]) MapStrAny() map[string]V {
	m := make(map[string]V, tree.Size())
	tree.ForEachAsc(func(key K, value V) bool {
		m[gconv.String(key)] = value
		return true
	})
	return m
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *TreeMap[K, V]) Left() *RedBlackTreeNode[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.leftNode()
	if tree.mu.IsSafe() {
		return &RedBlackTreeNode[K, V]{
			key:   node.key,
			value: node.value,
		}
	}
	return node
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *TreeMap[K, V]) Right() *RedBlackTreeNode[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.rightNode()
	if tree.mu.IsSafe() {
		return &RedBlackTreeNode[K, V]{
			key:   node.key,
			value: node.value,
		}
	}
	return node
}

// leftNode returns the left-most (min) node or nil if tree is empty.
func (tree *TreeMap[K, V]) leftNode() *RedBlackTreeNode[K, V] {
	p := (*RedBlackTreeNode[K, V])(nil)
	n := tree.root
	for n != nil {
		p = n
		n = n.left
	}
	return p
}

// rightNode returns the right-most (max) node or nil if tree is empty.
func (tree *TreeMap[K, V]) rightNode() *RedBlackTreeNode[K, V] {
	p := (*RedBlackTreeNode[K, V])(nil)
	n := tree.root
	for n != nil {
		p = n
		n = n.right
	}
	return p
}

// FloorEntry returns the tree node associated with the greatest key less than or equal to the given key, or nil if there is no such key.
// Second return parameter is true if FloorEntry was found, otherwise false.
//
// A FloorEntry node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
func (tree *TreeMap[K, V]) FloorEntry(key K) MapEntry[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	p := tree.root
	for p != nil {
		cmp := tree.getComparator()(key, p.key)
		if cmp > 0 {
			if p.right != nil {
				p = p.right
			} else {
				return p
			}
		} else if cmp < 0 {
			if p.left != nil {
				p = p.left
			} else {
				parent := p.parent
				ch := p
				for parent != nil && ch == parent.left {
					ch = parent
					parent = parent.parent
				}
				if parent == nil {
					return nil
				}
				return parent
			}
		} else {
			return p
		}
	}
	return nil
}

// FloorKey returns the greatest key less than or equal to the given key, or empty of type K if there is no such key.
// The parameter `ok` indicates whether a non-empty `floorKey` is returned.
func (tree *TreeMap[K, V]) FloorKey(key K) (floorKey K, ok bool) {
	if entry := tree.FloorEntry(key); entry != nil {
		return entry.Key(), true
	}
	return
}

// CeilingEntry finds ceiling node of the input key, return the ceiling node or nil if no ceiling node is found.
// Second return parameter is true if ceiling was found, otherwise false.
//
// CeilingEntry node is defined as the smallest node that its key is larger than or equal to the given `key`.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
func (tree *TreeMap[K, V]) CeilingEntry(key K) MapEntry[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	p := tree.root
	for p != nil {
		cmp := tree.getComparator()(key, p.key)
		if cmp < 0 {
			if p.left != nil {
				p = p.left
			} else {
				return p
			}
		} else if cmp > 0 {
			if p.right != nil {
				p = p.right
			} else {
				parent := p.parent
				ch := p
				for parent != nil && ch == parent.right {
					ch = parent
					parent = parent.parent
				}
				if parent == nil {
					return nil
				}
				return parent
			}
		} else {
			return p
		}
	}
	return nil
}

// CeilingKey returns the least key greater than or equal to the given key, or empty of type K if there is no such key.
// The parameter `ok` indicates whether a non-empty `ceilingKey` is returned.
func (tree *TreeMap[K, V]) CeilingKey(key K) (ceilingKey K, ok bool) {
	if entry := tree.CeilingEntry(key); entry != nil {
		return entry.Key(), true
	}
	return
}

// HeadMap returns a view of the portion of this map whose keys are less than (or equal to, if inclusive is true) toKey.
func (tree *TreeMap[K, V]) HeadMap(toKey K, inclusive bool) SortedMap[K, V] {
	result := NewTreeMap[K, V](tree.Comparator(), tree.mu.IsSafe())
	tree.IteratorDescFrom(toKey, inclusive, func(key K, value V) bool {
		result.Put(key, value)
		return true
	})
	return result
}

// LowerEntry returns the tree node associated with the greatest key strictly less than the given key, or nil if there is no such key.
//
// A LowerEntry node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than or equal to the given node.
func (tree *TreeMap[K, V]) LowerEntry(key K) MapEntry[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	p := tree.root
	for p != nil {
		cmp := tree.getComparator()(key, p.key)
		if cmp > 0 {
			if p.right != nil {
				p = p.right
			} else {
				return p
			}
		} else {
			if p.left != nil {
				p = p.left
			} else {
				parent := p.parent
				ch := p
				for parent != nil && ch == parent.left {
					ch = parent
					parent = parent.parent
				}
				if parent == nil {
					return nil
				}
				return parent
			}
		}
	}
	return nil
}

// LowerKey returns the greatest key strictly less than the given key, or empty of type K if there is no such key.
// The parameter `ok` indicates whether a non-empty `lowerKey` is returned.
func (tree *TreeMap[K, V]) LowerKey(key K) (lowerKey K, ok bool) {
	if entry := tree.LowerEntry(key); entry != nil {
		return entry.Key(), true
	}
	return
}

// HigherEntry returns the tree node associated with the least key strictly greater than the given key, or nil if there is no such key.
// Second return parameter is true if HigherEntry was found, otherwise false.
//
// A HigherEntry node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than or equal to the given node.
func (tree *TreeMap[K, V]) HigherEntry(key K) MapEntry[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	p := tree.root
	for p != nil {
		cmp := tree.getComparator()(key, p.key)
		if cmp < 0 {
			if p.left != nil {
				p = p.left
			} else {
				return p
			}
		} else {
			if p.right != nil {
				p = p.right
			} else {
				parent := p.parent
				ch := p
				for parent != nil && ch == parent.right {
					ch = parent
					parent = parent.parent
				}
				if parent == nil {
					return nil
				}
				return parent
			}
		}
	}
	return nil
}

// HigherKey returns the least key strictly greater than the given key, or empty of type K if there is no such key.
// The parameter `ok` indicates whether a non-empty `higherKey` is returned.
func (tree *TreeMap[K, V]) HigherKey(key K) (higherKey K, ok bool) {
	if entry := tree.HigherEntry(key); entry != nil {
		return entry.Key(), true
	}
	return
}

// ForEach is alias of ForEachAsc.
func (tree *TreeMap[K, V]) ForEach(f func(key K, value V) bool) {
	tree.ForEachAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
func (tree *TreeMap[K, V]) IteratorFrom(key K, inclusive bool, f func(key K, value V) bool) {
	tree.IteratorAscFrom(key, inclusive, f)
}

// ForEachAsc iterates the tree readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *TreeMap[K, V]) ForEachAsc(f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorAsc(tree.leftNode(), f)
}

// IteratorAscFrom iterates the tree readonly in ascending order with given callback function `f`.
// The parameter `key` specifies the start entry for iterating. The `match` specifies whether
// starting iterating if the `key` is fully matched, or else using index searching iterating.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *TreeMap[K, V]) IteratorAscFrom(key K, inclusive bool, f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var entry MapEntry[K, V]
	if inclusive {
		entry = tree.CeilingEntry(key)
	} else {
		entry = tree.HigherEntry(key)
	}
	if entry == nil {
		return
	}
	tree.doIteratorAsc(entry.(*RedBlackTreeNode[K, V]), f)
}

func (tree *TreeMap[K, V]) doIteratorAsc(node *RedBlackTreeNode[K, V], f func(key K, value V) bool) {
loop:
	if node == nil {
		return
	}
	if !f(node.key, node.value) {
		return
	}
	if node.right != nil {
		node = node.right
		for node.left != nil {
			node = node.left
		}
		goto loop
	}
	if node.parent != nil {
		old := node
		for node.parent != nil {
			node = node.parent
			if tree.getComparator()(old.key, node.key) <= 0 {
				goto loop
			}
		}
	}
}

// ForEachDesc iterates the tree readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *TreeMap[K, V]) ForEachDesc(f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorDesc(tree.rightNode(), f)
}

// IteratorDescFrom iterates the tree readonly in descending order with given callback function `f`.
// The parameter `key` specifies the start entry for iterating. The `match` specifies whether
// starting iterating if the `key` is fully matched, or else using index searching iterating.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *TreeMap[K, V]) IteratorDescFrom(key K, inclusive bool, f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var entry MapEntry[K, V]
	if inclusive {
		entry = tree.FloorEntry(key)
	} else {
		entry = tree.LowerEntry(key)
	}
	if entry == nil {
		return
	}
	tree.doIteratorDesc(entry.(*RedBlackTreeNode[K, V]), f)
}

func (tree *TreeMap[K, V]) doIteratorDesc(node *RedBlackTreeNode[K, V], f func(key K, value V) bool) {
loop:
	if node == nil {
		return
	}
	if !f(node.key, node.value) {
		return
	}
	if node.left != nil {
		node = node.left
		for node.right != nil {
			node = node.right
		}
		goto loop
	}
	if node.parent != nil {
		old := node
		for node.parent != nil {
			node = node.parent
			if tree.getComparator()(old.key, node.key) >= 0 {
				goto loop
			}
		}
	}
}

func (tree *TreeMap[K, V]) LastEntry() MapEntry[K, V] {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.rightNode()
}

// SubMap returns a view of the portion of this map whose keys range from fromKey to toKey.
func (tree *TreeMap[K, V]) SubMap(fromKey K, fromInclusive bool, toKey K, toInclusive bool) SortedMap[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var (
		startElement *RedBlackTreeNode[K, V]
		endElement   *RedBlackTreeNode[K, V]
		outOfBound   bool
		result       = NewTreeMap[K, V](tree.getComparator(), tree.mu.IsSafe())
	)
	if fromInclusive {
		entry := tree.CeilingEntry(fromKey)
		if entry != nil {
			startElement = entry.(*RedBlackTreeNode[K, V])
		} else {
			outOfBound = true
		}
	} else {
		entry := tree.HigherEntry(fromKey)
		if entry != nil {
			startElement = entry.(*RedBlackTreeNode[K, V])
		} else {
			outOfBound = true
		}
	}
	if outOfBound {
		return result
	}
	if toInclusive {
		entry := tree.FloorEntry(toKey)
		if entry != nil {
			endElement = entry.(*RedBlackTreeNode[K, V])
		} else {
			outOfBound = true
		}
	} else {
		entry := tree.LowerEntry(toKey)
		if entry != nil {
			endElement = entry.(*RedBlackTreeNode[K, V])
		} else {
			outOfBound = true
		}
	}
	if outOfBound {
		return result
	}
	tree.doIteratorAsc(startElement, func(key K, value V) bool {
		if tree.getComparator()(key, endElement.key) > 0 {
			return false
		}
		result.Put(key, value)
		return true
	})
	return result
}

// Clear removes all nodes from the tree.
func (tree *TreeMap[K, V]) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
}

// Replace the data of the tree with given `data`.
func (tree *TreeMap[K, V]) Replace(data map[K]V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
	for k, v := range data {
		tree.insertEntry(k, v)
	}
}

// String returns a string representation of container.
func (tree *TreeMap[K, V]) String() string {
	if tree == nil {
		return ""
	}
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	str := ""
	if tree.size != 0 {
		tree.output(tree.root, "", true, &str)
	}
	return str
}

// Print prints the tree to stdout.
func (tree *TreeMap[K, V]) Print() {
	fmt.Println(tree.String())
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *TreeMap[K, V]) Search(key K) (value V, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.getEntry(key)
	if node != nil {
		return node.value, true
	}
	return
}

// TailMap returns a view of the portion of this map whose keys are greater than (or equal to, if inclusive is true) fromKey.
func (tree *TreeMap[K, V]) TailMap(fromKey K, inclusive bool) SortedMap[K, V] {
	result := NewTreeMap[K, V](tree.Comparator(), tree.mu.IsSafe())
	tree.IteratorAscFrom(fromKey, inclusive, func(key K, value V) bool {
		result.Put(key, value)
		return true
	})
	return result
}

func (tree *TreeMap[K, V]) output(node *RedBlackTreeNode[K, V], prefix string, isTail bool, str *string) {
	if node.right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		tree.output(node.right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += fmt.Sprintf("%v\n", node.key)
	if node.left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		tree.output(node.left, newPrefix, true, str)
	}
}

func (tree *TreeMap[K, V]) getEntry(key K) *RedBlackTreeNode[K, V] {
	p := tree.root
	for p != nil {
		compare := tree.getComparator()(key, p.key)
		switch {
		case compare == 0:
			return p
		case compare < 0:
			p = p.left
		case compare > 0:
			p = p.right
		}
	}
	return p
}

func (tree *TreeMap[K, V]) rotateLeft(p *RedBlackTreeNode[K, V]) {
	if p == nil {
		return
	}
	r := p.right
	p.right = r.left
	if r.left != nil {
		r.left.parent = p
	}
	r.parent = p.parent
	if p.parent == nil {
		tree.root = r
	} else if p.parent.left == p {
		p.parent.left = r
	} else {
		p.parent.right = r
	}
	r.left = p
	p.parent = r
}

func (tree *TreeMap[K, V]) rotateRight(p *RedBlackTreeNode[K, V]) {
	if p == nil {
		return
	}
	l := p.left
	p.left = l.right
	if l.right != nil {
		l.right.parent = p
	}
	l.parent = p.parent
	if p.parent == nil {
		tree.root = l
	} else if p.parent.right == p {
		p.parent.right = l
	} else {
		p.parent.left = l
	}
	l.right = p
	p.parent = l
}

func (tree *TreeMap[K, V]) fixAfterInsertion(x *RedBlackTreeNode[K, V]) {
	x.color = red

	for x != nil && x != tree.root && x.parent.color == red {
		if parentOf(x) == leftOf(parentOf(parentOf(x))) {
			y := rightOf(parentOf(parentOf(x)))
			if colorOf(y) == red {
				setColor(parentOf(x), black)
				setColor(y, black)
				setColor(parentOf(parentOf(x)), red)
				x = parentOf(parentOf(x))
			} else {
				if x == rightOf(parentOf(x)) {
					x = parentOf(x)
					tree.rotateLeft(x)
				}
				setColor(parentOf(x), black)
				setColor(parentOf(parentOf(x)), red)
				tree.rotateRight(parentOf(parentOf(x)))
			}
		} else {
			y := leftOf(parentOf(parentOf(x)))
			if colorOf(y) == red {
				setColor(parentOf(x), black)
				setColor(y, black)
				setColor(parentOf(parentOf(x)), red)
				x = parentOf(parentOf(x))
			} else {
				if x == leftOf(parentOf(x)) {
					x = parentOf(x)
					tree.rotateRight(x)
				}
				setColor(parentOf(x), black)
				setColor(parentOf(parentOf(x)), red)
				tree.rotateLeft(parentOf(parentOf(x)))
			}
		}
	}
	tree.root.color = black
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree TreeMap[K, V]) MarshalJSON() (jsonBytes []byte, err error) {
	if tree.root == nil {
		return []byte("null"), nil
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('{')
	tree.ForEach(func(key K, value V) bool {
		valueBytes, valueJsonErr := json.Marshal(value)
		if valueJsonErr != nil {
			err = valueJsonErr
			return false
		}
		if buffer.Len() > 1 {
			buffer.WriteByte(',')
		}
		buffer.WriteString(fmt.Sprintf(`"%v":%s`, key, valueBytes))
		return true
	})
	buffer.WriteByte('}')
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (tree *TreeMap[K, V]) UnmarshalJSON(b []byte) error {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.comparator == nil {
		tree.comparator = comparators.ComparatorAny[K]
	}
	var data map[K]V
	if err := json.UnmarshalUseNumber(b, &data); err != nil {
		return err
	}
	for k, v := range data {
		tree.insertEntry(k, v)
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (tree *TreeMap[K, V]) UnmarshalValue(value interface{}) (err error) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.comparator == nil {
		tree.comparator = comparators.ComparatorAny[K]
	}
	for k, v := range gconv.Map(value) {
		kt := gconv.ConvertGeneric[K](k)
		var vt V
		switch v.(type) {
		case string, []byte, json2.Number:
			var ok bool
			if vt, ok = v.(V); !ok {
				if err = json.UnmarshalUseNumber(gconv.Bytes(v), &vt); err != nil {
					return err
				}
			}
		default:
			vt, _ = v.(V)
		}
		tree.insertEntry(kt, vt)
	}
	return
}

// getComparator returns the comparator if it's previously set,
// or else it panics.
func (tree *TreeMap[K, V]) getComparator() func(a, b K) int {
	if tree.comparator == nil {
		return comparators.ComparatorAny[K]
	}
	return tree.comparator
}
