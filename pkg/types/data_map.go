package types

import (
	"sync"
	"sync/atomic"
)

// DataMap is string to any value map to carry datas
// It provides many accessors
type DataMap interface {
	Set(string, interface{})
	Get(string) (interface{}, bool)
	Has(string) bool
	Remove(string)
	Len() int
}

// SimpleDataMap with no sync operations
type SimpleDataMap struct {
	container map[string]interface{}
}

// NewSimpleDataMap returns a new SimpleDataMap instance
func NewSimpleDataMap() *SimpleDataMap {
	return &SimpleDataMap{
		container: make(map[string]interface{}),
	}
}

// Set is a proxy method to map set
func (sdm *SimpleDataMap) Set(key string, value interface{}) {
	sdm.container[key] = value
}

// Get is a proxy method to map get
func (sdm *SimpleDataMap) Get(key string) (interface{}, bool) {
	val, ok := sdm.container[key]
	return val, ok
}

// Has returns if there is an item with given key
func (sdm *SimpleDataMap) Has(key string) bool {
	_, ok := sdm.container[key]
	return ok
}

// Remove is a proxy method to map delete
func (sdm *SimpleDataMap) Remove(key string) {
	delete(sdm.container, key)
}

// Len returns the size of the container
func (sdm *SimpleDataMap) Len() int {
	return len(sdm.container)
}

// LockedDataMap contains a sync.RWMutex to ensure thread safety
type LockedDataMap struct {
	container map[string]interface{}
	mutex     sync.RWMutex
}

// NewLockedDataMap returns a new LockedDataMap instance
func NewLockedDataMap() *LockedDataMap {
	return &LockedDataMap{
		container: make(map[string]interface{}),
	}
}

// Set is a proxy method to map set with writelock
func (ldm *LockedDataMap) Set(key string, value interface{}) {
	ldm.mutex.Lock()
	defer ldm.mutex.Unlock()
	ldm.container[key] = value
}

// Get is a proxy method to map get with read lock
func (ldm *LockedDataMap) Get(key string) (interface{}, bool) {
	ldm.mutex.RLock()
	defer ldm.mutex.RUnlock()
	val, ok := ldm.container[key]
	return val, ok
}

// Has returns if there is an item with given key with read lock
func (ldm *LockedDataMap) Has(key string) bool {
	ldm.mutex.RLock()
	defer ldm.mutex.RUnlock()
	_, ok := ldm.container[key]
	return ok
}

// Remove is a proxy method to map delete with write lock
func (ldm *LockedDataMap) Remove(key string) {
	ldm.mutex.Lock()
	defer ldm.mutex.Unlock()
	delete(ldm.container, key)
}

// Len returns the size of the container with read lock
func (ldm *LockedDataMap) Len() int {
	ldm.mutex.RLock()
	defer ldm.mutex.RUnlock()
	return len(ldm.container)
}

// SyncDataMap is just a proxy to sync.Map which with higher memory consumption but faster
// It also contains a size to count the actual size of SyncDataMap
type SyncDataMap struct {
	container sync.Map
	size      int32
}

// NewSyncDataMap returns a new SyncDataMap instance
func NewSyncDataMap() *SyncDataMap {
	return &SyncDataMap{}
}

// Set will do insert and size increament
// It will do a check if the key is exist to determine if we need to increase the size
func (sdm *SyncDataMap) Set(key string, value interface{}) {
	if !sdm.Has(key) {
		atomic.AddInt32(&sdm.size, 1)
	}

	sdm.container.Store(key, value)
}

// Get is a simple proxy to Load
func (sdm *SyncDataMap) Get(key string) (interface{}, bool) {
	return sdm.container.Load(key)
}

// Has will simply check if the key is exist
func (sdm *SyncDataMap) Has(key string) bool {
	_, ok := sdm.container.Load(key)
	return ok
}

// Remove will first check if the key is exist
// Then do remove works, it requires one more read overhead
func (sdm *SyncDataMap) Remove(key string) {
	if sdm.Has(key) {
		sdm.container.Delete(key)
		atomic.StoreInt32(&sdm.size, atomic.LoadInt32(&sdm.size)-1)
	}
}

// Len will load size atomically so it will be safe and fast
func (sdm *SyncDataMap) Len() int {
	return int(atomic.LoadInt32(&sdm.size))
}
