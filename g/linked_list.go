// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with l file,
// You can obtain one at https://github.com/gogf/gf.
//

package g

import (
	"bytes"
	json2 "encoding/json"

	"github.com/wesleywu/gcontainer/internal/deepcopy"
	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/internal/rwmutex"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

// LinkedList represents a doubly linked list.
// The zero value for LinkedList is an empty list ready to use.
type LinkedList[T comparable] struct {
	mu   rwmutex.RWMutex
	root Element[T] // sentinel list element, only &root, root.prev, and root.next are used
	len  int        // current list length excluding (this) sentinel element
}

// Element is an element of a linked list.
type Element[T comparable] struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element[T]

	// The list to which this element belongs.
	list *LinkedList[T]

	// The value stored with this element.
	Value T
}

// Init initializes or clears list l.
func (l *LinkedList[T]) Init() *LinkedList[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// NewLinkedList returns an initialized list.
func NewLinkedList[T comparable](safe ...bool) *LinkedList[T] {
	l := new(LinkedList[T]).Init()
	l.mu = rwmutex.Create(safe...)
	return l
}

// NewLinkedListFrom creates and returns a list from a copy of given slice `array`.
// The parameter `safe` is used to specify whether using list in concurrent-safety,
// which is false in default.
func NewLinkedListFrom[T comparable](array []T, safe ...bool) *LinkedList[T] {
	l := NewLinkedList[T](safe...)
	for _, v := range array {
		l.PushBack(v)
	}
	return l
}

// Next returns the next list element or nil.
func (e *Element[T]) Next() *Element[T] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element[T]) Prev() *Element[T] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Add append a new element e with value v at the back of list l and returns true.
func (l *LinkedList[T]) Add(values ...T) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	for _, value := range values {
		_ = l.insertValue(value, l.root.prev)
	}
	return true
}

// AddAll adds all the elements in the specified collection to this list.
// Returns true if this collection changed as a result of the call
func (l *LinkedList[T]) AddAll(values Collection[T]) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	values.ForEach(func(value T) bool {
		_ = l.insertValue(value, l.root.prev)
		return true
	})
	return true
}

// Contains returns true if this collection contains the specified element.
func (l *LinkedList[T]) Contains(value T) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	found := false
	length := l.len
	if length > 0 {
		for i, e := 0, l.root.next; i < length; i, e = i+1, e.Next() {
			if e.Value == value {
				found = true
				break
			}
		}
	}
	return found
}

// ContainsAll returns true if this collection contains all the elements in the specified collection.
func (l *LinkedList[T]) ContainsAll(values Collection[T]) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	foundMap := make(map[T]bool, 0)
	values.ForEach(func(t T) bool {
		foundMap[t] = false
		return true
	})
	length := l.len
	if length > 0 {
		for i, e := 0, l.root.next; i < length; i, e = i+1, e.Next() {
			if _, ok := foundMap[e.Value]; ok {
				foundMap[e.Value] = true
				break
			}
		}
	}
	for _, found := range foundMap {
		if !found {
			return false
		}
	}
	return true
}

// ForEach iterates all elements in this collection readonly with custom callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (l *LinkedList[T]) ForEach(f func(T) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	for i, e := 0, l.root.next; i < l.len; i, e = i+1, e.Next() {
		if !f(e.Value) {
			break
		}
	}
}

// IsEmpty returns true if this collection contains no elements.
func (l *LinkedList[T]) IsEmpty() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	return l.len == 0
}

// Slice returns an array containing shadow copy of all the elements in this list.
func (l *LinkedList[T]) Slice() []T {
	return l.FrontAll()
}

// search returns the matching element in this list, or nil if the element can not be found.
func (l *LinkedList[T]) search(value T) *Element[T] {
	if l.len > 0 {
		for i, e := 0, l.root.next; i < l.len; i, e = i+1, e.Next() {
			if e.Value == value {
				return e
			}
		}
	}
	return nil
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *LinkedList[T]) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.len
}

// Size is alias of Len.
func (l *LinkedList[T]) Size() int {
	return l.Len()
}

