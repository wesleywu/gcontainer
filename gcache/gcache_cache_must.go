// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"time"
)

// MustGet acts like Get, but it panics if any error occurs.
func (c *Cache[K, V]) MustGet(ctx context.Context, key K) V {
	v, _, err := c.Get(ctx, key)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetOrSet acts like GetOrSet, but it panics if any error occurs.
func (c *Cache[K, V]) MustGetOrSet(ctx context.Context, key K, value V, duration time.Duration) V {
	v, err := c.GetOrSet(ctx, key, value, duration)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetOrSetFunc acts like GetOrSetFunc, but it panics if any error occurs.
func (c *Cache[K, V]) MustGetOrSetFunc(ctx context.Context, key K, f Func[V], duration time.Duration) V {
	v, err := c.GetOrSetFunc(ctx, key, f, duration)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetOrSetFuncLock acts like GetOrSetFuncLock, but it panics if any error occurs.
func (c *Cache[K, V]) MustGetOrSetFuncLock(ctx context.Context, key K, f Func[V], duration time.Duration) V {
	v, err := c.GetOrSetFuncLock(ctx, key, f, duration)
	if err != nil {
		panic(err)
	}
	return v
}

// MustContains acts like Contains, but it panics if any error occurs.
func (c *Cache[K, V]) MustContains(ctx context.Context, key K) bool {
	v, err := c.Contains(ctx, key)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetExpire acts like GetExpire, but it panics if any error occurs.
func (c *Cache[K, V]) MustGetExpire(ctx context.Context, key K) time.Duration {
	v, err := c.GetExpire(ctx, key)
	if err != nil {
		panic(err)
	}
	return v
}

// MustSize acts like Size, but it panics if any error occurs.
func (c *Cache[K, V]) MustSize(ctx context.Context) int {
	v, err := c.Size(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// MustData acts like Data, but it panics if any error occurs.
func (c *Cache[K, V]) MustData(ctx context.Context) map[K]V {
	v, err := c.Data(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// MustKeys acts like Keys, but it panics if any error occurs.
func (c *Cache[K, V]) MustKeys(ctx context.Context) []K {
	v, err := c.Keys(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// MustKeyStrings acts like KeyStrings, but it panics if any error occurs.
func (c *Cache[K, V]) MustKeyStrings(ctx context.Context) []string {
	v, err := c.KeyStrings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// MustValues acts like Values, but it panics if any error occurs.
func (c *Cache[K, V]) MustValues(ctx context.Context) []V {
	v, err := c.Values(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
