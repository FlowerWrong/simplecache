package simplecache

import (
	"errors"
	"log"
	"sync"
	"time"
)

const (
	// ErrEmptyKey is for empty error key
	ErrEmptyKey string = "Empty key"
	// ErrMemLimit is for memory limit
	ErrMemLimit string = "Reach the memory limit size"
)

// Cache is interface for many cache adapter, now just for simple
// type Cache interface {
// 	SetMaxMemory(size string) bool

// 	Set(key string, val interface{}, expire time.Duration) error

// 	Get(key string) (interface{}, bool)

// 	Del(key string)

// 	Exists(key string) bool

// 	Flush()

// 	Size() int

// 	Keys() []string

// 	GC()
// }

// CacheMap is a cache adapter struct, using map structure
type CacheMap struct {
	items     map[string]item // TODO use ring
	rw        sync.RWMutex
	maxMemory uint64
	gcPercent uint64 // GC percent in second
	nbytes    uint64 // bytes for key + value
}

type item struct {
	data   interface{}
	expire int64
}

// New a simple cache
func New(gcPercent uint64) *CacheMap {
	c := &CacheMap{
		items:     make(map[string]item),
		gcPercent: gcPercent,
		nbytes:    0,
	}
	go c.GC()
	return c
}

// SetMaxMemory set the simple max memory limit
func (c *CacheMap) SetMaxMemory(size string) bool {
	maxMemory, err := ToBytes(size)
	if err != nil {
		return false
	}
	c.maxMemory = maxMemory
	return true
}

// Set key value and expire time
func (c *CacheMap) Set(key string, val interface{}, expire time.Duration) error {
	if key == "" {
		return errors.New(ErrEmptyKey)
	}

	c.rw.Lock()
	_, found := c.items[key]
	if !found {
		c.nbytes += uint64(len(key) + GetRealSizeOf(val))

		if c.nbytes > c.maxMemory {
			return errors.New(ErrMemLimit)
		}
	}

	c.items[key] = item{
		data:   val,
		expire: time.Now().Add(expire).UnixNano(),
	}
	c.rw.Unlock()
	return nil
}

// Get by key
func (c *CacheMap) Get(key string) (interface{}, bool) {
	c.rw.RLock()
	item, found := c.items[key]
	if !found {
		c.rw.RUnlock()
		return nil, false
	}
	if item.isExpire() {
		c.rw.RUnlock()
		return nil, false
	}
	c.rw.RUnlock()
	return item.data, true
}

// Del key
func (c *CacheMap) Del(key string) {
	c.rw.Lock()
	item, found := c.items[key]
	if !found {
		c.rw.Unlock()
		return
	}
	delete(c.items, key)
	c.nbytes -= uint64(len(key) + GetRealSizeOf(item.data))
	c.rw.Unlock()
}

// Exists or not for key
func (c *CacheMap) Exists(key string) bool {
	c.rw.RLock()
	item, found := c.items[key]
	if !found {
		c.rw.RUnlock()
		return false
	}
	if item.isExpire() {
		c.rw.RUnlock()
		return false
	}
	c.rw.RUnlock()
	return true
}

// Flush all keys
func (c *CacheMap) Flush() {
	c.rw.Lock()
	c.items = make(map[string]item)
	c.nbytes = 0
	c.rw.Unlock()
}

// Size of cache
func (c *CacheMap) Size() int {
	c.rw.RLock()
	size := len(c.items)
	c.rw.RUnlock()
	return size
}

// Keys list
func (c *CacheMap) Keys() (keys []string) {
	c.rw.RLock()
	for key := range c.items {
		keys = append(keys, key)
	}
	c.rw.RUnlock()
	return keys
}

// GC to delete expire key
func (c *CacheMap) GC() {
	for {
		select {
		case <-time.After(time.Duration(c.gcPercent) * time.Second):
			if keys := c.expiredKeys(); len(keys) != 0 {
				for _, key := range keys {
					c.Del(key)
				}
			}
			log.Println("GC")
		}
	}
}

func (c *CacheMap) expiredKeys() (keys []string) {
	c.rw.RLock()
	for key, item := range c.items {
		if item.isExpire() {
			keys = append(keys, key)
		}
	}
	c.rw.RUnlock()
	return keys
}

func (i *item) isExpire() bool {
	if time.Now().UnixNano() > i.expire {
		return true
	}
	return false
}
