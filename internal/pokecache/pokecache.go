package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]cacheEntry
	mutex sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	data []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache {
		cacheMap: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, data []byte) {
	newEntry := cacheEntry{time.Now(), data}
	c.mutex.Lock()
	c.cacheMap[key] = newEntry
	c.mutex.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, ok := c.cacheMap[key]
	if !ok {
		return []byte{}, false
	}
	return entry.data, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		c.mutex.Lock()
		for key, entry := range c.cacheMap {
			if t.Sub(entry.createdAt) > interval {
				delete(c.cacheMap, key)
			}
		}
		c.mutex.Unlock()
	}
}