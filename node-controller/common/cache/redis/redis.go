package redis

import (
	"node-controller/common/cache"
	"context"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/bsm/redislock"
	"github.com/pkg/errors"
)

const TYPE = "redis"

type Cache struct {
	redisCli   redis.Cmdable
	defaultTll time.Duration
}

func New(cli redis.Cmdable, defaultExpiration time.Duration) *Cache {
	return &Cache{
		redisCli:   cli,
		defaultTll: defaultExpiration,
	}
}

func (c *Cache) Type() string {
	return TYPE
}

func (c *Cache) Set(ctx context.Context, key string, data interface{}, dur time.Duration) error {
	if dur == cache.DefaultExpiration {
		dur = c.defaultTll
	}
	return c.redisCli.Set(ctx, key, data, dur).Err()
}

func (c *Cache) Get(ctx context.Context, key string, data interface{}) error {
	res := c.redisCli.Get(ctx, key)
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return cache.ErrCacheMiss
		}
		return err
	}

	return res.Scan(data)
}

func (c *Cache) SetString(ctx context.Context, key string, data string, dur time.Duration) error {
	return c.Set(ctx, key, data, dur)
}

func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	res := c.redisCli.Get(ctx, key)
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", cache.ErrCacheMiss
		}
		return "", err
	}

	return res.Val(), nil
}

func (c *Cache) UnSet(ctx context.Context, key string) error {
	return c.redisCli.Del(ctx, key).Err()
}

func (c *Cache) Exist(ctx context.Context, key string) (bool, error) {
	res := c.redisCli.Exists(ctx, key)
	if err := res.Err(); err != nil {
		return false, err
	}
	return res.Val() == 1, nil
}

func lockKey(key string) string {
	return "lock:" + key
}

func formatMs(dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}

func (c *Cache) Lock(ctx context.Context, name string, dur time.Duration) (cache.Locker, error) {
	if dur == 0 {
		dur = time.Second
	}
	k := lockKey(name)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		lock, err := redislock.Obtain(ctx, c.redisCli, k, dur, nil)
		if err != nil {
			if errors.Is(err, redislock.ErrNotObtained) {
				time.Sleep(time.Microsecond * 10)
				continue
			}
			return nil, err
		}
		return &locker{l: lock}, nil
	}
}

type locker struct {
	l *redislock.Lock
}

func (l *locker) TTL(ctx context.Context) (time.Duration, error) {
	return l.l.TTL(ctx)
}

func (l *locker) Refresh(ctx context.Context, ttl time.Duration) error {
	err := l.l.Refresh(ctx, ttl, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return cache.ErrUnlock
	}
	return err
}

func (l *locker) Release(ctx context.Context) error {
	return l.l.Release(ctx)
}
