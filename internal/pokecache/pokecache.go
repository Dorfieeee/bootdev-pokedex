package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	ttl     time.Duration
	mu      sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
		ttl:     interval,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, prst := c.entries[key]
	if !prst {
		return nil, prst
	}
	return entry.val, prst
}

// Clear cached entries
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for name, entry := range c.entries {
			if now.Sub(entry.createdAt) > c.ttl {
				delete(c.entries, name)
				// fmt.Printf("Removed %s from cache\n", name)
			}
		}
		c.mu.Unlock()
	}
}
