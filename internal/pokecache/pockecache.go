package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}
type cacheEntry struct {
	value     []byte
	createdAt time.Time
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		cache:    make(map[string]cacheEntry),
		interval: interval,
	}
	go cache.readLoop()
	return &cache

}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		value:     value,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, found := c.cache[key]
	if !found {
		return nil, false
	}

	if time.Since(entry.createdAt) > c.interval {
		delete(c.cache, key)
		return nil, false
	}

	return entry.value, true
}

func (c *Cache) readLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, entry := range c.cache {
				if time.Since(entry.createdAt) > c.interval {
					delete(c.cache, key)
				}
			}
			c.mu.Unlock()
		}

	}
}
