package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewEventsCache(5 * time.Minute)

	c.Set("test", []byte("123"))
	_, _ = c.Get("test")
	c.Invalidate()
}
