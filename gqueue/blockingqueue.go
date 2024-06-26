// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gqueue provides dynamic/static concurrent-safe queue.
//
// Features:
//
// 1. FIFO queue(data -> list -> chan);
//
// 2. Fast creation and initialization;
//
// 3. Support dynamic queue size(unlimited queue size);
//
// 4. Blocking when reading data from queue;
package gqueue

import (
	"math"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/gtype"
)

// BlockingQueue is a concurrent-safe queue built on doubly linked list and channel.
type BlockingQueue[T any] struct {
	limit  int              // Limit for queue size.
	list   *g.LinkedList[T] // Underlying list structure for data maintaining.
	closed *gtype.Bool      // Whether queue is closed.
	events chan struct{}    // Events for data writing.
	C      chan T           // Underlying channel for data reading.
}

const (
	defaultQueueSize = 10000 // Size for queue buffer.
	defaultBatchSize = 10    // Max batch size per-fetching from list.
)

// New returns an empty queue object.
// Optional parameter `limit` is used to limit the size of the queue, which is unlimited in default.
// When `limit` is given, the queue will be static and high performance which is any with stdlib channel.
func New[T any](limit ...int) *BlockingQueue[T] {
	q := &BlockingQueue[T]{
		closed: gtype.NewBool(),
	}
	if len(limit) > 0 && limit[0] > 0 {
		q.limit = limit[0]
		q.C = make(chan T, limit[0])
	} else {
		q.list = g.NewLinkedList[T](true)
		q.events = make(chan struct{}, math.MaxInt32)
		q.C = make(chan T, defaultQueueSize)
		go q.asyncLoopFromListToChannel()
	}
	return q
}

// Push pushes the data `v` into the queue.
// Note that it would panic if Push is called after the queue is closed.
func (q *BlockingQueue[T]) Push(v T) {
	if q.limit > 0 {
		q.C <- v
	} else {
		q.list.PushBack(v)
		if len(q.events) < defaultQueueSize {
			q.events <- struct{}{}
		}
	}
}

// MustPop pops an item from the queue in FIFO way.
// Note that it would return empty value of T or nil if T is a pointer, when Pop is called after the queue is closed.
func (q *BlockingQueue[T]) MustPop() T {
	return <-q.C
}

// Pop pops an item from the queue in FIFO way, and a bool value indicating whether the channel is still open.
func (q *BlockingQueue[T]) Pop() (result T, ok bool) {
	result, ok = <-q.C
	return
}

// Close closes the queue.
// Notice: It would notify all goroutines return immediately,
// which are being blocked reading using Pop method.
func (q *BlockingQueue[T]) Close() {
	if !q.closed.Cas(false, true) {
		return
	}
	if q.events != nil {
		close(q.events)
	}
	if q.limit > 0 {
		close(q.C)
	} else {
		for i := 0; i < defaultBatchSize; i++ {
			q.Pop()
		}
	}
}

// Len returns the length of the queue.
// Note that the result might not be accurate if using unlimited queue size as there's an
// asynchronous channel reading the list constantly.
func (q *BlockingQueue[T]) Len() (length int64) {
	bufferedSize := int64(len(q.C))
	if q.limit > 0 {
		return bufferedSize
	}
	return int64(q.list.Size()) + bufferedSize
}

// Size is alias of Len.
// Deprecated: use Len instead.
func (q *BlockingQueue[T]) Size() int64 {
	return q.Len()
}

// asyncLoopFromListToChannel starts an asynchronous goroutine,
// which handles the data synchronization from list `q.list` to channel `q.C`.
func (q *BlockingQueue[T]) asyncLoopFromListToChannel() {
	defer func() {
		if q.closed.Val() {
			_ = recover()
		}
	}()
	for !q.closed.Val() {
		<-q.events
		for !q.closed.Val() {
			if bufferLength := q.list.Len(); bufferLength > 0 {
				// When q.C is closed, it will panic here, especially q.C is being blocked for writing.
				// If any error occurs here, it will be caught by recover and be ignored.
				for i := 0; i < bufferLength; i++ {
					if front, ok := q.list.PopFront(); ok {
						q.C <- front
					}
				}
			} else {
				break
			}
		}
		// Clear q.events to remain just one event to do the next synchronization check.
		for i := 0; i < len(q.events)-1; i++ {
			<-q.events
		}
	}
	// It should be here to close `q.C` if `q` is unlimited size.
	// It's the sender's responsibility to close channel when it should be closed.
	close(q.C)
}
