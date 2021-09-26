package domain

import (
	"time"
)

const (
	// interval
	// CsrfInterval = 30 * time.Minute
	CSRF_INTERVAL_MINUTE = 30

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
	v, ok := c.Items[kind][key]
	if ok {
		isValid = v.Valid(time.Now().UnixNano())
	}
	return isValid
}

func (c *Cache) delete(kind, key string) {
	delete(c.Items[kind], key)
}

func (c *Cache) DeleteExpired(kind string) {
	if v, ok := c.Items[kind]; !ok {
		return
	} else {
		for key, ci := range v {
			if !ci.Valid(time.Now().UnixNano()) {
				c.delete(kind, key)
			}
		}
	}
}

func NewCahce() *Cache {
	c := &Cache{
		Items: map[string]map[string]*CacheItem{
			CacheTypeCsrf: {},
		},
	}
	return c
}

func (c *Cache) AddCacheItem(kind, key string, d time.Duration) {
	c.Items[kind][key] = &CacheItem{
		Expires: time.Now().Add(d).UnixNano(),
	}
}

func (c *Cache) DeleteRegularly(kind string, d time.Duration) {
	// t := time.NewTicker(d)
	// defer t.Stop()

	// delete := make(chan bool)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-delete:
	// 			fmt.Print("aaaa")
	// 			c.DeleteExpired(kind)
	// 		}
	// 	}
	// }()

	go func() {
		t := time.NewTicker(d)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				// fmt.Print(c.Items[CacheTypeCsrf])
				c.DeleteExpired(kind)
			}
		}
	}()
}
