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

func (c *Cache) NewCache() Cache {
	return Cache {
		cacheMap: map[string]cacheEntry{},
	}
}