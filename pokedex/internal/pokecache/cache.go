package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu      sync.Mutex
	content map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	if interval < 0 {
		return nil
	}
	c := &Cache{
		content: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	c.content[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.content[key]
	if !found {
		return nil, false
	}

	return item.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	if interval < time.Millisecond {
		interval = time.Millisecond
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C

		c.mu.Lock()

		for key, entry := range c.content {
			time_since_creation := time.Since(entry.createdAt)
			if time_since_creation > interval {
				delete(c.content, key)
			}
		}

		c.mu.Unlock()
	}
}
