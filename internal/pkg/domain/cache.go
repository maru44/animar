package domain

import (
	"sync"
	"time"
)

// // cached data
// type CacheItem struct {
// 	Value   string
// 	Expires int64
// }

// // struct for caching
// type Cache struct {
// 	Items map[string][]*CacheItem
// 	mu    sync.Mutex
// }

// func (c *CacheItem) Expired(time int64) bool {
// 	if c.Expires == 0 {
// 		return false
// 	}
// 	return time > c.Expires
// }

// func (c *Cache) Get(key string) (items []*CacheItem, err error) {
// 	c.mu.Lock()
// 	if v, ok := c.Items[key]; ok {
// 		items = v
// 	} else {
// 		// err = errors.New(fmt.Sprintf("%s has Expired", key))
// 	}
// 	c.mu.Unlock()
// 	return
// }

type CacheItem struct {
	Expires int64
}

type Cache struct {
	Items map[string]*CacheItem
	mu    sync.Mutex
}

func (i *CacheItem) Valid(time int64) bool {
	if i.Expires == 0 {
		return false
	}
	return time < i.Expires
}

func (c *Cache) Get(key string) bool {
	isValid := false
	c.mu.Lock()
	v, ok := c.Items[key]
	if ok {
		isValid = v.Valid(time.Now().UnixNano())
	}
	c.mu.Unlock()
	return isValid
}

func (c *Cache) Delete(key string) {
	_, ok := c.Items[key]
	if ok {
		delete(c.Items, key)
	}
}

// func NewCache() *Cache {
// 	c := &Cache{
// 		Items: CacheItem{
// 			Expires: ,
// 		},
// 	}
// }
