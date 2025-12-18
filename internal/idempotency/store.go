package idempotency

import (
    "sync"
    "time"
)

type cached struct {
    status    int
    body      []byte
    expiresAt time.Time
}

type Store struct {
    mu       sync.RWMutex
    store    map[string]cached
    locks    map[string]*sync.Mutex
    locksMu  sync.Mutex
    ttl      time.Duration
}

func NewStore(ttl time.Duration) *Store {
    s := &Store{
        store:   make(map[string]cached),
        locks:   make(map[string]*sync.Mutex),
        ttl:     ttl,
    }
    go s.gc()
    return s
}

var Default *Store

func init() {
    Default = NewStore(24 * time.Hour)
}

// Set stores a response for the given idempotency key.
func (s *Store) Set(key string, status int, body []byte) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.store[key] = cached{status: status, body: body, expiresAt: time.Now().Add(s.ttl)}
}

// Get retrieves a cached response if the key exists and has not expired.
func (s *Store) Get(key string) (int, []byte, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    c, ok := s.store[key]
    if !ok || time.Now().After(c.expiresAt) {
        return 0, nil, false
    }
    return c.status, c.body, true
}

// TryLock acquires a per-key lock to prevent concurrent execution of the same idempotent request.
// If successful, caller must call Unlock(key) when done.
// Returns true if lock acquired, false if already locked (concurrent request in flight).
func (s *Store) TryLock(key string) bool {
    s.locksMu.Lock()
    mu, exists := s.locks[key]
    if !exists {
        mu = &sync.Mutex{}
        s.locks[key] = mu
    }
    s.locksMu.Unlock()
    // Try to lock with timeout to avoid indefinite blocking
    locked := mu.TryLock()
    return locked
}

// Unlock releases the per-key lock.
func (s *Store) Unlock(key string) {
    s.locksMu.Lock()
    mu, exists := s.locks[key]
    s.locksMu.Unlock()
    if exists {
        mu.Unlock()
    }
}

// gc runs periodically to clean expired entries
func (s *Store) gc() {
    for now := range time.Tick(time.Minute) {
        s.mu.Lock()
        for k, v := range s.store {
            if now.After(v.expiresAt) {
                delete(s.store, k)
            }
        }
        s.mu.Unlock()
    }
}
