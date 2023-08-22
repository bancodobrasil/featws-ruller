package types

import (
	"sync"
	"time"
)

type CacheItem[C any] struct {
	Value      C
	Expiration time.Time
}

type Cache[K comparable, V any] struct {
	items map[K]CacheItem[V]
	mu    sync.RWMutex
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]CacheItem[V]),
	}
}

func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem[V]{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.Expiration) {
		delete(c.items, key)
		return nil, false
	}
	return &item.Value, true
}
