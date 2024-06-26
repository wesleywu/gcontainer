// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcache_test

import (
	"context"
	"log"
	"math"
	"testing"
	"time"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/gcache"
	"github.com/wesleywu/gcontainer/internal/gtest"
)

var (
	ctx = context.Background()
)

func TestCache_GCache_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		t.AssertNil(gcache.Set(ctx, 1, 11, 0))
		defer gcache.Remove(ctx, []any{1, 2, 3}...)
		v, ok, _ := gcache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
		b, _ := gcache.Contains(ctx, 1)
		t.Assert(b, true)
	})
}

func TestCache_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		c := gcache.New[int, int]()
		defer c.Close(ctx)
		t.Assert(c.Set(ctx, 1, 11, 0), nil)
		v, ok, _ := c.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
		b, _ := c.Contains(ctx, 1)
		t.Assert(b, true)
	})
}

func TestCache_Set_Expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[int, int]()
		t.Assert(cache.Set(ctx, 2, 22, 100*time.Millisecond), nil)
		v, ok, _ := cache.Get(ctx, 2)
		t.Assert(v, 22)
		t.Assert(ok, true)
		time.Sleep(200 * time.Millisecond)
		v, ok, _ = cache.Get(ctx, 2)
		t.Assert(ok, false)
		time.Sleep(3 * time.Second)
		n, _ := cache.Size(ctx)
		t.Assert(n, 0)
		t.Assert(cache.Close(ctx), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[int, int]()
		t.Assert(cache.Set(ctx, 1, 11, 100*time.Millisecond), nil)
		v, ok, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
		time.Sleep(200 * time.Millisecond)
		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(ok, false)
	})
}

func TestCache_Update(t *testing.T) {
	// gcache
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		key := "12345678"
		t.AssertNil(gcache.Set(ctx, key, 11, 3*time.Second))
		expire1, _ := gcache.GetExpire(ctx, key)
		oldValue, exist, err := gcache.Update(ctx, key, 12)
		t.AssertNil(err)
		t.Assert(oldValue, 11)
		t.Assert(exist, true)

		expire2, _ := gcache.GetExpire(ctx, key)
		v, ok, _ := gcache.Get(ctx, key)
		t.Assert(v, 12)
		t.Assert(ok, true)
		t.Assert(math.Ceil(expire1.Seconds()), math.Ceil(expire2.Seconds()))
	})
	// gcache.Cache
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[int, int]()
		t.AssertNil(cache.Set(ctx, 1, 11, 3*time.Second))

		oldValue, exist, err := cache.Update(ctx, 1, 12)
		t.AssertNil(err)
		t.Assert(oldValue, 11)
		t.Assert(exist, true)

		expire1, _ := cache.GetExpire(ctx, 1)
		expire2, _ := cache.GetExpire(ctx, 1)
		v, ok, _ := cache.Get(ctx, 1)
		t.Assert(v, 12)
		t.Assert(ok, true)
		t.Assert(math.Ceil(expire1.Seconds()), math.Ceil(expire2.Seconds()))
	})
}

func TestCache_UpdateExpire(t *testing.T) {
	// gcache
	gtest.C(t, func(t *gtest.T) {
		key := "12345678"
		t.AssertNil(gcache.Set(ctx, key, 11, 3*time.Second))
		defer gcache.Remove(ctx, key)
		oldExpire, _ := gcache.GetExpire(ctx, key)
		newExpire := 10 * time.Second
		oldExpire2, err := gcache.UpdateExpire(ctx, key, newExpire)
		t.AssertNil(err)
		t.AssertIN(oldExpire2, []any{oldExpire, `2.999s`})

		e, _ := gcache.GetExpire(ctx, key)
		t.AssertNE(e, oldExpire)
		e, _ = gcache.GetExpire(ctx, key)
		t.Assert(math.Ceil(e.Seconds()), 10)
	})
	// gcache.Cache
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New[any, any]()
		t.AssertNil(cache.Set(ctx, 1, 11, 3*time.Second))
		oldExpire, _ := cache.GetExpire(ctx, 1)
		newExpire := 10 * time.Second
		oldExpire2, err := cache.UpdateExpire(ctx, 1, newExpire)
		t.AssertNil(err)
		t.AssertIN(oldExpire2, []any{oldExpire, `2.999s`})

		e, _ := cache.GetExpire(ctx, 1)
		t.AssertNE(e, oldExpire)

		e, _ = cache.GetExpire(ctx, 1)
		t.Assert(math.Ceil(e.Seconds()), 10)
	})
}

