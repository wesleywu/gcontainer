// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gqueue_test

import (
	"testing"
	"time"

	"github.com/wesleywu/gcontainer/gqueue"
	"github.com/wesleywu/gcontainer/internal/gtest"
)

func TestBlockingQueue_Len(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			maxNum   = 100
			maxTries = 100
		)
		for n := 10; n < maxTries; n++ {
			q1 := gqueue.New[int](maxNum)
			for i := 0; i < maxNum; i++ {
				q1.Push(i)
			}
			t.Assert(q1.Len(), maxNum)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			maxNum   = 100
			maxTries = 100
		)
		for n := 10; n < maxTries; n++ {
			q1 := gqueue.New[int]()
			for i := 0; i < maxNum; i++ {
				q1.Push(i)
			}
			t.AssertLE(q1.Len(), maxNum)
		}
	})
}

func TestBlockingQueue_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New[int]()
		for i := 0; i < 100; i++ {
			q.Push(i)
		}
		t.Assert(q.MustPop(), 0)
		t.Assert(q.MustPop(), 1)
	})
}

func TestBlockingQueue_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q1 := gqueue.New[int]()
		q1.Push(1)
		q1.Push(2)
		q1.Push(3)
		q1.Push(4)
		i1 := q1.MustPop()
		t.Assert(i1, 1)
	})
}

func TestBlockingQueue_Close(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q1 := gqueue.New[int]()
		q1.Push(1)
		q1.Push(2)
		time.Sleep(time.Millisecond)
		t.Assert(q1.Len(), 2)
		q1.Close()
	})
	gtest.C(t, func(t *gtest.T) {
		q1 := gqueue.New[int](2)
		q1.Push(1)
		q1.Push(2)
		time.Sleep(time.Millisecond)
		t.Assert(q1.Len(), 2)
		q1.Close()
	})
}

func Test_Issue2509(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New[int]()
		q.Push(1)
		q.Push(2)
		q.Push(3)
		t.AssertLE(q.Len(), 3)
		t.Assert(<-q.C, 1)
		t.AssertLE(q.Len(), 2)
		t.Assert(<-q.C, 2)
		t.AssertLE(q.Len(), 1)
		t.Assert(<-q.C, 3)
		t.Assert(q.Len(), 0)
	})
}
