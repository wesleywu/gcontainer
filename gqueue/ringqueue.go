package gqueue

import (
	"github.com/wesleywu/gcontainer/internal/rwmutex"
)

const _defaultCapacity = 16

// RingQueue is a FIFO queue backed by a ring buffer, designed for minimal allocation.
// The zero value for RingQueue is an empty queue ready to use, though not thread-safe.
//
// RingQueue is thread-safe for concurrent use by calling NewRingQueue with safe = true.
type RingQueue[T any] struct {
	mu rwmutex.RWMutex
	// buff is the ring buffer.
	//
	// The first item in the queue is at buff[head].
	// The last item in the queue is at buff[tail-1].
	// The queue is empty if head == tail.
	buff []T

	// head is the index of the first item in the queue.
	head int // inv: 0 <= head < len(buff)

	// tail is the index of the next empty slot in buff.
	tail int // inv: 0 <= tail < len(buff)
}

// NewRingQueue returns a new queue with the given capacity.
// If capacity is zero, the queue is initialized with a default capacity.
// If safe is true, the queue is thread-safe for concurrent use
//
// The capacity defines the leeway for bursts of pushes
// before the queue needs to grow.
func NewRingQueue[T any](capacity int, safe ...bool) *RingQueue[T] {
	if capacity == 0 {
		capacity = _defaultCapacity
	}
	return &RingQueue[T]{
		mu: rwmutex.Create(safe...),
		// Allocate requested capacity plus one slot
		// so that filling the queue to exactly the requested capacity
		// doesn't require resizing.
		buff: make([]T, capacity+1),
		head: 0,
		tail: 0,
	}
}

// IsEmpty returns true if the queue is empty.
//
// This is an O(1) operation and does not allocate.
func (q *RingQueue[T]) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.head == q.tail
}

// Cap returns the current capacity of items.
//
// This is an O(1) operation and does not allocate.
func (q *RingQueue[T]) Cap() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.buff)
}

// Len returns the number of items in the queue.
//
// This is an O(1) operation and does not allocate.
func (q *RingQueue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if q.head <= q.tail {
		return q.tail - q.head
	}
	return len(q.buff) - q.head + q.tail
}

// Clear removes all items from the queue.
// It does not adjust its internal capacity.
//
// This is an O(1) operation and does not allocate.
func (q *RingQueue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.head = 0
	q.tail = 0
}

// Push adds x to the back of the queue.
//
// This operation is O(n) in the worst case if the queue needs to grow.
// However, for target use cases, it's amortized O(1).
// See package documentation for details.
func (q *RingQueue[T]) Push(x T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.buff) == 0 {
		q.buff = make([]T, _defaultCapacity)
	}

	q.buff[q.tail] = x
	q.tail++

	if q.tail == len(q.buff) {
		// Wrap around.
		q.tail = 0
	}

	// We'll hit this only if the tail has wrapped around
	// and has caught up with the head (the queue is full).
	// In that case, we need to grow the queue
	// copying buff[head:] and buff[:tail] to the new buffer.
	if q.head == q.tail {
		// The queue is full. Make room.
		buff := make([]T, 2*len(q.buff))
		n := copy(buff, q.buff[q.head:])
		n += copy(buff[n:], q.buff[:q.tail])
		q.head = 0
		q.tail = n
		q.buff = buff
	}
}

// MustPop removes and returns the item at the front of the queue.
// It returns empty value of T if the queue is empty.
//
// This is an O(1) operation and does not allocate.
func (q *RingQueue[T]) MustPop() T {
	q.mu.Lock()
	defer q.mu.Unlock()
	t, _ := q.Pop()
	return t

}

// Pop removes and returns the item at the front of the queue.
// It returns false if the queue is empty.
// Otherwise, it returns true and the item.
//
// This is an O(1) operation and does not allocate.
func (q *RingQueue[T]) Pop() (x T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.head == q.tail {
		return x, false
	}

	x = q.buff[q.head]
	q.head++
	if q.head == len(q.buff) {
		// Wrap around.
		//
		// If tail has wrapped around too,
		// the next MustPop will catch it when head == tail.
		q.head = 0
	}
	return x, true
}

// BatchPop removes and returns specified number of items at the front of the queue.
// It returns empty slice if the queue is empty.
// Otherwise, it returns true and the item.
//
// This is an O(1) operation and does not allocate.
func (q *RingQueue[T]) BatchPop(numberOfItems int) []T {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.head == q.tail {
		return []T{}
	}

	if q.tail > q.head || numberOfItems <= len(q.buff)-q.head { // the items need to be popped are continuous
		var popCount = numberOfItems
		if q.tail > q.head {
			popCount = min(numberOfItems, q.tail-q.head)
		}
		result := make([]T, popCount)
		copy(result, q.buff[q.head:q.head+popCount])
		q.head += popCount

		if q.head == len(q.buff) {
			// Wrap around.
			//
			// If tail has wrapped around too,
			// the next MustPop will catch it when head == tail.
			q.head = 0
		}
		return result
	} else { // the items need to be popped are discontinuous
		tailedCount := len(q.buff) - q.head
		headCount := q.tail
		popCount := min(numberOfItems, tailedCount+headCount)
		result := make([]T, popCount)
		copy(result, q.buff[q.head:len(q.buff)])
		copy(result[tailedCount:], q.buff[:q.tail])
		q.head += popCount - len(q.buff)
		return result
	}
}

// MustPeek returns the item at the front of the queue without removing it.
// It returns empty value of T if the queue is empty.
//
// This is an O(1) operation and does not allocate.
func (q *RingQueue[T]) MustPeek() T {
	q.mu.Lock()
	defer q.mu.Unlock()
	x, ok := q.Peek()
	if !ok {
		panic("empty queue")
	}
	return x
}

// Peek returns the item at the front of the queue.
// It returns false if the queue is empty.
// Otherwise, it returns true and the item.
//
// This is an O(1) operation and does not allocate.
func (q *RingQueue[T]) Peek() (x T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.head == q.tail {
		return x, false
	}
	return q.buff[q.head], true
}

// Snapshot appends the contents of the queue to dst and returns the result.
// Use dst to avoid allocations when you know the capacity of the queue
//
//	dst := make([]T, 0, q.Len())
//	dst = q.Snapshot(dst)
//
// Pass nil to let the function allocate a new slice.
//
//	q.Snapshot(nil) // allocates a new slice
//
// The returned slice is a copy of the internal buffer and is safe to modify.
func (q *RingQueue[T]) Snapshot(dst []T) []T {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if q.head <= q.tail {
		return append(dst, q.buff[q.head:q.tail]...)
	}

	dst = append(dst, q.buff[q.head:]...)
	return append(dst, q.buff[:q.tail]...)
}