func TestCache_Keys_Values(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := gcache.New[any, any]()
		for i := 0; i < 10; i++ {
			t.Assert(c.Set(ctx, i, i*10, 0), nil)
		}
		var (
			keys, _   = c.Keys(ctx)
			values, _ = c.Values(ctx)
		)
		t.Assert(len(keys), 10)
		t.Assert(len(values), 10)
		t.AssertIN(0, keys)
		t.AssertIN(90, values)
	})
}

func TestCache_LRU(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[int, int](2)
		for i := 0; i < 10; i++ {
			t.AssertNil(cache.Set(ctx, i, i, 0))
		}
		n, _ := cache.Size(ctx)
		t.Assert(n, 10)
		v, ok, _ := cache.Get(ctx, 6)
		t.Assert(v, 6)
		t.Assert(ok, true)
		time.Sleep(4 * time.Second)
		log.Printf(`items after lru: %+v`, cache.MustData(ctx))
		n, _ = cache.Size(ctx)
		t.Assert(n, 2)
		v, ok, _ = cache.Get(ctx, 6)
		t.Assert(v, 6)
		t.Assert(ok, true)
		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(ok, false)
		t.Assert(cache.Close(ctx), nil)
	})
}

func TestCache_LRU_expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[int, any](2)
		t.Assert(cache.Set(ctx, 1, nil, 1000), nil)
		n, _ := cache.Size(ctx)
		t.Assert(n, 1)
		v, ok, _ := cache.Get(ctx, 1)
		t.Assert(ok, true)
		t.Assert(v, nil)
	})
}

func TestCache_SetIfNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New[any, any]()
		ok, err := cache.SetIfNotExist(ctx, 1, 11, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, ok, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		ok, err = cache.SetIfNotExist(ctx, 1, 22, 0)
		t.AssertNil(err)
		t.Assert(ok, false)

		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		ok, err = cache.SetIfNotExist(ctx, 2, 22, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, ok, _ = cache.Get(ctx, 2)
		t.Assert(v, 22)
		t.Assert(ok, true)

		gcache.Remove(ctx, []any{1, 2, 3}...)
		ok, err = gcache.SetIfNotExist(ctx, 1, 11, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, ok, _ = gcache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		ok, err = gcache.SetIfNotExist(ctx, 1, 22, 0)
		t.AssertNil(err)
		t.Assert(ok, false)

		v, ok, _ = gcache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
	})
}

func TestCache_SetIfNotExistFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[int, int]()
		exist, err := cache.SetIfNotExistFunc(ctx, 1, func() (value int, err error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, true)

		v, ok, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		exist, err = cache.SetIfNotExistFunc(ctx, 1, func() (value int, err error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, false)

		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
	})
	gtest.C(t, func(t *gtest.T) {
		gcache.Remove(ctx, []any{1, 2, 3}...)

		ok, err := gcache.SetIfNotExistFunc(ctx, 1, func() (value interface{}, err error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, ok, _ := gcache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		ok, err = gcache.SetIfNotExistFunc(ctx, 1, func() (value interface{}, err error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(ok, false)

		v, ok, _ = gcache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
	})
}

func TestCache_SetIfNotExistFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[any, any]()
		exist, err := cache.SetIfNotExistFuncLock(ctx, 1, func() (value interface{}, err error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, true)

		v, ok, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		exist, err = cache.SetIfNotExistFuncLock(ctx, 1, func() (value interface{}, err error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, false)

		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
	})
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		gcache.Remove(ctx, []any{1, 2, 3}...)

		exist, err := gcache.SetIfNotExistFuncLock(ctx, 1, func() (value interface{}, err error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, true)

		v, ok, _ := gcache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		exist, err = gcache.SetIfNotExistFuncLock(ctx, 1, func() (value interface{}, err error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, false)

		v, ok, _ = gcache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
	})
}

func TestCache_SetMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[int, int]()
		t.AssertNil(cache.SetMap(ctx, map[int]int{1: 11, 2: 22}, 0))
		v, ok, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)

		gcache.Remove(ctx, []any{1, 2, 3}...)
		t.AssertNil(gcache.SetMap(ctx, map[any]any{1: 11, 2: 22}, 0))
		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
	})
}

