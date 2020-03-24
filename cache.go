package easycache

import (
	"sync"
	"time"
)

/**
 * @Author: eggsy
 * @Description:
 * @File:  cache
 * @Version: 1.0.0
 * @Date: 2020-03-24 16:34
 */

type EasyCache struct {
	lock  sync.RWMutex
	items map[interface{}]item
	close chan struct{}
}

type item struct {
	expires int64
	value   interface{}
}

func New(cleaningInterval time.Duration) *EasyCache {
	var cache = new(EasyCache)
	cache.items = map[interface{}]item{}
	cache.close = make(chan struct{})
	go func() {
		ticker := time.NewTicker(cleaningInterval)
		for {
			select {
			case <-ticker.C:
				now := time.Now().Unix()
				for k, v := range cache.items {
					if v.expires > 0 && now > v.expires {
						cache.Delete(k)
					}
				}
			case <-cache.close:
				return
			}
		}
	}()
	return cache
}

func (cache *EasyCache) Set(key interface{}, value interface{}, duration time.Duration) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	var expires int64
	if duration > 0 {
		expires = time.Now().Add(duration).Unix()
	}
	cache.items[key] = item{
		expires: expires,
		value:   value,
	}
}

func (cache *EasyCache) Get(key interface{}) (interface{}, bool) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	if value, ok := cache.items[key]; ok {
		return value.value, true
	} else {
		return nil, false
	}
}

func (cache *EasyCache) Delete(key interface{}) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	delete(cache.items, key)
}

func (cache *EasyCache) Close() {
	cache.close <- struct{}{}
	cache.items = map[interface{}]item{}
}
