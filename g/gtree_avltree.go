// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
	"bytes"
	"fmt"

	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/internal/rwmutex"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

// AVLTree holds elements of the AVL tree.
type AVLTree[K comparable, V any] struct {
	mu         rwmutex.RWMutex
	root       *AVLTreeNode[K, V]
	comparator func(v1, v2 K) int
	size       int
}

// AVLTreeNode is a single element within the tree.
type AVLTreeNode[K comparable, V any] struct {
	key      K
	value    V
	parent   *AVLTreeNode[K, V]
	children [2]*AVLTreeNode[K, V]
	b        int8
}

// NewAVLTree instantiates an AVL tree with the custom key comparators.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewAVLTree[K comparable, V any](comparator func(v1, v2 K) int, safe ...bool) *AVLTree[K, V] {
	return &AVLTree[K, V]{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
	}
}

// NewAVLTreeFrom instantiates an AVL tree with the custom key comparators and data map.
// The parameter `safe` is used to specify whether using tree in concurrent-safety,
// which is false in default.
func NewAVLTreeFrom[K comparable, V any](comparator func(v1, v2 K) int, data map[K]V, safe ...bool) *AVLTree[K, V] {
	tree := NewAVLTree[K, V](comparator, safe...)
	for k, v := range data {
		tree.put(k, v, nil, &tree.root)
	}
	return tree
}

func (n *AVLTreeNode[K, V]) Key() K {
	return n.key
}

func (n *AVLTreeNode[K, V]) Value() V {
	return n.value
}

// Clone returns a new tree with a copy of current tree.
func (tree *AVLTree[K, V]) Clone(safe ...bool) Map[K, V] {
	newTree := NewAVLTree[K, V](tree.comparator, safe...)
	newTree.Puts(tree.Map())
	return newTree
}

// Put inserts node into the tree.
func (tree *AVLTree[K, V]) Put(key K, value V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.put(key, value, nil, &tree.root)
}

// Puts batch sets key-values to the tree.
func (tree *AVLTree[K, V]) Puts(data map[K]V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for key, value := range data {
		tree.put(key, value, nil, &tree.root)
	}
}

// Search searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *AVLTree[K, V]) Search(key K) (value V, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if node, found := tree.doSearch(key); found {
		return node.value, true
	}
	return
}

// doSearch searches the tree with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (tree *AVLTree[K, V]) doSearch(key K) (node *AVLTreeNode[K, V], found bool) {
	node = tree.root
	for node != nil {
		cmp := tree.getComparator()(key, node.key)
		switch {
		case cmp == 0:
			return node, true
		case cmp < 0:
			node = node.children[0]
		case cmp > 0:
			node = node.children[1]
		}
	}
	return nil, false
}

// Get returns the value by given `key`, or empty value of type K if the key is not found in the map.
func (tree *AVLTree[K, V]) Get(key K) (value V) {
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
func (tree *AVLTree[K, V]) doSetWithLockCheck(key K, value V) V {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if node, found := tree.doSearch(key); found {
		return node.value
	}
	if any(value) != nil {
		tree.put(key, value, nil, &tree.root)
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
func (tree *AVLTree[K, V]) doSetWithLockCheckFunc(key K, f func() V) V {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if node, found := tree.doSearch(key); found {
		return node.value
	}
	value := f()
	if any(value) != nil {
		tree.put(key, value, nil, &tree.root)
	}
	return value
}

// GetOrPut returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (tree *AVLTree[K, V]) GetOrPut(key K, value V) V {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrPutFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
func (tree *AVLTree[K, V]) GetOrPutFunc(key K, f func() V) V {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheckFunc(key, f)
	} else {
		return v
	}
}

// PutIfAbsent sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (tree *AVLTree[K, V]) PutIfAbsent(key K, value V) bool {
	if !tree.ContainsKey(key) {
		tree.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// PutIfAbsentFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (tree *AVLTree[K, V]) PutIfAbsentFunc(key K, f func() V) bool {
	if !tree.ContainsKey(key) {
		tree.doSetWithLockCheckFunc(key, f)
		return true
	}
	return false
}

// ContainsKey checks whether `key` exists in the tree.
func (tree *AVLTree[K, V]) ContainsKey(key K) bool {
	_, ok := tree.Search(key)
	return ok
}

// Remove removes the node from the tree by key.
// key should adhere to the comparators's type assertion, otherwise method panics.
func (tree *AVLTree[K, V]) Remove(key K) (value V, removed bool) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.remove(key, &tree.root)
}

// Removes batch deletes values of the tree by `keys`.
func (tree *AVLTree[K, V]) Removes(keys []K) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.remove(key, &tree.root)
	}
}