// Front returns the first element of list l or nil if the list is empty.
func (l *LinkedList[T]) Front() *Element[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *LinkedList[T]) Back() *Element[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero LinkedList value.
func (l *LinkedList[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *LinkedList[T]) insert(e, at *Element[T]) *Element[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{value: v}, at).
func (l *LinkedList[T]) insertValue(v T, at *Element[T]) *Element[T] {
	return l.insert(&Element[T]{Value: v}, at)
}

// remove removes e from its list, decrements l.len
func (l *LinkedList[T]) remove(e *Element[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
}

// move moves e to next to at.
func (l *LinkedList[T]) move(e, at *Element[T]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// Remove removes all of this list's elements that are also contained in the specified slice
// if it is present.
// Returns true if this collection changed as a result of the call
func (l *LinkedList[T]) Remove(values ...T) (changed bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	changed = false
	for _, value := range values {
		existing := l.search(value)
		if existing != nil {
			l.remove(existing)
			changed = true
		}
	}
	return
}

// RemoveAll removes all of this list's elements that are also contained in the specified collection
// Returns true if this collection changed as a result of the call
func (l *LinkedList[T]) RemoveAll(values Collection[T]) (changed bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	changed = false
	values.ForEach(func(value T) bool {
		existing := l.search(value)
		if existing != nil {
			l.remove(existing)
			changed = true
		}
		return true
	})
	return
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *LinkedList[T]) PushBack(v T) *Element[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *LinkedList[T]) PushFront(v T) *Element[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBacks inserts multiple new elements with values `values` at the back of list `l`.
func (l *LinkedList[T]) PushBacks(values []T) {
	l.mu.Lock()
	l.mu.Unlock()
	l.lazyInit()
	for _, v := range values {
		l.PushBack(v)
	}
}

// PushFronts inserts multiple new elements with values `values` at the front of list `l`.
func (l *LinkedList[T]) PushFronts(values []T) {
	l.mu.Lock()
	defer l.mu.RUnlock()
	l.lazyInit()
	for _, v := range values {
		l.PushFront(v)
	}
}

// PopBack removes the element from back of `l` and returns the value of the element.
func (l *LinkedList[T]) PopBack() (value T, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	if e := l.root.prev; e != nil {
		value = e.Value
		if e.list == l {
			// if e.list == l, l must have been initialized when e was inserted
			// in l or l == nil (e is a zero Element) and l.remove will crash
			l.remove(e)
			return value, true
		}
	}
	return value, false
}

// PopFront removes the element from front of `l` and returns the value of the element.
func (l *LinkedList[T]) PopFront() (value T, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	if e := l.root.next; e != nil {
		value = e.Value
		if e.list == l {
			// if e.list == l, l must have been initialized when e was inserted
			// in l or l == nil (e is a zero Element) and l.remove will crash
			l.remove(e)
			return value, true
		}
	}
	return value, false
}

// PopBacks removes `max` elements from back of `l`
// and returns values of the removed elements as slice.
func (l *LinkedList[T]) PopBacks(max int) (values []T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	length := l.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]T, length)
		for i := 0; i < length; i++ {
			back := l.root.prev
			values[i] = back.Value
			if back.list == l {
				// if e.list == l, l must have been initialized when e was inserted
				// in l or l == nil (e is a zero Element) and l.remove will crash
				l.remove(back)
			}
		}
	}
	return
}

// PopFronts removes `max` elements from front of `l`
// and returns values of the removed elements as slice.
func (l *LinkedList[T]) PopFronts(max int) (values []T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	length := l.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]T, length)
		for i := 0; i < length; i++ {
			front := l.root.next
			values[i] = front.Value
			if front.list == l {
				// if e.list == l, l must have been initialized when e was inserted
				// in l or l == nil (e is a zero Element) and l.remove will crash
				l.remove(front)
			}
		}
	}
	return
}

// PopBackAll removes all elements from back of `l`
// and returns values of the removed elements as slice.
func (l *LinkedList[T]) PopBackAll() []T {
	return l.PopBacks(-1)
}

// PopFrontAll removes all elements from front of `l`
// and returns values of the removed elements as slice.
func (l *LinkedList[T]) PopFrontAll() []T {
	return l.PopFronts(-1)
}

// FrontAll copies and returns values of all elements from front of `l` as slice.
func (l *LinkedList[T]) FrontAll() (values []T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	length := l.Len()
	if length > 0 {
		values = make([]T, length)
		for i, e := 0, l.root.next; i < length; i, e = i+1, e.Next() {
			values[i] = e.Value
		}
	}
	return
}

// BackAll copies and returns values of all elements from back of `l` as slice.
func (l *LinkedList[T]) BackAll() (values []T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	length := l.Len()
	if length > 0 {
		values = make([]T, length)
		for i, e := 0, l.root.prev; i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value
		}
	}
	return
}

// FrontValue returns value of the first element of `l` or nil if the list is empty.
func (l *LinkedList[T]) FrontValue() (value T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	if e := l.root.next; e != nil {
		value = e.Value
	}
	return
}

