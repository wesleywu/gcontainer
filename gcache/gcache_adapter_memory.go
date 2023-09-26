// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"math"
	"time"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/gtimer"
	"github.com/wesleywu/gcontainer/gtype"
	"github.com/wesleywu/gcontainer/utils/empty"
)

// AdapterMemory is an adapter implements using memory.
type AdapterMemory[K comparable, V any] struct {
	// cap limits the size of the cache pool.
	// If the size of the cache exceeds the cap,
	// the cache expiration process performs according to the LRU algorithm.
	// It is 0 in default which means no limits.
	cap         int
	data        *adapterMemoryData[K, V]              // data is the underlying cache data which is stored in a hash table.
	expireTimes *adapterMemoryExpireTimes[K]          // expireTimes is the expiring key to its timestamp mapping, which is used for quick indexing and deleting.
	expireSets  *adapterMemoryExpireSets[K]           // expireSets is the expiring timestamp to its key set mapping, which is used for quick indexing and deleting.
	lru         *adapterMemoryLru[K, V]               // lru is the LRU manager, which is enabled when attribute cap > 0.
	lruGetList  *g.LinkedList[K]                      // lruGetList is the LRU history according to Get function.
	eventList   *g.LinkedList[*adapterMemoryEvent[K]] // eventList is the asynchronous event list for internal data synchronization.
	closed      *gtype.Bool                           // closed controls the cache closed or not.
}

// Internal cache item.
type adapterMemoryItem[V any] struct {
	v V     // Value.
	e int64 // Expire timestamp in milliseconds.
}

// Internal event item.
type adapterMemoryEvent[K comparable] struct {
	k K     // Key.
	e int64 // Expire time in milliseconds.
}

const (
	// defaultMaxExpire is the default expire time for no expiring items.
	// It equals to math.MaxInt64/1000000.
	defaultMaxExpire = 9223372036854
)

// NewAdapterMemory creates and returns a new memory cache object.
func NewAdapterMemory[K comparable, V any](lruCap ...int) *AdapterMemory[K, V] {
	c := &AdapterMemory[K, V]{
		data:        newAdapterMemoryData[K, V](),
		lruGetList:  g.NewLinkedList[K](true),
		expireTimes: newAdapterMemoryExpireTimes[K](),
		expireSets:  newAdapterMemoryExpireSets[K](),
		eventList:   g.NewLinkedList[*adapterMemoryEvent[K]](true),
		closed:      gtype.NewBool(),
	}
	if len(lruCap) > 0 {
		c.cap = lruCap[0]
		c.lru = newMemCacheLru[K, V](c)
	}
	return c
}

// Set sets cache with `key`-`value` pair, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory[K, V]) Set(ctx context.Context, key K, value V, duration time.Duration) error {
	expireTime := c.getInternalExpire(duration)
	c.data.Set(key, adapterMemoryItem[V]{
		v: value,
		e: expireTime,
	})
	c.eventList.PushBack(&adapterMemoryEvent[K]{
		k: key,
		e: expireTime,
	})
	return nil
}

// SetMap batch sets cache with key-value pairs by `data` map, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory[K, V]) SetMap(ctx context.Context, data map[K]V, duration time.Duration) error {
	expireTime := c.getInternalExpire(duration)
	c.data.SetMap(data, expireTime)
	for k := range data {
		c.eventList.PushBack(&adapterMemoryEvent[K]{
			k: k,
			e: expireTime,
		})
	}
	return nil
}

