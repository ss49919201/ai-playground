package ttlcache

import (
	"errors"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var newCacheStoreString = sync.OnceValue(func() *ttlcache.Cache[string, string] {
	c := ttlcache.New[string, string]()
	go c.Start()
	return c
})

func Set(k string, v string, ttl time.Duration) {
	newCacheStoreString().Set(k, v, ttl)
}

func Load(k string) (string, error) {
	item := newCacheStoreString().Get(k)
	if item == nil {
		return "", errors.New("not found")
	}
	return item.Value(), nil
}
