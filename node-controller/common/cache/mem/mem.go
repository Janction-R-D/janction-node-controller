package mem

import (
	"node-controller/common/cache"
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	memCache "github.com/eleztian/go-cache"
	"github.com/pkg/errors"
)

const TYPE = "mem"

type Cache struct {
	cache *memCache.Cache

	mu *sync.Mutex
}

func lockKey(key string) string {
	return "lock:" + key
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	return &Cache{
		cache: memCache.New(defaultExpiration, cleanupInterval),
		mu:    &sync.Mutex{},
	}
}

func (c *Cache) Type() string {
	return TYPE
}

func copyValue(src interface{}, to interface{}) error {
	toValue := reflect.ValueOf(to)
	if toValue.Kind() != reflect.Ptr || toValue.IsNil() {
		return errors.Errorf("invalid type %s or is nil", reflect.TypeOf(to).String())
	}

	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() == reflect.Ptr {
		srcValue = reflect.Indirect(srcValue)
	}

	if toValue.Elem().Type() != srcValue.Type() {
		return errors.Errorf("invalid type %s, should %s", toValue.Elem().Type(), srcValue.Type())
	}

	toValue.Elem().Set(srcValue)

	return nil
}

func (c *Cache) Set(_ context.Context, key string, data interface{}, dur time.Duration) error {
	c.cache.Set(key, data, dur)
	return nil
}

func (c *Cache) Get(_ context.Context, key string, data interface{}) error {
	d, ok := c.cache.Get(key)
	if !ok {
		return cache.ErrCacheMiss
	}

	return copyValue(d, data)
}

func (c *Cache) SetString(ctx context.Context, key string, data string, dur time.Duration) error {
	return c.Set(ctx, key, data, dur)
}

func (c *Cache) GetString(_ context.Context, key string) (string, error) {
	data, ok := c.cache.Get(key)
	if !ok {
		return "", cache.ErrCacheMiss
	}
	ds, ok := data.(string)
	if !ok {
		return fmt.Sprintf("%v", ds), nil
	}
	return ds, nil
}

func (c *Cache) UnSet(_ context.Context, key string) error {
	c.cache.Delete(key)
	return nil
}

func (c *Cache) Exist(_ context.Context, key string) (bool, error) {
	_, ok := c.cache.Get(key)
	return ok, nil
}

func (c *Cache) Lock(ctx context.Context, name string, dur time.Duration) (cache.Locker, error) {
	if dur == 0 {
		dur = time.Second
	}

	var k = lockKey(name)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		c.mu.Lock()
		_, ok := c.cache.Get(k)
		if !ok {
			c.cache.Set(k, struct{}{}, dur)
			c.mu.Unlock()

			lm := &sync.Mutex{}
			once := &sync.Once{}
			lm.Lock()

			return &locker{
				key:   k,
				o:     once,
				m:     lm,
				cache: c.cache,
			}, nil
		}
		c.mu.Unlock()
		time.Sleep(time.Microsecond * 100)
	}
}

type locker struct {
	key   string
	o     *sync.Once
	m     *sync.Mutex
	cache *memCache.Cache
}

func (l *locker) TTL(ctx context.Context) (time.Duration, error) {
	_, ok := l.cache.Get(l.key)
	if !ok {
		return time.Second, nil
	}
	return 0, nil
}

func (l *locker) Refresh(ctx context.Context, ttl time.Duration) error {
	if ttl == 0 {
		ttl = time.Second
	}
	_, ok := l.cache.Get(l.key)
	if !ok {
		return cache.ErrUnlock
	}
	l.cache.Set(l.key, struct{}{}, ttl)

	return nil
}

func (l *locker) Release(ctx context.Context) error {
	unlock := false
	var f = func() {
		l.m.Unlock()
		unlock = true
	}

	l.o.Do(f)

	if unlock {
		l.cache.Delete(l.key)
	}

	return nil
}
