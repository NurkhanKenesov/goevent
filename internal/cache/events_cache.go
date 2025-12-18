package cache

import (
	"sync"
	"time"
)

type EventsCache struct {
	mu      sync.RWMutex
	data    map[string]any
	expires map[string]time.Time
	ttl     time.Duration
}

func NewEventsCache(ttl time.Duration) *EventsCache {
	return &EventsCache{
		data:    make(map[string]any),
		expires: make(map[string]time.Time),
		ttl:     ttl,
	}
}

func (c *EventsCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	exp, ok := c.expires[key]
	if !ok || time.Now().After(exp) {
		return nil, false
	}

	return c.data[key], true
}

func (c *EventsCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.expires[key] = time.Now().Add(c.ttl)
}

func (c *EventsCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]any)
	c.expires = make(map[string]time.Time)
}
