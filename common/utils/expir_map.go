package utils

import (
	"sync"
	"time"
)

type ExpiringMap struct {
	sync.RWMutex
	m map[string]*entry
}

type entry struct {
	value      interface{}
	expiration int64
}

var expiringMap ExpiringMap

func NewExpiringMap() *ExpiringMap {
	m := make(map[string]*entry)
	expiringMap = ExpiringMap{m: m}
	return &expiringMap
}

func GetExpiringMap() *ExpiringMap {
	return &expiringMap
}

func (em *ExpiringMap) Set(key string, value interface{}, expiration time.Duration) {
	em.Lock()
	defer em.Unlock()
	em.m[key] = &entry{
		value:      value,
		expiration: time.Now().Add(expiration).UnixNano(),
	}
}

func (em *ExpiringMap) Get(key string) (interface{}, bool) {
	em.RLock()
	defer em.RUnlock()
	e, ok := em.m[key]
	if !ok {
		return nil, false
	}
	if e.expiration < time.Now().UnixNano() {
		em.delete(key)
		return nil, false
	}
	return e.value, true
}

func (em *ExpiringMap) delete(key string) {
	em.Lock()
	defer em.Unlock()
	delete(em.m, key)
}
