package pokecache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		entries: make(map[string]cacheEntry),
		mu:      &sync.Mutex{},
	}

	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				cache.reapLoop(t.Add(-interval))
			}
		}
	}()

	return &cache
}

func (c Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{}

	if entry, ok := c.entries[key]; ok {
		entry.val = val
		c.entries[key] = entry

		return nil
	}

	return errors.New("failed to make new cache entry")
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.entries[key]; ok {
		return entry.val, true
	}

	return nil, false
}

func (c Cache) reapLoop(time time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, val := range c.entries {
		if val.createdAt.Before(time) {
			delete(c.entries, key)
		}
	}
}
