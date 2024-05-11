package gqueue_test

import (
	"sync"
	"testing"

	"github.com/wesleywu/gcontainer/gqueue"
)

// Runs a few goroutines calling the different methods on RingQueue concurrently.
func TestRingQueue_race(t *testing.T) {
	const (
		// Number of times each function should be called.
		Steps = 1000

		// Number of goroutines calling each function.
		Workers = 10
	)

	q := gqueue.NewRingQueue[int](0, true)
	funcs := []func(){
		func() { q.IsEmpty() },
		func() { q.Len() },
		q.Clear,
		func() { q.Push(0) },
		func() { q.Pop() },
		func() { q.Peek() },
		func() { q.Snapshot(nil) },
	}

	var (
		ready sync.WaitGroup // to block start
		done  sync.WaitGroup // to wait for end
	)
	for _, fn := range funcs {
		fn := fn
		done.Add(Workers)
		ready.Add(Workers)
		for i := 0; i < Workers; i++ {
			go func() {
				defer done.Done()

				ready.Done() // I'm ready...
				ready.Wait() // ...is everone else?

				for step := 0; step < Steps; step++ {
					fn()
				}
			}()
		}
	}

	done.Wait()
}
