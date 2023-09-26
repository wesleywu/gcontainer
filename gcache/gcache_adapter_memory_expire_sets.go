// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"sync"

	"github.com/wesleywu/gcontainer/g"
)

type adapterMemoryExpireSets[K comparable] struct {
	mu         sync.RWMutex       // expireSetMu ensures the concurrent safety of expireSets map.
	expireSets map[int64]g.Set[K] // expireSets is the expiring timestamp to its key set mapping, which is used for quick indexing and deleting.
}

func newAdapterMemoryExpireSets[K comparable]() *adapterMemoryExpireSets[K] {
	return &adapterMemoryExpireSets[K]{
		expireSets: make(map[int64]g.Set[K]),
	}
}

func (d *adapterMemoryExpireSets[K]) Get(key int64) (result g.Set[K]) {
	d.mu.RLock()
	result = d.expireSets[key]
	d.mu.RUnlock()
	return
}

func (d *adapterMemoryExpireSets[K]) GetOrNew(key int64) (result g.Set[K]) {
	if result = d.Get(key); result != nil {
		return
	}
	d.mu.Lock()
	if es, ok := d.expireSets[key]; ok {
		result = es
	} else {
		result = g.NewHashSet[K](true)
		d.expireSets[key] = result
	}
	d.mu.Unlock()
	return
}

func (d *adapterMemoryExpireSets[K]) Delete(key int64) {
	d.mu.Lock()
	delete(d.expireSets, key)
	d.mu.Unlock()
}
