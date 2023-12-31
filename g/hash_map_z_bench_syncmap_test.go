// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package g_test

import (
	"sync"
	"testing"

	"github.com/wesleywu/gcontainer/g"
)

var gm = g.NewHashMap[int, int](true)

var sm = sync.Map{}

func Benchmark_GMapSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			gm.Put(i, i)
			i++
		}
	})
}

func Benchmark_SyncMapSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Store(i, i)
			i++
		}
	})
}

func Benchmark_GMapGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			gm.Get(i)
			i++
		}
	})
}

func Benchmark_SyncMapGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Load(i)
			i++
		}
	})
}

func Benchmark_GMapRemove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			gm.Remove(i)
			i++
		}
	})
}

func Benchmark_SyncMapRmove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Delete(i)
			i++
		}
	})
}
