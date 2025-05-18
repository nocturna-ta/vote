package cache

import (
	"context"
	"github.com/nocturna-ta/golib/cache/redlock"
	"github.com/nocturna-ta/golib/tracing"
	"time"
)

type DistributedLock interface {
	AcquireLock(ctx context.Context, key string, attempts int) error
	AcquireLockWithTTL(ctx context.Context, key string, ttl time.Duration, attempts int) error
	ReleaseLock(ctx context.Context, key string) error
}

type RedisLock struct {
	redLock redlock.RedLock
	prefix  string
}

func NewRedisLock(redLock redlock.RedLock, prefix string) DistributedLock {
	if prefix == "" {
		prefix = "vote:"
	}

	return &RedisLock{
		redLock: redLock,
		prefix:  prefix,
	}
}

func (r *RedisLock) AcquireLockWithTTL(ctx context.Context, key string, ttl time.Duration, attempts int) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisLock.AcquireLockWithTTL")
	defer span.End()

	fullKey := r.prefix + key
	return r.redLock.AcquireLockWithTTL(ctx, fullKey, ttl, attempts)
}

func (r *RedisLock) ReleaseLock(ctx context.Context, key string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisLock.ReleaseLock")
	defer span.End()

	fullKey := r.prefix + key
	return r.redLock.ReleaseLock(ctx, fullKey)
}

func (r *RedisLock) AcquireLock(ctx context.Context, key string, attempts int) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisLock.AcquireLock")
	defer span.End()

	fullKey := r.prefix + key
	return r.redLock.AcquireLock(ctx, fullKey, attempts)
}
