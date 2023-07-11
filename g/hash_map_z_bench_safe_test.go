// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package g_test

import (
	"strconv"
	"testing"

	"github.com/wesleywu/gcontainer/g"
)

var intIntMap = g.NewHashMap[int, int](true)

var intAnyMap = g.NewHashMap[int, any](true)

var intStrMap = g.NewHashMap[int, string](true)

var strIntMap = g.NewHashMap[string, int](true)

var strAnyMap = g.NewHashMap[string, any](true)

var strStrMap = g.NewHashMap[string, string](true)

func Benchmark_IntIntMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intIntMap.Put(i, i)
			i++
		}
	})
}

func Benchmark_IntAnyMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intAnyMap.Put(i, i)
			i++
		}
	})
}

func Benchmark_IntStrMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intStrMap.Put(i, "123456789")
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrIntMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strIntMap.Put(strconv.Itoa(i), i)
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrAnyMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strAnyMap.Put(strconv.Itoa(i), i)
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrStrMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strStrMap.Put(strconv.Itoa(i), "123456789")
			i++
		}
	})
}

func Benchmark_IntIntMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intIntMap.Get(i)
			i++
		}
	})
}

func Benchmark_IntAnyMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intAnyMap.Get(i)
			i++
		}
	})
}

func Benchmark_IntStrMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intStrMap.Get(i)
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrIntMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strIntMap.Get(strconv.Itoa(i))
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrAnyMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strAnyMap.Get(strconv.Itoa(i))
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrStrMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strStrMap.Get(strconv.Itoa(i))
			i++
		}
	})
}
