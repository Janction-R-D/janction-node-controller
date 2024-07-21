package utils

import (
	"context"
	"errors"
	"math"
	"time"
)

var (
	ErrExceedMaxRetryCount = errors.New("exceed max retry count")
)

type RetryDo[T any] func(ctx context.Context, retryCount int) (*T, bool, error)

type retryOpts struct {
	maxRetryCount      int
	intervalGrowFactor int
	maxInterval        int
}

func newRetryOpts() *retryOpts {
	return &retryOpts{
		maxRetryCount:      math.MaxInt,
		intervalGrowFactor: 2,
		maxInterval:        3600, // 1h
	}
}

type RetryOption func(*retryOpts)

func WithRetryMaxIntervalSecond(maxIntervalSecond int) RetryOption {
	return func(ro *retryOpts) { ro.maxInterval = maxIntervalSecond }
}

func WithRetryMaxCount(maxRetryCount int) RetryOption {
	return func(ro *retryOpts) { ro.maxRetryCount = maxRetryCount }
}

func WithRetryIntervalGrowFactor(intervalGrowFactor int) RetryOption {
	return func(ro *retryOpts) { ro.intervalGrowFactor = intervalGrowFactor }
}

func RetryWithTimeout[T any](timeout time.Duration, retryDo RetryDo[T], opts ...RetryOption) (*T, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return RetryWithContext(ctx, retryDo, opts...)
}

func RetryWithContext[T any](ctx context.Context, retryDo RetryDo[T], opts ...RetryOption) (*T, bool, error) {
	ro := newRetryOpts()
	for _, opt := range opts {
		opt(ro)
	}

	ticker := time.NewTicker(time.Second * 1)
	tick := 0
	interval := 1
	retryCount := 1
	for {
		tick++

		if tick%interval == 0 {
			tick = 0
			interval = Min(interval*ro.intervalGrowFactor, ro.maxInterval)

			v, ok, err := retryDo(ctx, retryCount)
			if err == nil && ok {
				return v, true, nil
			}

			if retryCount > ro.maxRetryCount {
				return nil, false, ErrExceedMaxRetryCount
			}
			retryCount++
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil, false, ctx.Err()
		}
	}
}
