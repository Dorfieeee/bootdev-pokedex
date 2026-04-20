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
	mu      sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{}
	c.reapLoop(interval)
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
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			c.mu.Lock()
			defer c.mu.Unlock()
			now := time.Now()
			for name, entry := range c.entries {
				if now.Compare(entry.createdAt.Add(interval)) == -1 {
					delete(c.entries, name)
				}
			}
		}
	}()
}
