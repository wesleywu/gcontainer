// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/gtimer"
	"github.com/wesleywu/gcontainer/gtype"
)

// LRU cache object.
// It uses list.List from stdlib for its underlying doubly linked list.
type adapterMemoryLru[K comparable, V any] struct {
	cache   *AdapterMemory[K, V]    // Parent cache object.
	data    g.Map[K, *g.Element[K]] // Key mapping to the item of the list.
	list    *g.LinkedList[K]        // Key list.
	rawList *g.LinkedList[K]        // History for key adding.
	closed  *gtype.Bool             // Closed or not.
}

// newMemCacheLru creates and returns a new LRU object.
func newMemCacheLru[K comparable, V any](cache *AdapterMemory[K, V]) *adapterMemoryLru[K, V] {
	lru := &adapterMemoryLru[K, V]{
		cache:   cache,
		data:    g.NewHashMap[K, *g.Element[K]](true),
		list:    g.NewLinkedList[K](true),
		rawList: g.NewLinkedList[K](true),
		closed:  gtype.NewBool(),
	}
	return lru
}

// Close closes the LRU object.
func (lru *adapterMemoryLru[K, V]) Close() {
	lru.closed.Set(true)
}

// Remove deletes the `key` FROM `lru`.
func (lru *adapterMemoryLru[K, V]) Remove(key K) {
	if v := lru.data.Get(key); v != nil {
		lru.data.Remove(key)
		lru.list.Remove(v.Value)
	}
}

// Size returns the size of `lru`.
func (lru *adapterMemoryLru[K, V]) Size() int {
	return lru.data.Size()
}

// Push pushes `key` to the tail of `lru`.
func (lru *adapterMemoryLru[K, V]) Push(key K) {
	lru.rawList.PushBack(key)
}

// Pop deletes and returns the key from tail of `lru`.
func (lru *adapterMemoryLru[K, V]) Pop() (k K, ok bool) {
	if k, ok = lru.list.PopBack(); ok {
		lru.data.Remove(k)
		return
	}
	return
}

// SyncAndClear synchronizes the keys from `rawList` to `list` and `data`
// using Least Recently Used algorithm.
func (lru *adapterMemoryLru[K, V]) SyncAndClear(ctx context.Context) {
	if lru.closed.Val() {
		gtimer.Exit()
		return
	}
	// Data synchronization.
	var alreadyExistItem *g.Element[K]
	for {
		if rawListItem, ok := lru.rawList.PopFront(); ok {
			// Deleting the key from list.
			if alreadyExistItem = lru.data.Get(rawListItem); alreadyExistItem != nil {
				lru.list.Remove(alreadyExistItem.Value)
			}
			// Pushing key to the head of the list
			// and setting its list item to hash table for quick indexing.
			lru.data.Put(rawListItem, lru.list.PushFront(rawListItem))
		} else {
			break
		}
	}
	// Data cleaning up.
	for clearLength := lru.Size() - lru.cache.cap; clearLength > 0; clearLength-- {
		if topKey, ok := lru.Pop(); ok {
			lru.cache.clearByKey(topKey, true)
		}
	}
}
