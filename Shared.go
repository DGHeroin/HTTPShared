package HTTPShared

import (
	"sync"
	"errors"
)

// Shared
type Shared struct {
	sync.RWMutex
	values map[string] *Storage
}

var (
	KeyNotFound = errors.New("Key Not Found")
)

// New A Shared
func NewShared() *Shared {
	var shared Shared
	shared.values = make(map[string] *Storage)
	return &shared
}
// Get Value Of Shared
// return nil, when no data
func (this *Shared) Get(key string) (*Result) {
	this.RLock()
	defer this.RUnlock()
	if v, ok := this.values[key]; ok {
		return v.Get()
	}
	return nil
}
// Set Value To Shared
// return the modified version
func (this *Shared) Set(key, value string) uint64 {
	this.Lock()
	defer this.Unlock()
	var storage *Storage
	var ok bool
	if storage, ok = this.values[key]; !ok {
		storage = &Storage{Key:key}
		this.values[key] = storage
	}
	return storage.Put(value)
}

// Watch Value Modify Action Of Key
// return error, when no such key
func (this* Shared) Watch(key string, callback func(result *Result)(bool)) (uint64, error) {
	this.Lock()
	defer this.Unlock()
	if storage, ok := this.values[key]; ok {
		id := storage.Watch(callback)
		return id, nil
	}
	return 0, KeyNotFound
}
// Watch Value Modify Action Of Key
// can watch a key, whatever it's has not be set
func (this* Shared) Watch2(key string, callback func(result *Result)(bool)) (uint64) {
	this.Lock()
	defer this.Unlock()
	var storage *Storage
	var ok bool
	if storage, ok = this.values[key]; !ok {
		storage = &Storage{Key:key}
		this.values[key] = storage
	}
	id := storage.Watch(callback)
	return id
}
// Unwatch A Modify Action Callback
func (this*Shared) UnWatch(key string, id uint64) {
	this.Lock()
	defer this.Unlock()
	if storage, ok := this.values[key]; ok {
		storage.UnWatch(id)
	}
}