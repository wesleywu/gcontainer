// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Timer Operations

package gtimer_test

import (
	"context"
	"testing"
	"time"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/gtimer"
	"github.com/wesleywu/gcontainer/internal/gtest"
)

func TestTimer_Add_Close(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		//fmt.Println("start", time.Now())
		timer.Add(ctx, 200*time.Millisecond, func(ctx context.Context) error {
			//fmt.Println("job1", time.Now())
			array.Add(1)
			return nil
		})
		timer.Add(ctx, 200*time.Millisecond, func(ctx context.Context) error {
			//fmt.Println("job2", time.Now())
			array.Add(1)
			return nil
		})
		timer.Add(ctx, 400*time.Millisecond, func(ctx context.Context) error {
			//fmt.Println("job3", time.Now())
			array.Add(1)
			return nil
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 5)
		timer.Close()
		time.Sleep(250 * time.Millisecond)
		fixedLength := array.Len()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), fixedLength)
	})
}

func TestTimer_Start_Stop_Close(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.Add(ctx, 1000*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			return nil
		})
		t.Assert(array.Len(), 0)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		timer.Stop()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		timer.Start()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 2)
		timer.Close()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestJob_Reset(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		job := timer.AddSingleton(ctx, 500*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			return nil
		})
		time.Sleep(300 * time.Millisecond)
		job.Reset()
		time.Sleep(300 * time.Millisecond)
		job.Reset()
		time.Sleep(300 * time.Millisecond)
		job.Reset()
		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_AddSingleton(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.AddSingleton(ctx, 200*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			time.Sleep(10 * time.Second)
			return nil
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_AddSingletonWithQuick(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New(gtimer.TimerOptions{
			Interval: 100 * time.Millisecond,
			Quick:    true,
		})
		array := g.NewArrayList[int](true)
		timer.AddSingleton(ctx, 5*time.Second, func(ctx context.Context) error {
			array.Add(1)
			time.Sleep(10 * time.Second)
			return nil
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_AddSingletonWithoutQuick(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New(gtimer.TimerOptions{
			Interval: 100 * time.Millisecond,
			Quick:    false,
		})
		array := g.NewArrayList[int](true)
		timer.AddSingleton(ctx, 5*time.Second, func(ctx context.Context) error {
			array.Add(1)
			time.Sleep(10 * time.Second)
			return nil
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 0)
	})
}

func TestTimer_AddOnce(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.AddOnce(ctx, 200*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			return nil
		})
		timer.AddOnce(ctx, 200*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			return nil
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		timer.Close()
		time.Sleep(250 * time.Millisecond)
		fixedLength := array.Len()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), fixedLength)
	})
}

func TestTimer_AddTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.AddTimes(ctx, 200*time.Millisecond, 2, func(ctx context.Context) error {
			array.Add(1)
			return nil
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_DelayAdd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.DelayAdd(ctx, 200*time.Millisecond, 200*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			return nil
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddJob(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.DelayAddEntry(ctx, 200*time.Millisecond, 200*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			return nil
		}, false, 100, gtimer.StatusReady)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddSingleton(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.DelayAddSingleton(ctx, 200*time.Millisecond, 200*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			time.Sleep(10 * time.Second)
			return nil
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddOnce(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.DelayAddOnce(ctx, 200*time.Millisecond, 200*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			return nil
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.DelayAddTimes(ctx, 200*time.Millisecond, 500*time.Millisecond, 2, func(ctx context.Context) error {
			array.Add(1)
			return nil
		})
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 2)

		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_AddLessThanInterval(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New(gtimer.TimerOptions{
			Interval: 100 * time.Millisecond,
		})
		array := g.NewArrayList[int](true)
		timer.Add(ctx, 20*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			return nil
		})
		time.Sleep(50 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(110 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(110 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_AddLeveledJob1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.DelayAdd(ctx, 1000*time.Millisecond, 1000*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			return nil
		})
		time.Sleep(1500 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_Exit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := gtimer.New()
		array := g.NewArrayList[int](true)
		timer.Add(ctx, 200*time.Millisecond, func(ctx context.Context) error {
			array.Add(1)
			gtimer.Exit()
			return nil
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}
