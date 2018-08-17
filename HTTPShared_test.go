package HTTPShared

import (
	"testing"
	"time"
)

func TestSharedGetSet(t *testing.T) {
	shared := NewShared()
	shared.Set("foo", "bar1")
	shared.Get("foo")
	shared.Get("no-such-key")

	var watchId uint64
	go func() {
		watchId = shared.Watch2("foo", func(result *Result) bool {
			return false
		})
	}()
	shared.Set("foo", "bar2")
	shared.UnWatch("foo", watchId)

	go func() {
		shared.Watch2("no-such-key", func(result *Result) bool {
			return true
		})
	}()
	time.Sleep(time.Millisecond * 200)
	shared.Set("no-such-key", "value")
}

func TestSharedWatch(t *testing.T) {
	shared := NewShared()
	shared.Set("foo", "bar1")

	go func() {
		shared.Watch("foo", func(result *Result) bool {
			return false
		})
	}()
	time.Sleep(time.Millisecond * 200)

	shared.Watch("no-such-key", func(result *Result) bool {
		return false
	})
}
