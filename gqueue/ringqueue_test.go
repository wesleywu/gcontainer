package gqueue_test

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wesleywu/gcontainer/gqueue"
)

func TestRingQueue(t *testing.T) {
	t.Parallel()

	testQueueSuite(t, func(capacity int) queue[int] {
		return gqueue.NewRingQueue[int](capacity)
	})
}

func TestRingQueueThreadSafe(t *testing.T) {
	t.Parallel()

	testQueueSuite(t, func(capacity int) queue[int] {
		return gqueue.NewRingQueue[int](capacity, true)
	})
}

type queue[T any] interface {
	IsEmpty() bool
	Cap() int
	Len() int
	Clear()
	Push(x T)
	Pop() (T, bool)
	PopMulti(int) []T
	PopAll() []T
	Peek() (T, bool)
	Snapshot([]T) []T
}

var (
	_ queue[int] = (*gqueue.RingQueue[int])(nil)
)

func testQueueSuite(t *testing.T, newWithCap func(capacity int) queue[int]) {
	capacities := []int{
		-1, // special case: use zero value
		0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024,
	}
	sizes := []int{1, 10, 100, 1000, 10000}
	batchSizes := []int{1, 3, 7, 15, 31, 63, 127, 255, 511, 1023, 2047, 4095}

	for _, capacity := range capacities {
		for _, size := range sizes {
			require.Greater(t, size, 0,
				"invalid test: sizes must be greater than 0")

			capacity, size := capacity, size
			name := fmt.Sprintf("Capacity=%d/Size=%d", capacity, size)
			newEmpty := func() queue[int] {
				if capacity < 0 {
					return new(gqueue.RingQueue[int])
				}
				return newWithCap(capacity)
			}

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				suite := &queueSuite{
					NewEmpty:   newEmpty,
					NumItems:   size,
					BatchSizes: batchSizes,
				}

				suitev := reflect.ValueOf(suite)
				suitet := suitev.Type()
				for i := 0; i < suitet.NumMethod(); i++ {
					name, ok := cutPrefix(suitet.Method(i).Name, "Test")
					if !ok {
						continue
					}

					testfn, ok := suitev.Method(i).Interface().(func(*testing.T))
					if !ok {
						continue
					}

					t.Run(name, testfn)
				}
			})
		}
	}
}

type queueSuite struct {
	NewEmpty   func() queue[int]
	NumItems   int
	BatchSizes []int
}

func (s *queueSuite) TestEmpty(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	assert.True(t, q.IsEmpty(), "empty")
	assert.Zero(t, q.Len(), "length")

	t.Run("TryPeekPop", func(t *testing.T) {
		_, ok := q.Peek()
		assert.False(t, ok, "peek should fail")

		_, ok = q.Pop()
		assert.False(t, ok, "pop should fail")
	})

	assert.Empty(t, q.Snapshot(nil), "snapshot")
}

func (s *queueSuite) TestPushPop(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
	}
	assert.False(t, q.IsEmpty(), "empty")
	assert.Equal(t, s.NumItems, q.Len(), "length")

	for i := 0; i < s.NumItems; i++ {
		assert.Equal(t, i, requirePeek(t, q), "peek")
		assert.Equal(t, i, requirePop(t, q), "pop")
	}

	assert.True(t, q.IsEmpty(), "empty")
	assert.Zero(t, q.Len(), "length")
}

func (s *queueSuite) TestPushPopInterleaved(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
		assert.Equal(t, i, requirePeek(t, q), "peek")
		assert.Equal(t, i, requirePop(t, q), "pop")
	}

	assert.True(t, q.IsEmpty(), "empty")
	assert.Zero(t, q.Len(), "length")
}

func (s *queueSuite) TestPushPopWraparound(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
		q.Push(requirePop(t, q))
	}

	got := make([]int, 0, q.Len())
	for !q.IsEmpty() {
		got = append(got, requirePop(t, q))
	}
	sort.Ints(got)

	want := make([]int, 0, s.NumItems)
	for i := 0; i < s.NumItems; i++ {
		want = append(want, i)
	}

	assert.Equal(t, want, got, "items")
}

func (s *queueSuite) TestPushPopMultiWithPushingBack(t *testing.T) {
	t.Parallel()

	for _, batchSize := range s.BatchSizes {
		q := s.NewEmpty()
		for i := 0; i < s.NumItems; i++ {
			q.Push(i)
		}
		assert.False(t, q.IsEmpty(), "empty")
		assert.Equal(t, s.NumItems, q.Len(), "length")
		capacity := q.Cap()
		currentLength := q.Len()

		for i := 0; i < (capacity/batchSize)+3; i++ { // test batchPop 3 times
			poppedItems := q.PopMulti(batchSize)
			poppedCount := len(poppedItems)
			assert.True(t, poppedCount <= batchSize, "popped count")
			if poppedCount < batchSize {
				assert.True(t, q.IsEmpty(), "empty")
			}
			assert.Equal(t, currentLength, q.Len()+poppedCount, "length")
			// insert the same count of popped items into queue
			for j := 0; j < poppedCount; j++ {
				q.Push(j)
			}
			assert.Equal(t, capacity, q.Cap(), "capacity")
			currentLength = q.Len()
		}
	}
}

