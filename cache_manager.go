package cache_manager

import (
	"sync"
	"time"
)

type Cache struct {
	list map[string]CacheValue
	sync.RWMutex
}

type CacheValue struct {
	data       any
	expires_at time.Time
}

func New() *Cache {
	return &Cache{
		list: make(map[string]CacheValue),
	}
}

func (c *Cache) Set(key string, value any, expires int) {
	c.Lock()
	defer c.Unlock()

	c.list[key] = CacheValue{value, time.Now().Add(time.Duration(expires) * time.Second)}
}

func (c *Cache) Get(key string) (any, bool) {
	c.RLock()
	defer c.RUnlock()

	cacheValue, ok := c.list[key]
	if !ok || cacheValue.Expires() {
		return nil, false
	}

	return cacheValue.data, ok
}

func (c *Cache) Clear(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.list, key)
}

func (c *Cache) ClearAll(key string) {
	c.Lock()
	defer c.Unlock()

	c.list = make(map[string]CacheValue)
}

func (cv *CacheValue) Expires() bool {
	return cv.expires_at.Before(time.Now())
}
