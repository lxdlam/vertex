package types

import (
	"sync"
	"sync/atomic"
)

// DataMap is string to any value map to carry data
// It provides many accessors
type DataMap interface {
	Set(string, interface{})
	Get(string) (interface{}, bool)
	Has(string) bool
	Remove(string)
	Len() int
}

type simpleDataMap struct {
	container map[string]interface{}
}

// NewSimpleDataMap returns a new simpleDataMap instance
func NewSimpleDataMap() *simpleDataMap {
	return &simpleDataMap{
		container: make(map[string]interface{}),
	}
}

// Set is a proxy method to map set
func (sdm *simpleDataMap) Set(key string, value interface{}) {
	sdm.container[key] = value
}

// Get is a proxy method to map get
func (sdm *simpleDataMap) Get(key string) (interface{}, bool) {
	val, ok := sdm.container[key]
	return val, ok
}

// Has returns if there is an item with given key
func (sdm *simpleDataMap) Has(key string) bool {
	_, ok := sdm.container[key]
	return ok
}

// Remove is a proxy method to map delete
func (sdm *simpleDataMap) Remove(key string) {
	delete(sdm.container, key)
}

// Len returns the size of the container
func (sdm *simpleDataMap) Len() int {
	return len(sdm.container)
}

// lockedDataMap contains a sync.RWMutex to ensure thread safety
type lockedDataMap struct {
	container map[string]interface{}
	mutex     sync.RWMutex
}

// NewLockedDataMap returns a new lockedDataMap instance
func NewLockedDataMap() *lockedDataMap {
	return &lockedDataMap{
		container: make(map[string]interface{}),
	}
}

// Set is a proxy method to map set with writelock
func (ldm *lockedDataMap) Set(key string, value interface{}) {
	ldm.mutex.Lock()
	defer ldm.mutex.Unlock()
	ldm.container[key] = value
}

// Get is a proxy method to map get with read lock
func (ldm *lockedDataMap) Get(key string) (interface{}, bool) {
	ldm.mutex.RLock()
	defer ldm.mutex.RUnlock()
	val, ok := ldm.container[key]
	return val, ok
}

// Has returns if there is an item with given key with read lock
func (ldm *lockedDataMap) Has(key string) bool {
	ldm.mutex.RLock()
	defer ldm.mutex.RUnlock()
	_, ok := ldm.container[key]
	return ok
}

// Remove is a proxy method to map delete with write lock
func (ldm *lockedDataMap) Remove(key string) {
	ldm.mutex.Lock()
	defer ldm.mutex.Unlock()
	delete(ldm.container, key)
}

// Len returns the size of the container with read lock
func (ldm *lockedDataMap) Len() int {
	ldm.mutex.RLock()
	defer ldm.mutex.RUnlock()
	return len(ldm.container)
}

// syncDataMap is just a proxy to sync.Map which with higher memory consumption but faster
// It also contains a size to count the actual size of syncDataMap
type syncDataMap struct {
	container sync.Map
	size      int32
}

// NewSyncDataMap returns a new syncDataMap instance
func NewSyncDataMap() *syncDataMap {
	return &syncDataMap{}
}

// Set will do insert and size increment
// It will do a check if the key is exist to determine if we need to increase the size
func (sdm *syncDataMap) Set(key string, value interface{}) {
	if !sdm.Has(key) {
		atomic.AddInt32(&sdm.size, 1)
	}

	sdm.container.Store(key, value)
}

// Get is a simple proxy to Load
func (sdm *syncDataMap) Get(key string) (interface{}, bool) {
	return sdm.container.Load(key)
}

// Has will simply check if the key is exist
func (sdm *syncDataMap) Has(key string) bool {
	_, ok := sdm.container.Load(key)
	return ok
}

// Remove will first check if the key is exist
// Then do remove works, it requires one more read overhead
func (sdm *syncDataMap) Remove(key string) {
	if sdm.Has(key) {
		sdm.container.Delete(key)
		atomic.AddInt32(&sdm.size, -1)
	}
}

// Len will load size atomically so it will be safe and fast
func (sdm *syncDataMap) Len() int {
	return int(atomic.LoadInt32(&sdm.size))
}