func (s *queueSuite) TestPushPopMultiWithoutPushingBack(t *testing.T) {
	t.Parallel()

	for _, batchSize := range s.BatchSizes {
		q := s.NewEmpty()
		for i := 0; i < s.NumItems; i++ {
			q.Push(i)
		}
		assert.False(t, q.IsEmpty(), "empty")
		assert.Equal(t, s.NumItems, q.Len(), "length")
		capacity := q.Cap()
		currentLength := q.Len()

		for i := 0; i < (capacity/batchSize)+3; i++ { // test batchPop 3 times
			poppedItems := q.PopMulti(batchSize)
			poppedCount := len(poppedItems)
			assert.True(t, poppedCount <= batchSize, "popped count")
			if poppedCount < batchSize {
				assert.True(t, q.IsEmpty(), "empty")
			}
			assert.Equal(t, currentLength, q.Len()+poppedCount, "length")
			assert.Equal(t, capacity, q.Cap(), "capacity")
			currentLength = q.Len()
		}
	}
}

func (s *queueSuite) TestPushPopAll(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
	}
	assert.False(t, q.IsEmpty(), "empty")
	assert.Equal(t, s.NumItems, q.Len(), "length")
	capacity := q.Cap()
	currentLength := q.Len()

	for i := 0; i < capacity; i += 10 { // test PopAll
		poppedItems := q.PopAll()
		poppedCount := len(poppedItems)
		assert.True(t, poppedCount == currentLength, "popped count")
		assert.True(t, q.IsEmpty(), "empty")
		// insert some items into queue again
		for j := 0; j < i; j++ {
			q.Push(j)
		}
		assert.Equal(t, capacity, q.Cap(), "capacity")
		currentLength = q.Len()
	}
}

func (s *queueSuite) TestSnapshot(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
	}

	snap := q.Snapshot(nil /* dst */)
	assert.Len(t, snap, q.Len(), "length")
	for _, item := range snap {
		assert.Equal(t, item, requirePop(t, q), "item")
	}
}

func (s *queueSuite) TestSnapshotReuse(t *testing.T) {
	t.Parallel()

	q := s.NewEmpty()
	for i := 0; i < s.NumItems; i++ {
		q.Push(i)
	}

	snap := []int{42}
	snap = q.Snapshot(snap)
	assert.Len(t, snap, q.Len()+1, "length")

	assert.Equal(t, 42, snap[0], "item")
	for _, item := range snap[1:] {
		assert.Equal(t, item, requirePop(t, q), "item")
	}
}

func TestRingQueue_empty(t *testing.T) {
	t.Parallel()

	var (
		q  gqueue.RingQueue[int]
		ok bool
	)
	assert.True(t, q.IsEmpty(), "empty")
	assert.Zero(t, q.Len(), "len")
	_, ok = q.Peek()
	assert.False(t, ok, "peek")
	_, ok = q.Pop()
	assert.False(t, ok, "pop")
	assert.Empty(t, q.Snapshot(nil), "snapshot")
}

func TestRingQueue_PeekPop(t *testing.T) {
	t.Parallel()

	var q gqueue.RingQueue[int]
	q.Push(42)
	assert.Equal(t, 42, q.MustPeek(), "peek")
	assert.Equal(t, 42, q.MustPop(), "pop")
	assert.True(t, q.IsEmpty(), "empty")
}

func TestRingQueue_PopMulti(t *testing.T) {
	t.Parallel()

	q := gqueue.NewRingQueue[int](1026)
	size := 1000
	for i := 0; i < size; i++ {
		q.Push(i)
	}
	batchSize := 15
	capacity := q.Cap()
	currentLength := q.Len()

	for i := 0; i < 350; i++ { // test batchPop 3 times
		poppedItems := q.PopMulti(batchSize)
		poppedCount := len(poppedItems)
		assert.True(t, poppedCount <= batchSize, "popped count")
		if poppedCount < batchSize {
			assert.True(t, q.IsEmpty(), "empty")
		}
		assert.Equal(t, currentLength, q.Len()+poppedCount, "length")
		//insert the same count of popped items into queue
		for j := 0; j < poppedCount; j++ {
			q.Push(j)
		}
		assert.Equal(t, capacity, q.Cap(), "capacity")
		currentLength = q.Len()
	}
}

func requirePeek[T any](t require.TestingT, q queue[T]) T {
	v, ok := q.Peek()
	require.True(t, ok, "peek")
	return v
}

func requirePop[T any](t require.TestingT, q queue[T]) T {
	v, ok := q.Pop()
	require.True(t, ok, "pop")
	return v
}

// Copy of strings.CutPrefix for Go 1.19.
// Delete once Go 1.20 is minimum supported version.
func cutPrefix(s, prefix string) (after string, found bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}
