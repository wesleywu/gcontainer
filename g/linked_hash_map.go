// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
	"bytes"
	json2 "encoding/json"
	"fmt"

	"github.com/wesleywu/gcontainer/internal/deepcopy"
	"github.com/wesleywu/gcontainer/internal/json"
	"github.com/wesleywu/gcontainer/internal/rwmutex"
	"github.com/wesleywu/gcontainer/utils/empty"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

// LinkedHashMap is a map that preserves insertion-order.
//
// It is backed by a hash table to store values and doubly-linked list to store ordering.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
type LinkedHashMap[K comparable, V any] struct {
	mu   rwmutex.RWMutex
	data map[K]*Element[*gListMapNode[K, V]]
	list *LinkedList[*gListMapNode[K, V]]
}

type gListMapNode[K comparable, V any] struct {
	key   K
	value V
}

// NewListMap returns an empty link map.
// LinkedHashMap is backed by a hash table to store values and doubly-linked list to store ordering.
// The parameter `safe` is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewListMap[K comparable, V any](safe ...bool) *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{
		mu:   rwmutex.Create(safe...),
		data: make(map[K]*Element[*gListMapNode[K, V]]),
		list: NewLinkedList[*gListMapNode[K, V]](),
	}
}

// NewListMapFrom returns a link map from given map `data`.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewListMapFrom[K comparable, V any](data map[K]V, safe ...bool) *LinkedHashMap[K, V] {
	m := NewListMap[K, V](safe...)
	m.Puts(data)
	return m
}

// ForEach is alias of ForEachAsc.
func (m *LinkedHashMap[K, V]) ForEach(f func(key K, value V) bool) {
	m.ForEachAsc(f)
}

// ForEachAsc iterates the map readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *LinkedHashMap[K, V]) ForEachAsc(f func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.list != nil {
		m.list.ForEachAsc(func(node *gListMapNode[K, V]) bool {
			return f(node.key, node.value)
		})
	}
}

// ForEachDesc iterates the map readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *LinkedHashMap[K, V]) ForEachDesc(f func(key K, value interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.list != nil {
		m.list.ForEachDesc(func(node *gListMapNode[K, V]) bool {
			return f(node.key, node.value)
		})
	}
}

// Clone returns a new link map with copy of current map data.
func (m *LinkedHashMap[K, V]) Clone(safe ...bool) Map[K, V] {
	return NewListMapFrom[K, V](m.Map(), safe...)
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *LinkedHashMap[K, V]) Clear() {
	m.mu.Lock()
	m.data = make(map[K]*Element[*gListMapNode[K, V]])
	m.list = NewLinkedList[*gListMapNode[K, V]]()
	m.mu.Unlock()
}

// Replace the data of the map with given `data`.
func (m *LinkedHashMap[K, V]) Replace(data map[K]V) {
	m.mu.Lock()
	m.data = make(map[K]*Element[*gListMapNode[K, V]])
	m.list = NewLinkedList[*gListMapNode[K, V]]()
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
		} else {
			e.Value = &gListMapNode[K, V]{key, value}
		}
	}
	m.mu.Unlock()
}

// Map returns a copy of the underlying data of the map.
func (m *LinkedHashMap[K, V]) Map() map[K]V {
	m.mu.RLock()
	var data map[K]V
	if m.list != nil {
		data = make(map[K]V, len(m.data))
		m.list.ForEachAsc(func(node *gListMapNode[K, V]) bool {
			data[node.key] = node.value
			return true
		})
	}
	m.mu.RUnlock()
	return data
}

// MapStrAny returns a copy of the underlying data of the map as map[string]V.
func (m *LinkedHashMap[K, V]) MapStrAny() map[string]V {
	m.mu.RLock()
	var data map[string]V
	if m.list != nil {
		data = make(map[string]V, len(m.data))
		m.list.ForEachAsc(func(node *gListMapNode[K, V]) bool {
			data[gconv.String(node.key)] = node.value
			return true
		})
	}
	m.mu.RUnlock()
	return data
}

// FilterEmpty deletes all key-value pair of which the value is empty.
func (m *LinkedHashMap[K, V]) FilterEmpty() {
	m.mu.Lock()
	if m.list != nil {
		var (
			keys = make([]K, 0)
		)
		m.list.ForEachAsc(func(node *gListMapNode[K, V]) bool {
			if empty.IsEmpty(node.value) {
				keys = append(keys, node.key)
			}
			return true
		})
		if len(keys) > 0 {
			for _, key := range keys {
				if e, ok := m.data[key]; ok {
					delete(m.data, key)
					m.list.Remove(e.Value)
				}
			}
		}
	}
	m.mu.Unlock()
}

