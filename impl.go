package localcache

import (
	"sync"
	"time"
)

var (
	defaultExpiredTime = 30 * time.Second
)

type cache struct {
	m       *sync.RWMutex
	storage map[string]*data
}

type data struct {
	value         interface{}
	expiredHandle *time.Timer
}

// New a cache instance
func New() Cache {
	return &cache{
		m:       &sync.RWMutex{},
		storage: make(map[string]*data),
	}
}

// Set value into cache with key
func (c *cache) Set(key string, val interface{}) error {
	c.m.Lock()
	defer c.m.Unlock()

	c.storage[key] = &data{
		value: val,
		expiredHandle: time.AfterFunc(defaultExpiredTime, func() {
			c.del(key)
		}),
	}

	return nil
}

// Get value form cache by key
func (c *cache) Get(key string) (interface{}, error) {
	c.m.RLock()
	defer c.m.RUnlock()

	if data, ok := c.storage[key]; ok {
		return data.value, nil
	}

	return nil, ErrDataNotFound
}

func (c *cache) del(key string) {
	c.m.Lock()
	defer c.m.Unlock()
	delete(c.storage, key)
}
