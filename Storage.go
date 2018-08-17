package HTTPShared

import (
	"sync/atomic"
	"sync"
)

// Storage
type Storage struct {
	sync.RWMutex
	Key             string
	Value           string
	ModifiedVersion uint64
	WatcherVersion  uint64
	Watcher        []*Watcher
}
// Watcher
type Watcher struct {
	Version uint64
	Callback func(*Result)bool
}
// Get Value From Storage
// when no value, return nil
func (this *Storage) Get() *Result{
	return &Result{Key: this.Key,Value:this.Value, Version:this.ModifiedVersion}
}
// Put Value To Storage
// return the modified version
func (this *Storage) Put(value string) uint64 {
	this.Lock()
	defer this.Unlock()
	this.Value = value
	atomic.AddUint64(&this.ModifiedVersion, 1)

	var rs = &Result{
		Key:this.Key,
		Value:this.Value,
		Version:this.ModifiedVersion,
	}

	for k, w := range this.Watcher {
		v := w.Callback(rs)
		if v == false {
			this.Watcher = append(this.Watcher[:k], this.Watcher[k+1:]...)
		}
	}
	return this.ModifiedVersion
}
// Watch Modify Action Of Value
// return the watcher's id
func (this *Storage) Watch(callback func(result *Result)(bool)) uint64 {
	this.Lock()
	defer this.Unlock()
	atomic.AddUint64(&this.WatcherVersion, 1)
	v := this.WatcherVersion
	this.Watcher = append(this.Watcher, &Watcher{Version:v, Callback:callback})
	return v
}
// Unwatch Modify Action By Id
func (this *Storage) UnWatch(id uint64)  {
	this.Lock()
	defer this.Unlock()
	for k, w := range this.Watcher {
		if w.Version == id {
			this.Watcher = append(this.Watcher[:k], this.Watcher[k+1:]...)
			break
		}
	}
}