// Put sets key-value to the map.
func (m *LinkedHashMap[K, V]) Put(key K, value V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]*Element[*gListMapNode[K, V]])
		m.list = NewLinkedList[*gListMapNode[K, V]]()
	}
	if e, ok := m.data[key]; !ok {
		m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
	} else {
		e.Value = &gListMapNode[K, V]{key, value}
	}
	m.mu.Unlock()
}

// Puts batch sets key-values to the map.
func (m *LinkedHashMap[K, V]) Puts(data map[K]V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]*Element[*gListMapNode[K, V]])
		m.list = NewLinkedList[*gListMapNode[K, V]]()
	}
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
		} else {
			e.Value = &gListMapNode[K, V]{key, value}
		}
	}
	m.mu.Unlock()
}

// Search searches the map with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (m *LinkedHashMap[K, V]) Search(key K) (value V, found bool) {
	m.mu.RLock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.value
			found = ok
		}
	}
	m.mu.RUnlock()
	return
}

// Get returns the value by given `key`, or empty value of type K if the key is not found in the map.
func (m *LinkedHashMap[K, V]) Get(key K) (value V) {
	m.mu.RLock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.value
		}
	}
	m.mu.RUnlock()
	return
}

// Pop retrieves and deletes an item from the map.
func (m *LinkedHashMap[K, V]) Pop() (key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, e := range m.data {
		value = e.Value.value
		delete(m.data, k)
		m.list.Remove(e.Value)
		return k, value
	}
	return
}

// Pops retrieves and deletes `size` items from the map.
// It returns all items if size == -1.
func (m *LinkedHashMap[K, V]) Pops(size int) map[K]V {
	m.mu.Lock()
	defer m.mu.Unlock()
	if size > len(m.data) || size == -1 {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	index := 0
	newMap := make(map[K]V, size)
	for k, e := range m.data {
		value := e.Value.value
		delete(m.data, k)
		m.list.Remove(e.Value)
		newMap[k] = value
		index++
		if index == size {
			break
		}
	}
	return newMap
}

// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given `key`,
// or else just return the existing value.
//
// When setting value, if `value` is type of `func() interface {}`,
// it will be executed with mutex.Lock of the map,
// and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (m *LinkedHashMap[K, V]) doSetWithLockCheck(key K, value V) V {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*Element[*gListMapNode[K, V]])
		m.list = NewLinkedList[*gListMapNode[K, V]]()
	}
	if e, ok := m.data[key]; ok {
		return e.Value.value
	}
	if f, ok := any(value).(func() V); ok {
		value = f()
	}
	if any(value) != nil {
		m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
	}
	return value
}

// doSetWithLockCheckFunc checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given `key`,
// or else just return the existing value.
//
// When setting value, if `value` is type of `func() interface {}`,
// it will be executed with mutex.Lock of the map,
// and its return value will be set to the map with `key`.
//
// It returns value with given `key`.
func (m *LinkedHashMap[K, V]) doSetWithLockCheckFunc(key K, f func() V) V {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*Element[*gListMapNode[K, V]])
		m.list = NewLinkedList[*gListMapNode[K, V]]()
	}
	if e, ok := m.data[key]; ok {
		return e.Value.value
	}
	var value V
	value = f()
	if any(value) != nil {
		m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
	}
	return value
}

// GetOrPut returns the value by key,
// or sets value with given `value` if it does not exist and then returns this value.
func (m *LinkedHashMap[K, V]) GetOrPut(key K, value V) V {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrPutFunc returns the value by key,
// or sets value with returned value of callback function `f` if it does not exist
// and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
// with mutex.Lock of the map.
func (m *LinkedHashMap[K, V]) GetOrPutFunc(key K, f func() V) V {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheckFunc(key, f)
	} else {
		return v
	}
}

