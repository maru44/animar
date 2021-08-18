package domain

import (
	"animar/v1/internal/pkg/tools/tools"
	"time"
)

const (
	// interval
	CsrfInterval = 120 * time.Second

	// type
	CacheTypeCsrf = "CSRF"
)

type CacheItem struct {
	Expires int64
}

type Cache struct {
	Items map[string]map[string]*CacheItem
}

func (ci *CacheItem) Valid(time int64) bool {
	if ci.Expires == 0 {
		return false
	}
	return time < ci.Expires
}

func (c *Cache) Get(kind, key string) bool {
	isValid := false
	v, ok := c.Items[kind]
	if ok {
		isValid = v[key].Valid(time.Now().UnixNano())
	}
	// c.Delete(key)
	return isValid
}

func (c *Cache) Delete(key string) {
	delete(c.Items, key)
}

// func NewCache(t string, d time.Duration) *Cache {
// 	c := &Cache{
// 		Items: map[string]*CacheItem{tools.GenRandSlug(48): {
// 			Expires:   time.Now().Add(d).UnixNano(),
// 			CacheType: t,
// 		}},
// 	}
// 	return c
// }

func NewCahce() *Cache {
	c := &Cache{
		Items: map[string]map[string]*CacheItem{},
	}
	return c
}

func (c *Cache) AddCacheItem(kind string, d time.Duration) {
	key := tools.GenRandSlug(32)
	c.Items[kind][key] = &CacheItem{
		Expires: time.Now().Add(d).UnixNano(),
	}
}