// BackValue returns value of the last element of `l` or nil if the list is empty.
func (l *LinkedList[T]) BackValue() (value T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	if e := l.root.prev; e != nil {
		value = e.Value
	}
	return
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *LinkedList[T]) InsertBefore(mark *Element[T], v T) *Element[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if mark.list != l {
		return nil
	}
	// see comment in LinkedList.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *LinkedList[T]) InsertAfter(mark *Element[T], v T) *Element[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if mark.list != l {
		return nil
	}
	// see comment in LinkedList.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *LinkedList[T]) MoveToFront(e *Element[T]) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in LinkedList.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *LinkedList[T]) MoveToBack(e *Element[T]) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in LinkedList.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *LinkedList[T]) MoveBefore(e, mark *Element[T]) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *LinkedList[T]) MoveAfter(e, mark *Element[T]) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *LinkedList[T]) PushBackList(other *LinkedList[T]) {
	if l != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	for i, e := other.len, other.root.next; i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *LinkedList[T]) PushFrontList(other *LinkedList[T]) {
	if l != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	for i, e := other.len, other.root.prev; i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

// Clear removes all the elements from this collection.
func (l *LinkedList[T]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Init()
}

func (l *LinkedList[T]) Clone() Collection[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	values := make([]T, l.len)
	for i, e := 0, l.root.next; i < l.len; i, e = i+1, e.Next() {
		values[i] = e.Value
	}
	return NewLinkedListFrom(values, l.mu.IsSafe())
}

func (l *LinkedList[T]) Equals(another Collection[T]) bool {
	if l == another {
		return true
	}
	var (
		ano *LinkedList[T]
		ok  bool
	)
	if ano, ok = another.(*LinkedList[T]); !ok {
		return false
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	ano.mu.RLock()
	defer ano.mu.RUnlock()
	if l.len != ano.len {
		return false
	}
	values := make([]T, l.len)
	for i, e := 0, l.root.next; i < l.len; i, e = i+1, e.Next() {
		values[i] = e.Value
	}
	valuesAno := make([]T, l.len)
	for i, e := 0, ano.root.next; i < ano.len; i, e = i+1, e.Next() {
		valuesAno[i] = e.Value
	}
	for i := 0; i < l.len; i++ {
		if values[i] != valuesAno[i] {
			return false
		}
	}
	return true
}

// ForEachAsc iterates the list readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (l *LinkedList[T]) ForEachAsc(f func(e T) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	for i, e := 0, l.root.next; i < l.len; i, e = i+1, e.Next() {
		if !f(e.Value) {
			break
		}
	}
}

// ForEachDesc iterates the list readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (l *LinkedList[T]) ForEachDesc(f func(e T) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	for i, e := 0, l.root.prev; i < l.len; i, e = i+1, e.Prev() {
		if !f(e.Value) {
			break
		}
	}
}

// Join joins list elements with a string `glue`.
func (l *LinkedList[T]) Join(glue string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.lazyInit()
	buffer := bytes.NewBuffer(nil)
	for i, e := 0, l.root.next; i < l.len; i, e = i+1, e.Next() {
		buffer.WriteString(gconv.String(e.Value))
		if i != l.len-1 {
			buffer.WriteString(glue)
		}
	}
	return buffer.String()
}

// String returns current list as a string.
func (l *LinkedList[T]) String() string {
	l.lazyInit()
	return "[" + l.Join(",") + "]"
}

// Sum returns the sum of values in an array.
func (l *LinkedList[T]) Sum() (sum int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for i, e := 0, l.root.next; i < l.len; i, e = i+1, e.Next() {
		sum += gconv.Int(e.Value)
	}
	return
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (l LinkedList[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.FrontAll())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (l *LinkedList[T]) UnmarshalJSON(b []byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	var array []T
	if err := json.UnmarshalUseNumber(b, &array); err != nil {
		return err
	}
	for _, v := range array {
		l.PushBack(v)
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for list.
func (l *LinkedList[T]) UnmarshalValue(value interface{}) (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lazyInit()
	var array []T
	switch value.(type) {
	case string, []byte, json2.Number:
		err = json.UnmarshalUseNumber(gconv.Bytes(value), &array)
	default:
		array = gconv.SliceAny[T](value)
	}
	for _, v := range array {
		l.PushBack(v)
	}
	return err
}

// DeepCopy implements interface for deep copy of current type.
func (l *LinkedList[T]) DeepCopy() Collection[T] {
	if l == nil {
		return nil
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	var (
		length = l.Len()
		values = make([]T, length)
	)
	if length > 0 {
		for i, e := 0, l.root.next; i < length; i, e = i+1, e.Next() {
			values[i] = deepcopy.Copy(e.Value).(T)
		}
	}
	return NewLinkedListFrom[T](values, l.mu.IsSafe())
}
