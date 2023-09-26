// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"sync"
	"time"
)

type adapterMemoryData[K comparable, V any] struct {
	mu   sync.RWMutex               // dataMu ensures the concurrent safety of underlying data map.
	data map[K]adapterMemoryItem[V] // data is the underlying cache data which is stored in a hash table.
}

func newAdapterMemoryData[K comparable, V any]() *adapterMemoryData[K, V] {
	return &adapterMemoryData[K, V]{
		data: make(map[K]adapterMemoryItem[V]),
	}
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func (d *adapterMemoryData[K, V]) Update(key K, value V) (oldValue V, exist bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if item, ok := d.data[key]; ok {
		d.data[key] = adapterMemoryItem[V]{
			v: value,
			e: item.e,
		}
		return item.v, true
	}
	return oldValue, false
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func (d *adapterMemoryData[K, V]) UpdateExpire(key K, expireTime int64) (oldDuration time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if item, ok := d.data[key]; ok {
		d.data[key] = adapterMemoryItem[V]{
			v: item.v,
			e: expireTime,
		}
		return time.Duration(item.e-time.Now().UnixMilli()) * time.Millisecond
	}
	return -1
}

// Remove deletes the one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the deleted last item.
func (d *adapterMemoryData[K, V]) Remove(keys ...K) (removedKeys []K, value V) {
	d.mu.Lock()
	defer d.mu.Unlock()
	removedKeys = make([]K, 0)
	for _, key := range keys {
		item, ok := d.data[key]
		if ok {
			value = item.v
			delete(d.data, key)
			removedKeys = append(removedKeys, key)
		}
	}
	return removedKeys, value
}

// Data returns a copy of all key-value pairs in the cache as map type.
func (d *adapterMemoryData[K, V]) Data() map[K]V {
	d.mu.RLock()
	m := make(map[K]V, len(d.data))
	for k, v := range d.data {
		if !v.IsExpired() {
			m[k] = v.v
		}
	}
	d.mu.RUnlock()
	return m
}

// Keys returns all keys in the cache as slice.
func (d *adapterMemoryData[K, V]) Keys() []K {
	d.mu.RLock()
	var (
		index = 0
		keys  = make([]K, len(d.data))
	)
	for k, v := range d.data {
		if !v.IsExpired() {
			keys[index] = k
			index++
		}
	}
	d.mu.RUnlock()
	return keys
}

// Values returns all values in the cache as slice.
func (d *adapterMemoryData[K, V]) Values() []V {
	d.mu.RLock()
	var (
		index  = 0
		values = make([]V, len(d.data))
	)
	for _, v := range d.data {
		if !v.IsExpired() {
			values[index] = v.v
			index++
		}
	}
	d.mu.RUnlock()
	return values
}

// Size returns the size of the cache.
func (d *adapterMemoryData[K, V]) Size() int {
	d.mu.RLock()
	size := len(d.data)
	d.mu.RUnlock()
	return size
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
func (d *adapterMemoryData[K, V]) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = make(map[K]adapterMemoryItem[V])
}

func (d *adapterMemoryData[K, V]) Get(key K) (item adapterMemoryItem[V], ok bool) {
	d.mu.RLock()
	item, ok = d.data[key]
	d.mu.RUnlock()
	return
}

func (d *adapterMemoryData[K, V]) Set(key K, value adapterMemoryItem[V]) {
	d.mu.Lock()
	d.data[key] = value
	d.mu.Unlock()
}

// SetMap batch sets cache with key-value pairs by `data`, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (d *adapterMemoryData[K, V]) SetMap(data map[K]V, expireTime int64) {
	d.mu.Lock()
	for k, v := range data {
		d.data[k] = adapterMemoryItem[V]{
			v: v,
			e: expireTime,
		}
	}
	d.mu.Unlock()
}

func (d *adapterMemoryData[K, V]) SetWithLock(ctx context.Context, key K, value V, expireTimestamp int64) V {
	d.mu.Lock()
	defer d.mu.Unlock()
	if v, ok := d.data[key]; ok && !v.IsExpired() {
		return v.v
	}
	d.data[key] = adapterMemoryItem[V]{v: value, e: expireTimestamp}
	return value
}

func (d *adapterMemoryData[K, V]) SetWithFuncLock(ctx context.Context, key K, f Func[V], expireTimestamp int64) V {
	d.mu.Lock()
	defer d.mu.Unlock()
	value, err := f()
	if err != nil {
		return value
	}
	d.data[key] = adapterMemoryItem[V]{v: value, e: expireTimestamp}
	return value
}

func (d *adapterMemoryData[K, V]) DeleteWithDoubleCheck(key K, force ...bool) {
	d.mu.Lock()
	// Doubly check before really deleting it from cache.
	if item, ok := d.data[key]; (ok && item.IsExpired()) || (len(force) > 0 && force[0]) {
		delete(d.data, key)
	}
	d.mu.Unlock()
}