func TestCache_GetOrSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[int, int]()
		value, err := cache.GetOrSet(ctx, 1, 11, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, ok, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
		value, err = cache.GetOrSet(ctx, 1, 111, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
	})

	gtest.C(t, func(t *gtest.T) {
		var ok bool
		gcache.Remove(ctx, []any{1, 2, 3}...)
		value, err := gcache.GetOrSet(ctx, 1, 11, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, ok, err := gcache.Get(ctx, 1)
		t.AssertNil(err)
		t.Assert(v, 11)
		t.Assert(ok, true)

		value, err = gcache.GetOrSet(ctx, 1, 111, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, ok, err = gcache.Get(ctx, 1)
		t.AssertNil(err)
		t.Assert(v, 11)
		t.Assert(ok, true)
	})
}

func TestCache_GetOrSetFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[any, any]()
		cache.GetOrSetFunc(ctx, 1, func() (value interface{}, err error) {
			return 11, nil
		}, 0)
		v, ok, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		cache.GetOrSetFunc(ctx, 1, func() (value interface{}, err error) {
			return 111, nil
		}, 0)
		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		gcache.Remove(ctx, []any{1, 2, 3}...)

		gcache.GetOrSetFunc(ctx, 1, func() (value interface{}, err error) {
			return 11, nil
		}, 0)
		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		gcache.GetOrSetFunc(ctx, 1, func() (value interface{}, err error) {
			return 111, nil
		}, 0)
		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
	})
}

func TestCache_GetOrSetFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool
		cache := gcache.New[int, int]()
		cache.GetOrSetFuncLock(ctx, 1, func() (value int, err error) {
			return 11, nil
		}, 0)
		v, ok, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		cache.GetOrSetFuncLock(ctx, 1, func() (value int, err error) {
			return 111, nil
		}, 0)
		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		gcache.Remove(ctx, []any{1, 2, 3}...)
		gcache.GetOrSetFuncLock(ctx, 1, func() (value any, err error) {
			return 11, nil
		}, 0)
		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)

		gcache.GetOrSetFuncLock(ctx, 1, func() (value any, err error) {
			return 111, nil
		}, 0)
		v, ok, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
		t.Assert(ok, true)
	})
}

func TestCache_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New[int, int]()
		cache.SetMap(ctx, map[int]int{1: 11, 2: 22}, 0)
		cache.Clear(ctx)
		n, _ := cache.Size(ctx)
		t.Assert(n, 0)
	})
}

//
//func TestCache_SetConcurrency(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		cache := gcache.New[any, any]()
//		pool := gpool.New[any](4, )
//		go func() {
//			for {
//				pool.(ctx, func() {
//					cache.SetIfNotExist(ctx, 1, 11, 10)
//				})
//			}
//		}()
//		select {
//		case <-time.After(2 * time.Second):
//			// t.Log("first part end")
//		}
//
//		go func() {
//			for {
//				pool.Add(ctx, func() {
//					cache.SetIfNotExist(ctx, 1, nil, 10)
//				})
//			}
//		}()
//		select {
//		case <-time.After(2 * time.Second):
//			// t.Log("second part end")
//		}
//	})
//}

