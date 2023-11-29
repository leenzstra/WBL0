package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type ICache[K comparable, V any] interface {
	GetItem(key K) (V, bool)
	SetItem(key K, value V) bool
	HasItem(key K) bool
}

type MemoryCache[K comparable, V any] struct {
	*ttlcache.Cache[K, V]
}

func New[K comparable, V any](ttl time.Duration) *MemoryCache[K, V] {
	return &MemoryCache[K, V]{
		Cache: ttlcache.New[K, V](
			ttlcache.WithTTL[K, V](ttl),
		),
	}
}

// Get автоматически продлевает время
func (c *MemoryCache[K, V]) GetItem(key K) (V, bool) {
	item := c.Get(key)
	
	if item == nil {
		return *new(V), false
	}

	return item.Value(), true
}

func (c *MemoryCache[K, V]) SetItem(key K, value V) bool {
	item := c.Set(key, value, ttlcache.DefaultTTL)
	return item != nil
}

func (c *MemoryCache[K, V]) HasItem(key K) bool {
	return c.Has(key)
}