// SetIfNotExist sets cache with `key`-`value` pair which is expired after `duration`
// if `key` does not exist in the cache. It returns true the `key` does not exist in the
// cache, and it sets `value` successfully to the cache, or else it returns false.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory[K, V]) SetIfNotExist(ctx context.Context, key K, value V, duration time.Duration) (bool, error) {
	isContained, err := c.Contains(ctx, key)
	if err != nil {
		return false, err
	}
	if !isContained {
		_, err := c.doSetWithLockCheck(ctx, key, value, duration)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// SetIfNotExistFunc sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// The parameter `value` can be type of `func() interface{}`, but it does nothing if its
// result is nil.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
func (c *AdapterMemory[K, V]) SetIfNotExistFunc(ctx context.Context, key K, f Func[V], duration time.Duration) (bool, error) {
	isContained, err := c.Contains(ctx, key)
	if err != nil {
		return false, err
	}
	if !isContained {
		value, err := f()
		if err != nil {
			return false, err
		}
		_, err = c.doSetWithLockCheck(ctx, key, value, duration)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// SetIfNotExistFuncLock sets `key` with result of function `f` and returns true
// if `key` does not exist in the cache, or else it does nothing and returns false if `key` already exists.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil.
//
// Note that it differs from function `SetIfNotExistFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func (c *AdapterMemory[K, V]) SetIfNotExistFuncLock(ctx context.Context, key K, f Func[V], duration time.Duration) (bool, error) {
	isContained, err := c.Contains(ctx, key)
	if err != nil {
		return false, err
	}
	if !isContained {
		_, err = c.doSetWithFuncLockCheck(ctx, key, f, duration)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Get retrieves and returns the associated value of given `key`.
// It returns nil if it does not exist, or its value is nil, or it's expired.
// If you would like to check if the `key` exists in the cache, it's better using function Contains.
func (c *AdapterMemory[K, V]) Get(ctx context.Context, key K) (v V, ok bool, err error) {
	item, ok := c.data.Get(key)
	if ok && !item.IsExpired() {
		// Adding to LRU history if LRU feature is enabled.
		if c.cap > 0 {
			c.lruGetList.PushBack(key)
		}
		v = item.v
		return v, true, nil
	}
	return v, false, nil
}

// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and
// returns `value` if `key` does not exist in the cache. The key-value pair expires
// after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (c *AdapterMemory[K, V]) GetOrSet(ctx context.Context, key K, value V, duration time.Duration) (V, error) {
	v, ok, err := c.Get(ctx, key)
	if err != nil {
		return v, err
	}
	if !ok {
		return c.doSetWithLockCheck(ctx, key, value, duration)
	}
	return v, nil
}

// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
func (c *AdapterMemory[K, V]) GetOrSetFunc(ctx context.Context, key K, f Func[V], duration time.Duration) (v V, err error) {
	var ok bool
	v, ok, err = c.Get(ctx, key)
	if err != nil {
		return v, err
	}
	if !ok {
		v, err = f()
		if err != nil {
			return v, err
		}
		if empty.IsNil(v) {
			return v, nil
		}
		return c.doSetWithLockCheck(ctx, key, v, duration)
	}
	return v, nil
}

// GetOrSetFuncLock retrieves and returns the value of `key`, or sets `key` with result of
// function `f` and returns its result if `key` does not exist in the cache. The key-value
// pair expires after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the `key` if `duration` < 0 or given `value` is nil, but it does nothing
// if `value` is a function and the function result is nil.
//
// Note that it differs from function `GetOrSetFunc` is that the function `f` is executed within
// writing mutex lock for concurrent safety purpose.
func (c *AdapterMemory[K, V]) GetOrSetFuncLock(ctx context.Context, key K, f Func[V], duration time.Duration) (V, error) {
	v, ok, err := c.Get(ctx, key)
	if err != nil {
		return v, err
	}
	if !ok {
		return c.doSetWithFuncLockCheck(ctx, key, f, duration)
	}
	return v, nil
}

// Contains checks and returns true if `key` exists in the cache, or else returns false.
func (c *AdapterMemory[K, V]) Contains(ctx context.Context, key K) (ok bool, err error) {
	_, ok, err = c.Get(ctx, key)
	return ok, err
}

// GetExpire retrieves and returns the expiration of `key` in the cache.
//
// Note that,
// It returns 0 if the `key` does not expire.
// It returns -1 if the `key` does not exist in the cache.
func (c *AdapterMemory[K, V]) GetExpire(ctx context.Context, key K) (time.Duration, error) {
	if item, ok := c.data.Get(key); ok {
		return time.Duration(item.e-time.Now().UnixMilli()) * time.Millisecond, nil
	}
	return -1, nil
}

// Remove deletes one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the last deleted item.
func (c *AdapterMemory[K, V]) Remove(ctx context.Context, keys ...K) (v V, err error) {
	var removedKeys []K
	removedKeys, value := c.data.Remove(keys...)
	for _, key := range removedKeys {
		c.eventList.PushBack(&adapterMemoryEvent[K]{
			k: key,
			e: time.Now().UnixMilli() - 1000000,
		})
	}
	return value, nil
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func (c *AdapterMemory[K, V]) Update(ctx context.Context, key K, value V) (oldValue V, exist bool, err error) {
	v, exist := c.data.Update(key, value)
	return v, exist, nil
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func (c *AdapterMemory[K, V]) UpdateExpire(ctx context.Context, key K, duration time.Duration) (oldDuration time.Duration, err error) {
	newExpireTime := c.getInternalExpire(duration)
	oldDuration = c.data.UpdateExpire(key, newExpireTime)
	if oldDuration != -1 {
		c.eventList.PushBack(&adapterMemoryEvent[K]{
			k: key,
			e: newExpireTime,
		})
	}
	return
}

// Size returns the size of the cache.
func (c *AdapterMemory[K, V]) Size(ctx context.Context) (size int, err error) {
	return c.data.Size(), nil
}

// Data returns a copy of all key-value pairs in the cache as map type.
func (c *AdapterMemory[K, V]) Data(ctx context.Context) (map[K]V, error) {
	return c.data.Data(), nil
}

// Keys returns all keys in the cache as slice.
func (c *AdapterMemory[K, V]) Keys(ctx context.Context) ([]K, error) {
	return c.data.Keys(), nil
}

// Values returns all values in the cache as slice.
func (c *AdapterMemory[K, V]) Values(ctx context.Context) ([]V, error) {
	return c.data.Values(), nil
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
func (c *AdapterMemory[K, V]) Clear(ctx context.Context) error {
	c.data.Clear()
	return nil
}

// Close closes the cache.
func (c *AdapterMemory[K, V]) Close(ctx context.Context) error {
	if c.cap > 0 {
		c.lru.Close()
	}
	c.closed.Set(true)
	return nil
}

// doSetWithFuncLockCheck sets cache with `key`-`value` pair if `key` does not exist in the
// cache, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// The parameter `value` can be type of <func() interface{}>, but it does nothing if the
// function result is nil.
//
// It doubly checks the `key` whether exists in the cache using mutex writing lock
// before setting it to the cache.
func (c *AdapterMemory[K, V]) doSetWithFuncLockCheck(ctx context.Context, key K, f Func[V], duration time.Duration) (v V, err error) {
	expireTimestamp := c.getInternalExpire(duration)
	v = c.data.SetWithFuncLock(ctx, key, f, expireTimestamp)
	c.eventList.PushBack(&adapterMemoryEvent[K]{k: key, e: expireTimestamp})
	return v, nil
}

// doSetWithLockCheck sets cache with `key`-`value` pair if `key` does not exist in the
// cache, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// The parameter `value` can be type of <func() interface{}>, but it does nothing if the
// function result is nil.
//
// It doubly checks the `key` whether exists in the cache using mutex writing lock
// before setting it to the cache.
func (c *AdapterMemory[K, V]) doSetWithLockCheck(ctx context.Context, key K, value V, duration time.Duration) (v V, err error) {
	expireTimestamp := c.getInternalExpire(duration)
	v = c.data.SetWithLock(ctx, key, value, expireTimestamp)
	c.eventList.PushBack(&adapterMemoryEvent[K]{k: key, e: expireTimestamp})
	return v, nil
}

// getInternalExpire converts and returns the expiration time with given expired duration in milliseconds.
func (c *AdapterMemory[K, V]) getInternalExpire(duration time.Duration) int64 {
	if duration == 0 {
		return defaultMaxExpire
	}
	return time.Now().UnixMilli() + duration.Nanoseconds()/1000000
}

// makeExpireKey groups the `expire` in milliseconds to its according seconds.
func (c *AdapterMemory[K, V]) makeExpireKey(expire int64) int64 {
	return int64(math.Ceil(float64(expire/1000)+1) * 1000)
}

// syncEventAndClearExpired does the asynchronous task loop:
// 1. Asynchronously process the data in the event list,
// and synchronize the results to the `expireTimes` and `expireSets` properties.
// 2. Clean up the expired key-value pair data.
func (c *AdapterMemory[K, V]) syncEventAndClearExpired(ctx context.Context) {
	if c.closed.Val() {
		gtimer.Exit()
		return
	}
	var (
		oldExpireTime int64
		newExpireTime int64
	)
	// ========================
	// Data Synchronization.
	// ========================
	for {
		event, ok := c.eventList.PopFront()
		if !ok {
			break
		}
		// Fetching the old expire set.
		oldExpireTime = c.expireTimes.Get(event.k)
		// Calculating the new expiration time set.
		newExpireTime = c.makeExpireKey(event.e)
		if newExpireTime != oldExpireTime {
			c.expireSets.GetOrNew(newExpireTime).Add(event.k)
			if oldExpireTime != 0 {
				c.expireSets.GetOrNew(oldExpireTime).Remove(event.k)
			}
			// Updating the expired time for <event.k>.
			c.expireTimes.Set(event.k, newExpireTime)
		}
		// Adding the key the LRU history by writing operations.
		if c.cap > 0 {
			c.lru.Push(event.k)
		}
	}
	// Processing expired keys from LRU.
	if c.cap > 0 {
		if c.lruGetList.Size() > 0 {
			for {
				if v, ok := c.lruGetList.PopFront(); ok {
					c.lru.Push(v)
				} else {
					break
				}
			}
		}
		c.lru.SyncAndClear(ctx)
	}
	// ========================
	// Data Cleaning up.
	// ========================
	var (
		expireSet g.Set[K]
		ek        = c.makeExpireKey(time.Now().UnixMilli())
		eks       = []int64{ek - 1000, ek - 2000, ek - 3000, ek - 4000, ek - 5000}
	)
	for _, expireTime := range eks {
		if expireSet = c.expireSets.Get(expireTime); expireSet != nil {
			// Iterating the set to delete all keys in it.
			expireSet.ForEach(func(key K) bool {
				c.clearByKey(key)
				return true
			})
			// Deleting the set after all of its keys are deleted.
			c.expireSets.Delete(expireTime)
		}
	}
}

// clearByKey deletes the key-value pair with given `key`.
// The parameter `force` specifies whether doing this deleting forcibly.
func (c *AdapterMemory[K, V]) clearByKey(key K, force ...bool) {
	// Doubly check before really deleting it from cache.
	c.data.DeleteWithDoubleCheck(key, force...)

	// Deleting its expiration time from `expireTimes`.
	c.expireTimes.Delete(key)

	// Deleting it from LRU.
	if c.cap > 0 {
		c.lru.Remove(key)
	}
}
