package gqueue_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wesleywu/gcontainer/gqueue"
)

// Verifies that a queue filled exactly to capacity does not resize.
func TestRingQueue_fillNoResize(t *testing.T) {
	t.Parallel()

	q := gqueue.NewRingQueue[int](3)
	initCap := q.Cap()
	q.Push(1)
	q.Push(2)
	q.Push(3)
	assert.Equal(t, initCap, q.Cap(), "capacity")
}
