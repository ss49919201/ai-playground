package ttlcache

import (
	"errors"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var cacheStore = sync.OnceValue(func() *ttlcache.Cache[string, string] {
	c := ttlcache.New[string, string](
		ttlcache.WithDisableTouchOnHit[string, string](),
	)
	go c.Start()
	return c
})

func Set(k string, v string, ttl time.Duration) {
	cacheStore().Set(k, v, ttl)
}

func Delete(k string) {
	cacheStore().Delete(k)
}

func Load(k string) (string, error) {
	item := cacheStore().Get(k)
	if item == nil {
		return "", errors.New("not found")
	}
	return item.Value(), nil
}
