package idempotency

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestGetSetBasic tests basic get/set functionality
func TestGetSetBasic(t *testing.T) {
	store := NewStore(5 * time.Minute)

	key := "test-key-1"
	status := 200
	body := []byte(`{"message": "success"}`)

	store.Set(key, status, body)

	// Get should find it
	gotStatus, gotBody, found := store.Get(key)
	if !found {
		t.Errorf("expected to find key, but did not")
	}
	if gotStatus != status {
		t.Errorf("expected status %d, got %d", status, gotStatus)
	}
	if string(gotBody) != string(body) {
		t.Errorf("expected body %s, got %s", body, gotBody)
	}
}

// TestGetMissing tests that Get returns false for missing keys
func TestGetMissing(t *testing.T) {
	store := NewStore(5 * time.Minute)

	_, _, found := store.Get("nonexistent")
	if found {
		t.Errorf("expected key not found, but it was")
	}
}

// TestExpiration tests that expired entries are not returned
func TestExpiration(t *testing.T) {
	store := NewStore(100 * time.Millisecond)

	key := "expire-key"
	store.Set(key, 200, []byte("data"))

	// Should exist immediately
	_, _, found := store.Get(key)
	if !found {
		t.Errorf("expected key to exist immediately after set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should not exist after expiration
	_, _, found = store.Get(key)
	if found {
		t.Errorf("expected key to expire, but it was still found")
	}
}

// TestConcurrentWrites tests that concurrent writes are safe
func TestConcurrentWrites(t *testing.T) {
	store := NewStore(5 * time.Minute)

	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			key := "concurrent-key"
			status := 200 + idx
			body := []byte("data")
			store.Set(key, status, body)
		}(i)
	}

	wg.Wait()

	// Final state should be consistent
	status, _, found := store.Get("concurrent-key")
	if !found {
		t.Errorf("expected key to exist after concurrent writes")
	}
	if status < 200 || status >= 200+numGoroutines {
		t.Errorf("unexpected status value: %d", status)
	}
}

// TestTryLockSuccess tests that TryLock succeeds when no lock exists
func TestTryLockSuccess(t *testing.T) {
	store := NewStore(5 * time.Minute)

	key := "lock-key"
	locked := store.TryLock(key)
	if !locked {
		t.Errorf("expected TryLock to succeed on fresh key")
	}

	store.Unlock(key)
}

// TestTryLockBlocked tests that TryLock fails when already locked
func TestTryLockBlocked(t *testing.T) {
	store := NewStore(5 * time.Minute)

	key := "blocked-key"
	locked1 := store.TryLock(key)
	if !locked1 {
		t.Errorf("expected first TryLock to succeed")
	}

	// Second TryLock on same key should fail
	locked2 := store.TryLock(key)
	if locked2 {
		t.Errorf("expected second TryLock to fail, but it succeeded")
	}

	store.Unlock(key)
}

// TestIdempotencyFlow tests the full idempotent request flow
func TestIdempotencyFlow(t *testing.T) {
	store := NewStore(5 * time.Minute)

	key := "idempotent-flow"
	executionCount := 0

	// First request: not cached, acquire lock, execute, store result
	if _, _, found := store.Get(key); !found {
		if !store.TryLock(key) {
			t.Fatal("expected lock acquisition to succeed")
		}
		defer store.Unlock(key)

		// Check cache again after lock (double-check)
		if _, _, found := store.Get(key); !found {
			executionCount++
			store.Set(key, 200, []byte(`{"result": "first"}`))
		}
	}

	if executionCount != 1 {
		t.Errorf("expected execution count 1, got %d", executionCount)
	}

	// Second request: should hit cache
	if cachedStatus, cachedBody, found := store.Get(key); found {
		if cachedStatus != 200 {
			t.Errorf("expected cached status 200, got %d", cachedStatus)
		}
		if string(cachedBody) != `{"result": "first"}` {
			t.Errorf("expected cached body to match")
		}
	} else {
		t.Errorf("expected second request to hit cache")
	}
}

// TestConcurrentIdempotentRequests tests that only one goroutine executes for the same key
func TestConcurrentIdempotentRequests(t *testing.T) {
	store := NewStore(5 * time.Minute)

	key := "concurrent-idempotent"
	executionCount := atomic.Int32{}
	numGoroutines := 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()

			// Try cache first
			if _, _, found := store.Get(key); found {
				return // cache hit, no execution
			}

			// Try lock
			if !store.TryLock(key) {
				// Lock failed, another goroutine is executing
				// Spin-wait for cache to be populated
				for attempts := 0; attempts < 100; attempts++ {
					if _, _, found := store.Get(key); found {
						return
					}
					time.Sleep(1 * time.Millisecond)
				}
				return
			}
			defer store.Unlock(key)

			// Check cache again (double-check)
			if _, _, found := store.Get(key); found {
				return
			}

			// Execute business logic
			executionCount.Add(1)
			store.Set(key, 200, []byte("data"))
		}()
	}

	wg.Wait()

	// Only one goroutine should have executed
	if count := executionCount.Load(); count != 1 {
		t.Errorf("expected execution count 1, got %d", count)
	}

	// Verify cached result exists
	status, _, found := store.Get(key)
	if !found {
		t.Errorf("expected cache to be populated")
	}
	if status != 200 {
		t.Errorf("expected cached status 200, got %d", status)
	}
}

// BenchmarkSet benchmarks Set performance
func BenchmarkSet(b *testing.B) {
	store := NewStore(5 * time.Minute)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Set("bench-key", 200, []byte("data"))
	}
}

// BenchmarkGet benchmarks Get performance
func BenchmarkGet(b *testing.B) {
	store := NewStore(5 * time.Minute)
	store.Set("bench-key", 200, []byte("data"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Get("bench-key")
	}
}

// BenchmarkConcurrentGetSet benchmarks concurrent access
func BenchmarkConcurrentGetSet(b *testing.B) {
	store := NewStore(5 * time.Minute)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				store.Set("key", 200, []byte("data"))
			} else {
				store.Get("key")
			}
			i++
		}
	})
}
