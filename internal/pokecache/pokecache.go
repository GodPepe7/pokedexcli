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
	data  map[string]cacheEntry
	mutex *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		data:  make(map[string]cacheEntry),
		mutex: &sync.Mutex{},
	}
	ticker := time.NewTicker(interval)
	go cache.reapLoop(ticker, interval)
	return cache
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cacheEntry, ok := cache.data[key]
	if !ok {
		return nil, false
	}
	return cacheEntry.val, true
}

func (cache *Cache) Add(key string, val []byte) {
	cache.data[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (cache *Cache) reapLoop(ticker *time.Ticker, interval time.Duration) {
	for {
		tick := <-ticker.C
		cache.mutex.Lock()
		for key, cacheEntry := range cache.data {
			if tick.Sub(cacheEntry.createdAt) >= interval {
				delete(cache.data, key)
			}
		}
		cache.mutex.Unlock()
	}
}