func TestCache_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		{
			var ok bool
			cache := gcache.New[int, int]()
			cache.SetMap(ctx, map[int]int{1: 11, 2: 22}, 0)
			b, _ := cache.Contains(ctx, 1)
			t.Assert(b, true)
			v, ok, _ := cache.Get(ctx, 1)
			t.Assert(v, 11)
			t.Assert(ok, true)
			data, _ := cache.Data(ctx)
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], 0)
			n, _ := cache.Size(ctx)
			t.Assert(n, 2)
			keys, _ := cache.Keys(ctx)
			t.Assert(g.NewHashSetFrom[int]([]int{1, 2}).Equals(g.NewHashSetFrom[int](keys)), true)
			keyStrs, _ := cache.KeyStrings(ctx)
			t.Assert(g.NewHashSetFrom[string]([]string{"1", "2"}).Equals(g.NewHashSetFrom[string](keyStrs)), true)
			values, _ := cache.Values(ctx)
			t.Assert(g.NewHashSetFrom[int]([]int{11, 22}).Equals(g.NewHashSetFrom[int](values)), true)
			removeData1, _ := cache.Remove(ctx, 1)
			t.Assert(removeData1, 11)
			n, _ = cache.Size(ctx)
			t.Assert(n, 1)

			cache.Remove(ctx, 2)
			n, _ = cache.Size(ctx)
			t.Assert(n, 0)
		}

		gcache.Remove(ctx, []any{1, 2, 3}...)
		{
			var ok bool
			gcache.SetMap(ctx, map[any]any{1: 11, 2: 22}, 0)
			b, _ := gcache.Contains(ctx, 1)
			t.Assert(b, true)
			v, ok, _ := gcache.Get(ctx, 1)
			t.Assert(v, 11)
			t.Assert(ok, true)
			data, _ := gcache.Data(ctx)
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			n, _ := gcache.Size(ctx)
			t.Assert(n, 2)
			keys, _ := gcache.Keys(ctx)
			t.Assert(g.NewHashSetFrom[any]([]any{1, 2}).Equals(g.NewHashSetFrom[any](keys)), true)
			keyStrs, _ := gcache.KeyStrings(ctx)
			t.Assert(g.NewHashSetFrom[string]([]string{"1", "2"}).Equals(g.NewHashSetFrom[string](keyStrs)), true)
			values, _ := gcache.Values(ctx)
			t.Assert(g.NewHashSetFrom[any]([]any{11, 22}).Equals(g.NewHashSetFrom[any](values)), true)
			removeData1, _ := gcache.Remove(ctx, 1)
			t.Assert(removeData1, 11)
			n, _ = gcache.Size(ctx)
			t.Assert(n, 1)
			gcache.Remove(ctx, 2)
			n, _ = gcache.Size(ctx)
			t.Assert(n, 0)
		}
	})
}

func TestCache_Removes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New[int, int]()
		t.AssertNil(cache.Set(ctx, 1, 11, 0))
		t.AssertNil(cache.Set(ctx, 2, 22, 0))
		t.AssertNil(cache.Set(ctx, 3, 33, 0))
		t.AssertNil(cache.Removes(ctx, []int{2, 3}))

		ok, err := cache.Contains(ctx, 1)
		t.AssertNil(err)
		t.Assert(ok, true)

		ok, err = cache.Contains(ctx, 2)
		t.AssertNil(err)
		t.Assert(ok, false)
	})

	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(gcache.Set(ctx, 1, 11, 0))
		t.AssertNil(gcache.Set(ctx, 2, 22, 0))
		t.AssertNil(gcache.Set(ctx, 3, 33, 0))
		t.AssertNil(gcache.Removes(ctx, []any{2, 3}))

		ok, err := gcache.Contains(ctx, 1)
		t.AssertNil(err)
		t.Assert(ok, true)

		ok, err = gcache.Contains(ctx, 2)
		t.AssertNil(err)
		t.Assert(ok, false)
	})
}

func TestCache_Basic_Must(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer gcache.Remove(ctx, []any{1, 2, 3, 4}...)

		t.AssertNil(gcache.Set(ctx, 1, 11, 0))
		v := gcache.MustGet(ctx, 1)
		t.Assert(v, 11)
		gcache.MustGetOrSet(ctx, 2, 22, 0)
		v = gcache.MustGet(ctx, 2)
		t.Assert(v, 22)

		gcache.MustGetOrSetFunc(ctx, 3, func() (value interface{}, err error) {
			return 33, nil
		}, 0)
		v = gcache.MustGet(ctx, 3)
		t.Assert(v, 33)

		gcache.GetOrSetFuncLock(ctx, 4, func() (value interface{}, err error) {
			return 44, nil
		}, 0)
		v = gcache.MustGet(ctx, 4)
		t.Assert(v, 44)

		t.Assert(gcache.MustContains(ctx, 1), true)

		t.AssertNil(gcache.Set(ctx, 1, 11, 3*time.Second))
		expire := gcache.MustGetExpire(ctx, 1)
		t.AssertGE(expire, 0)

		n := gcache.MustSize(ctx)
		t.Assert(n, 4)

		data := gcache.MustData(ctx)
		t.Assert(len(data), 4)

		keys := gcache.MustKeys(ctx)
		t.Assert(len(keys), 4)

		keyStrings := gcache.MustKeyStrings(ctx)
		t.Assert(len(keyStrings), 4)

		values := gcache.MustValues(ctx)
		t.Assert(len(values), 4)
	})
}

func TestCache_NewWithAdapter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.NewWithAdapter[any, any](gcache.NewAdapterMemory[any, any]())
		t.AssertNE(cache, nil)
	})
}
