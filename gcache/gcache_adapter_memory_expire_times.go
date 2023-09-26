// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"sync"
)

type adapterMemoryExpireTimes[K comparable] struct {
	mu          sync.RWMutex // expireTimeMu ensures the concurrent safety of expireTimes map.
	expireTimes map[K]int64  // expireTimes is the expiring key to its timestamp mapping, which is used for quick indexing and deleting.
}

func newAdapterMemoryExpireTimes[K comparable]() *adapterMemoryExpireTimes[K] {
	return &adapterMemoryExpireTimes[K]{
		expireTimes: make(map[K]int64),
	}
}

func (d *adapterMemoryExpireTimes[K]) Get(key K) (value int64) {
	d.mu.RLock()
	value = d.expireTimes[key]
	d.mu.RUnlock()
	return
}

func (d *adapterMemoryExpireTimes[K]) Set(key K, value int64) {
	d.mu.Lock()
	d.expireTimes[key] = value
	d.mu.Unlock()
}

func (d *adapterMemoryExpireTimes[K]) Delete(key K) {
	d.mu.Lock()
	delete(d.expireTimes, key)
	d.mu.Unlock()
}
