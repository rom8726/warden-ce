package cache

import (
	"context"
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
)

// Key represents a cache key.
type Key interface {
	String() string
}

// Value represents a cache value.
type Value interface {
	IsValid() bool
}

// Cache interface defines the basic cache operations.
type Cache[K Key, V Value] interface {
	Get(ctx context.Context, key K) (V, bool)
	Set(ctx context.Context, key K, value V) error
	Delete(ctx context.Context, key K) error
	Clear(ctx context.Context) error
	Size() int
	Capacity() int
}

// LRUCache implements Cache interface using golang-lru.
type LRUCache[K Key, V Value] struct {
	cache    *lru.Cache[string, V]
	capacity int
	mu       sync.RWMutex
}

// NewLRUCache creates a new LRU cache with the specified capacity.
func NewLRUCache[K Key, V Value](capacity int) (*LRUCache[K, V], error) {
	lruCache, err := lru.New[string, V](capacity)
	if err != nil {
		return nil, fmt.Errorf("create LRU cache: %w", err)
	}

	return &LRUCache[K, V]{
		cache:    lruCache,
		capacity: capacity,
	}, nil
}

// Get retrieves a value from the cache.
func (c *LRUCache[K, V]) Get(_ context.Context, key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.cache.Get(key.String())

	return value, ok
}

// Set stores a value in the cache.
func (c *LRUCache[K, V]) Set(_ context.Context, key K, value V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Add(key.String(), value)

	return nil
}

// Delete removes a value from the cache.
func (c *LRUCache[K, V]) Delete(_ context.Context, key K) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Remove(key.String())

	return nil
}

// Clear removes all values from the cache.
func (c *LRUCache[K, V]) Clear(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Purge()

	return nil
}

// Size returns the current number of items in the cache.
func (c *LRUCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cache.Len()
}

// Capacity returns the maximum capacity of the cache.
func (c *LRUCache[K, V]) Capacity() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.capacity
}
