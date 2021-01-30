package cache

import (
	"time"

	"github.com/dgraph-io/ristretto"
)

type cache struct {
	cache *ristretto.Cache
}

func New() *cache {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		panic(err)
	}

	return &cache{cache: c}
}

func (c *cache) Get(key interface{}) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *cache) Set(key, value interface{}) bool {
	return c.cache.Set(key, value, 1)
}

func (c *cache) SetWithTTL(key, value interface{}, ttl time.Duration) bool {
	return c.cache.SetWithTTL(key, value, 1, ttl)
}

func (c *cache) Del(key interface{}) {
	c.cache.Del(key)
}

func (c *cache) Clear() {
	c.cache.Clear()
}

func (c *cache) Close() {
	c.cache.Close()
}