// PutIfAbsent sets `value` to the map if the `key` does not exist, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *LinkedHashMap[K, V]) PutIfAbsent(key K, value V) bool {
	if !m.ContainsKey(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// PutIfAbsentFunc sets value with return value of callback function `f`, and then returns true.
// It returns false if `key` exists, and `value` would be ignored.
func (m *LinkedHashMap[K, V]) PutIfAbsentFunc(key K, f func() V) bool {
	if !m.ContainsKey(key) {
		m.doSetWithLockCheckFunc(key, f)
		return true
	}
	return false
}

// Remove deletes value from map by given `key`, and return this deleted value.
func (m *LinkedHashMap[K, V]) Remove(key K) (value V, removed bool) {
	m.mu.Lock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.value
			delete(m.data, key)
			m.list.Remove(e.Value)
			removed = true
		}
	}
	m.mu.Unlock()
	return
}

// Removes batch deletes values of the map by keys.
func (m *LinkedHashMap[K, V]) Removes(keys []K) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			if e, ok := m.data[key]; ok {
				delete(m.data, key)
				m.list.Remove(e.Value)
			}
		}
	}
	m.mu.Unlock()
}

// Keys returns all keys of the map as a slice in ascending order.
func (m *LinkedHashMap[K, V]) Keys() []K {
	m.mu.RLock()
	var (
		keys  = make([]K, m.list.Len())
		index = 0
	)
	if m.list != nil {
		m.list.ForEachAsc(func(node *gListMapNode[K, V]) bool {
			keys[index] = node.key
			index++
			return true
		})
	}
	m.mu.RUnlock()
	return keys
}

// Values returns all values of the map as a slice.
func (m *LinkedHashMap[K, V]) Values() []V {
	m.mu.RLock()
	var (
		values = make([]V, m.list.Len())
		index  = 0
	)
	if m.list != nil {
		m.list.ForEachAsc(func(node *gListMapNode[K, V]) bool {
			values[index] = node.value
			index++
			return true
		})
	}
	m.mu.RUnlock()
	return values
}

// ContainsKey checks whether a key exists.
// It returns true if the `key` exists, or else false.
func (m *LinkedHashMap[K, V]) ContainsKey(key K) (ok bool) {
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Size returns the size of the map.
func (m *LinkedHashMap[K, V]) Size() (size int) {
	m.mu.RLock()
	size = len(m.data)
	m.mu.RUnlock()
	return
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *LinkedHashMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// Merge merges two link maps.
// The `other` map will be merged into the map `m`.
func (m *LinkedHashMap[K, V]) Merge(other *LinkedHashMap[K, V]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*Element[*gListMapNode[K, V]])
		m.list = NewLinkedList[*gListMapNode[K, V]]()
	}
	if other != m {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	other.list.ForEachAsc(func(node *gListMapNode[K, V]) bool {
		if e, ok := m.data[node.key]; !ok {
			m.data[node.key] = m.list.PushBack(&gListMapNode[K, V]{node.key, node.value})
		} else {
			e.Value = &gListMapNode[K, V]{node.key, node.value}
		}
		return true
	})
}

// String returns the map as a string.
func (m *LinkedHashMap[K, V]) String() string {
	if m == nil {
		return ""
	}
	b, _ := m.MarshalJSON()
	return string(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m LinkedHashMap[K, V]) MarshalJSON() (jsonBytes []byte, err error) {
	if m.data == nil {
		return []byte("null"), nil
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte('{')
	m.ForEach(func(key K, value V) bool {
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
func (m *LinkedHashMap[K, V]) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*Element[*gListMapNode[K, V]])
		m.list = NewLinkedList[*gListMapNode[K, V]]()
	}
	var data map[K]V
	if err := json.UnmarshalUseNumber(b, &data); err != nil {
		return err
	}
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&gListMapNode[K, V]{key, value})
		} else {
			e.Value = &gListMapNode[K, V]{key, value}
		}
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *LinkedHashMap[K, V]) UnmarshalValue(value interface{}) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]*Element[*gListMapNode[K, V]])
		m.list = NewLinkedList[*gListMapNode[K, V]]()
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
		if e, ok := m.data[kt]; !ok {
			m.data[kt] = m.list.PushBack(&gListMapNode[K, V]{kt, vt})
		} else {
			e.Value = &gListMapNode[K, V]{kt, vt}
		}
	}
	return
}

// DeepCopy implements interface for deep copy of current type.
func (m *LinkedHashMap[K, V]) DeepCopy() interface{} {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[K]V, len(m.data))
	if m.list != nil {
		m.list.ForEachAsc(func(e *gListMapNode[K, V]) bool {
			data[e.key] = deepcopy.Copy(e.value).(V)
			return true
		})
	}
	return NewListMapFrom(data, m.mu.IsSafe())
}