// IsEmpty returns true if tree does not contain any nodes.
func (tree *AVLTree[K, V]) IsEmpty() bool {
	return tree.Size() == 0
}

// Size returns number of nodes in the tree.
func (tree *AVLTree[K, V]) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.size
}

// Keys returns all keys in asc order.
func (tree *AVLTree[K, V]) Keys() []K {
	keys := make([]K, tree.Size())
	index := 0
	tree.ForEachAsc(func(key K, value V) bool {
		keys[index] = key
		index++
		return true
	})
	return keys
}

// Values returns all values in asc order based on the key.
func (tree *AVLTree[K, V]) Values() []V {
	values := make([]V, tree.Size())
	index := 0
	tree.ForEachAsc(func(key K, value V) bool {
		values[index] = value
		index++
		return true
	})
	return values
}

// KeySet returns a set of the keys contained in the tree.
func (tree *AVLTree[K, V]) KeySet() Set[K] {
	return NewHashSetFrom(tree.Keys(), tree.mu.IsSafe())
}

// Left returns the minimum element of the AVL tree
// or nil if the tree is empty.
func (tree *AVLTree[K, V]) Left() *AVLTreeNode[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.bottom(0)
	if tree.mu.IsSafe() {
		return &AVLTreeNode[K, V]{
			key:   node.key,
			value: node.value,
		}
	}
	return node
}

// Right returns the maximum element of the AVL tree
// or nil if the tree is empty.
func (tree *AVLTree[K, V]) Right() *AVLTreeNode[K, V] {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.bottom(1)
	if tree.mu.IsSafe() {
		return &AVLTreeNode[K, V]{
			key:   node.key,
			value: node.value,
		}
	}
	return node
}

// Floor Finds floor node of the input key, return the floor node or nil if no floor node is found.
// Second return parameter is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree is larger than the given node.
//
// key should adhere to the comparators's type assertion, otherwise method panics.
func (tree *AVLTree[K, V]) Floor(key K) (floor *AVLTreeNode[K, V], found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	n := tree.root
	for n != nil {
		c := tree.getComparator()(key, n.key)
		switch {
		case c == 0:
			return n, true
		case c < 0:
			n = n.children[0]
		case c > 0:
			floor, found = n, true
			n = n.children[1]
		}
	}
	if found {
		return
	}
	return nil, false
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling node is found.
// Second return parameter is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree is smaller than the given node.
//
// key should adhere to the comparators's type assertion, otherwise method panics.
func (tree *AVLTree[K, V]) Ceiling(key K) (ceiling *AVLTreeNode[K, V], found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	n := tree.root
	for n != nil {
		c := tree.getComparator()(key, n.key)
		switch {
		case c == 0:
			return n, true
		case c > 0:
			n = n.children[1]
		case c < 0:
			ceiling, found = n, true
			n = n.children[0]
		}
	}
	if found {
		return
	}
	return nil, false
}

// Clear removes all nodes from the tree.
func (tree *AVLTree[K, V]) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
}

// Replace the data of the tree with given `data`.
func (tree *AVLTree[K, V]) Replace(data map[K]V) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
	for key, value := range data {
		tree.put(key, value, nil, &tree.root)
	}
}

// String returns a string representation of container
func (tree *AVLTree[K, V]) String() string {
	if tree == nil {
		return ""
	}
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	str := ""
	if tree.size != 0 {
		output(tree.root, "", true, &str)
	}
	return str
}

// Print prints the tree to stdout.
func (tree *AVLTree[K, V]) Print() {
	fmt.Println(tree.String())
}

