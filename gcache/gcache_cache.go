// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"time"

	"github.com/wesleywu/gcontainer/gtimer"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

// Cache struct.
type Cache[K comparable, V any] struct {
	Adapter[K, V]
}

// New creates and returns a new cache object using default memory adapter.
// Note that the LRU feature is only available using memory adapter.
func New[K comparable, V any](lruCap ...int) *Cache[K, V] {
	memAdapter := NewAdapterMemory[K, V](lruCap...)
	c := &Cache[K, V]{
		Adapter: memAdapter,
	}
	// Here may be a "timer leak" if adapter is manually changed from memory adapter.
	// Do not worry about this, as adapter is less changed, and it does nothing if it's not used.
	gtimer.AddSingleton(context.Background(), time.Second, memAdapter.syncEventAndClearExpired)
	return c
}

// NewWithAdapter creates and returns a Cache object with given Adapter implements.
func NewWithAdapter[K comparable, V any](adapter Adapter[K, V]) *Cache[K, V] {
	return &Cache[K, V]{
		Adapter: adapter,
	}
}

// SetAdapter changes the adapter for this cache.
// Be very note that, this setting function is not concurrent-safe, which means you should not call
// this setting function concurrently in multiple goroutines.
func (c *Cache[K, V]) SetAdapter(adapter Adapter[K, V]) {
	c.Adapter = adapter
}

// GetAdapter returns the adapter that is set in current Cache.
func (c *Cache[K, V]) GetAdapter() Adapter[K, V] {
	return c.Adapter
}

// Removes deletes `keys` in the cache.
func (c *Cache[K, V]) Removes(ctx context.Context, keys []K) error {
	_, err := c.Remove(ctx, keys...)
	return err
}

// KeyStrings returns all keys in the cache as string slice.
func (c *Cache[K, V]) KeyStrings(ctx context.Context) ([]string, error) {
	keys, err := c.Keys(ctx)
	if err != nil {
		return nil, err
	}
	return gconv.Strings(keys), nil
}
