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
	Expires   int64
	CacheType string
}

type Cache struct {
	Items map[string]*CacheItem
}

func (i *CacheItem) Valid(time int64) bool {
	if i.Expires == 0 {
		return false
	}
	return time < i.Expires
}

func (c *Cache) Get(key string) bool {
	isValid := false
	// c.mu.Lock()
	v, ok := c.Items[key]
	if ok {
		isValid = v.Valid(time.Now().UnixNano())
	}
	// c.mu.Unlock()
	c.Delete(key)
	return isValid
}

func (c *Cache) Delete(key string) {
	delete(c.Items, key)
}

func NewCache(t string, d time.Duration) *Cache {
	c := &Cache{
		Items: map[string]*CacheItem{tools.GenRandSlug(48): {
			Expires:   time.Now().Add(d).UnixNano(),
			CacheType: t,
		}},
	}
	return c
}