// Map returns all key-value items as map.
func (tree *AVLTree[K, V]) Map() map[K]V {
	m := make(map[K]V, tree.Size())
	tree.ForEachAsc(func(key K, value V) bool {
		m[key] = value
		return true
	})
	return m
}

// MapStrAny returns all key-value items as map[string]V.
func (tree *AVLTree[K, V]) MapStrAny() map[string]V {
	m := make(map[string]V, tree.Size())
	tree.ForEachAsc(func(key K, value V) bool {
		m[gconv.String(key)] = value
		return true
	})
	return m
}

// ForEach is alias of ForEachAsc.
func (tree *AVLTree[K, V]) ForEach(f func(key K, value V) bool) {
	tree.ForEachAsc(f)
}

// IteratorFrom is alias of IteratorAscFrom.
func (tree *AVLTree[K, V]) IteratorFrom(key K, match bool, f func(key K, value V) bool) {
	tree.IteratorAscFrom(key, match, f)
}

// ForEachAsc iterates the tree readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *AVLTree[K, V]) ForEachAsc(f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorAsc(tree.bottom(0), f)
}

// IteratorAscFrom iterates the tree readonly in ascending order with given callback function `f`.
// The parameter `key` specifies the start entry for iterating. The `match` specifies whether
// starting iterating if the `key` is fully matched, or else using index searching iterating.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *AVLTree[K, V]) IteratorAscFrom(key K, match bool, f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.doSearch(key)
	if match {
		if found {
			tree.doIteratorAsc(node, f)
		}
	} else {
		tree.doIteratorAsc(node, f)
	}
}

func (tree *AVLTree[K, V]) doIteratorAsc(node *AVLTreeNode[K, V], f func(key K, value V) bool) {
	for node != nil {
		if !f(node.key, node.value) {
			return
		}
		node = node.Next()
	}
}

// ForEachDesc iterates the tree readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *AVLTree[K, V]) ForEachDesc(f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorDesc(tree.bottom(1), f)
}

// IteratorDescFrom iterates the tree readonly in descending order with given callback function `f`.
// The parameter `key` specifies the start entry for iterating. The `match` specifies whether
// starting iterating if the `key` is fully matched, or else using index searching iterating.
// If `f` returns true, then it continues iterating; or false to stop.
func (tree *AVLTree[K, V]) IteratorDescFrom(key K, match bool, f func(key K, value V) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.doSearch(key)
	if match {
		if found {
			tree.doIteratorDesc(node, f)
		}
	} else {
		tree.doIteratorDesc(node, f)
	}
}

func (tree *AVLTree[K, V]) doIteratorDesc(node *AVLTreeNode[K, V], f func(key K, value V) bool) {
	for node != nil {
		if !f(node.key, node.value) {
			return
		}
		node = node.Prev()
	}
}

func (tree *AVLTree[K, V]) put(key K, value V, p *AVLTreeNode[K, V], qp **AVLTreeNode[K, V]) bool {
	q := *qp
	if q == nil {
		tree.size++
		*qp = &AVLTreeNode[K, V]{key: key, value: value, parent: p}
		return true
	}

	c := tree.getComparator()(key, q.key)
	if c == 0 {
		q.key = key
		q.value = value
		return false
	}

	if c < 0 {
		c = -1
	} else {
		c = 1
	}
	a := (c + 1) / 2
	if tree.put(key, value, q, &q.children[a]) {
		return putFix(int8(c), qp)
	}
	return false
}

func (tree *AVLTree[K, V]) remove(key K, qp **AVLTreeNode[K, V]) (value V, fix bool) {
	q := *qp
	if q == nil {
		return
	}

	c := tree.getComparator()(key, q.key)
	if c == 0 {
		tree.size--
		value = q.value
		fix = true
		if q.children[1] == nil {
			if q.children[0] != nil {
				q.children[0].parent = q.parent
			}
			*qp = q.children[0]
			return
		}
		if removeMin(&q.children[1], &q.key, &q.value) {
			removeFix(-1, qp)
			return value, true
		}
		return
	}

	if c < 0 {
		c = -1
	} else {
		c = 1
	}
	a := (c + 1) / 2
	value, fix = tree.remove(key, &q.children[a])
	if fix {
		removeFix(int8(-c), qp)
		return value, true
	}
	return value, false
}

