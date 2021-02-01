package cache

import (
	"time"

	"github.com/dgraph-io/ristretto"
)

type Cache struct {
	cache *ristretto.Cache
}

func New() *Cache {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		panic(err)
	}

	return &Cache{cache: c}
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *Cache) Set(key, value interface{}) bool {
	return c.cache.Set(key, value, 1)
}

func (c *Cache) SetWithTTL(key, value interface{}, ttl time.Duration) bool {
	return c.cache.SetWithTTL(key, value, 1, ttl)
}

func (c *Cache) Del(key interface{}) {
	c.cache.Del(key)
}

func (c *Cache) Clear() {
	c.cache.Clear()
}

func (c *Cache) Close() {
	c.cache.Close()
}
