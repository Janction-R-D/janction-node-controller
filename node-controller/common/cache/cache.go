package cache

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

const (
	// NoExpiration For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// DefaultExpiration For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or
	// NewFrom() when the cache was created (e.g. 5 minutes.)
	DefaultExpiration time.Duration = 0
)

type Cache interface {
	// Type 缓存类型.
	Type() string

	Set(ctx context.Context, key string, data interface{}, dur time.Duration) error
	Get(ctx context.Context, key string, data interface{}) error

	SetString(ctx context.Context, key string, data string, dur time.Duration) error
	GetString(ctx context.Context, key string) (string, error)

	Exist(ctx context.Context, key string) (bool, error)

	UnSet(ctx context.Context, key string) error

	Lock(ctx context.Context, name string, dur time.Duration) (Locker, error)
}

type UnLock func() error

var ErrUnlock = errors.New("un lock")

type Locker interface {
	TTL(ctx context.Context) (time.Duration, error)
	Refresh(ctx context.Context, ttl time.Duration) error
	Release(ctx context.Context) error
}