func removeMin[K comparable, V any](qp **AVLTreeNode[K, V], minKey *K, minVal *V) bool {
	q := *qp
	if q.children[0] == nil {
		*minKey = q.key
		*minVal = q.value
		if q.children[1] != nil {
			q.children[1].parent = q.parent
		}
		*qp = q.children[1]
		return true
	}
	fix := removeMin(&q.children[0], minKey, minVal)
	if fix {
		return removeFix(1, qp)
	}
	return false
}

func putFix[K comparable, V any](c int8, t **AVLTreeNode[K, V]) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return true
	}

	if s.b == -c {
		s.b = 0
		return false
	}

	if s.children[(c+1)/2].b == c {
		s = singleRotate(c, s)
	} else {
		s = doubleRotate(c, s)
	}
	*t = s
	return false
}

func removeFix[K comparable, V any](c int8, t **AVLTreeNode[K, V]) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return false
	}

	if s.b == -c {
		s.b = 0
		return true
	}

	a := (c + 1) / 2
	if s.children[a].b == 0 {
		s = rotate[K, V](c, s)
		s.b = -c
		*t = s
		return false
	}

	if s.children[a].b == c {
		s = singleRotate[K, V](c, s)
	} else {
		s = doubleRotate[K, V](c, s)
	}
	*t = s
	return true
}

func singleRotate[K comparable, V any](c int8, s *AVLTreeNode[K, V]) *AVLTreeNode[K, V] {
	s.b = 0
	s = rotate(c, s)
	s.b = 0
	return s
}

func doubleRotate[K comparable, V any](c int8, s *AVLTreeNode[K, V]) *AVLTreeNode[K, V] {
	a := (c + 1) / 2
	r := s.children[a]
	s.children[a] = rotate(-c, s.children[a])
	p := rotate(c, s)

	switch {
	default:
		s.b = 0
		r.b = 0
	case p.b == c:
		s.b = -c
		r.b = 0
	case p.b == -c:
		s.b = 0
		r.b = c
	}

	p.b = 0
	return p
}

func rotate[K comparable, V any](c int8, s *AVLTreeNode[K, V]) *AVLTreeNode[K, V] {
	a := (c + 1) / 2
	r := s.children[a]
	s.children[a] = r.children[a^1]
	if s.children[a] != nil {
		s.children[a].parent = s
	}
	r.children[a^1] = s
	r.parent = s.parent
	s.parent = r
	return r
}

func (tree *AVLTree[K, V]) bottom(d int) *AVLTreeNode[K, V] {
	n := tree.root
	if n == nil {
		return nil
	}

	for c := n.children[d]; c != nil; c = n.children[d] {
		n = c
	}
	return n
}

// Prev returns the previous element in an inorder
// walk of the AVL tree.
func (node *AVLTreeNode[K, V]) Prev() *AVLTreeNode[K, V] {
	return node.walk1(0)
}

// Next returns the next element in an inorder
// walk of the AVL tree.
func (node *AVLTreeNode[K, V]) Next() *AVLTreeNode[K, V] {
	return node.walk1(1)
}

func (node *AVLTreeNode[K, V]) walk1(a int) *AVLTreeNode[K, V] {
	if node == nil {
		return nil
	}
	n := node
	if n.children[a] != nil {
		n = n.children[a]
		for n.children[a^1] != nil {
			n = n.children[a^1]
		}
		return n
	}

	p := n.parent
	for p != nil && p.children[a] == n {
		n = p
		p = p.parent
	}
	return p
}

func output[K comparable, V any](node *AVLTreeNode[K, V], prefix string, isTail bool, str *string) {
	if node.children[1] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.children[1], newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += fmt.Sprintf("%v\n", node.Key())
	if node.children[0] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.children[0], newPrefix, true, str)
	}
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (tree AVLTree[K, V]) MarshalJSON() (jsonBytes []byte, err error) {
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

// getComparator returns the comparator if it's previously set,
// or else it panics.
func (tree *AVLTree[K, V]) getComparator() func(a, b K) int {
	if tree.comparator == nil {
		panic("comparators is missing for tree")
	}
	return tree.comparator
}
