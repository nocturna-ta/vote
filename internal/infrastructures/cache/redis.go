package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/nocturna-ta/golib/cache"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"time"
)

type RedisCache struct {
	client cache.Cache
	prefix string
}

func NewRedisCache(client cache.Cache, prefix string) VoteCache {
	if prefix == "" {
		prefix = "vote:"
	}

	return &RedisCache{
		client: client,
		prefix: prefix,
	}
}

func (r *RedisCache) IsVoterEligible(ctx context.Context, ktp string) (bool, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.IsVoterEligible")
	defer span.End()

	key := fmt.Sprintf("%seligible:%s", r.prefix, ktp)
	val, err := r.client.GetString(ctx, key)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return false, errors.New("voter eligibility not in cache")
		}
		return false, err
	}

	return val == "true", nil
}

func (r *RedisCache) SetVoterEligibility(ctx context.Context, ktp string, eligible bool, expiration time.Duration) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.SetVoterEligibility")
	defer span.End()

	key := fmt.Sprintf("%seligible:%s", r.prefix, ktp)
	val := "false"
	if eligible {
		val = "true"
	}

	err := r.client.Set(ctx, key, val, int(expiration.Seconds()))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ktp":   ktp,
		}).ErrorWithCtx(ctx, "[RedisCache.SetVoterEligibility] Failed to set voter eligibility")
		return err
	}

	return nil
}

func (r *RedisCache) GetVoteStatus(ctx context.Context, voterID string) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.GetVoteStatus")
	defer span.End()

	key := fmt.Sprintf("%sstatus:%s", r.prefix, voterID)
	val, err := r.client.GetString(ctx, key)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return "", errors.New("vote status not in cache")
		}
		return "", err
	}

	return val, nil
}

func (r *RedisCache) SetVoteStatus(ctx context.Context, voterID string, status string, expiration time.Duration) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.SetVoteStatus")
	defer span.End()

	key := fmt.Sprintf("%sstatus:%s", r.prefix, voterID)
	err := r.client.Set(ctx, key, status, int(expiration.Seconds()))
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"voterID": voterID,
		}).ErrorWithCtx(ctx, "[RedisCache.SetVoteStatus] Failed to set vote status")
		return err
	}

	return nil
}

func (r *RedisCache) GetVoterID(ctx context.Context, ktp string) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.GetVoterID")
	defer span.End()

	key := fmt.Sprintf("%svoterid:%s", r.prefix, ktp)
	val, err := r.client.GetString(ctx, key)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return "", errors.New("voter ID not in cache")
		}
		return "", err
	}

	return val, nil
}

func (r *RedisCache) SetVoterID(ctx context.Context, ktp string, voterID string, expiration time.Duration) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.SetVoterID")
	defer span.End()

	key := fmt.Sprintf("%svoterid:%s", r.prefix, ktp)
	err := r.client.Set(ctx, key, voterID, int(expiration.Seconds()))
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"ktp":     ktp,
			"voterID": voterID,
		}).ErrorWithCtx(ctx, "[RedisCache.SetVoterID] Failed to set voter ID")
		return err
	}

	return nil
}

func (r *RedisCache) GetElectionData(ctx context.Context, electionID string) ([]byte, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.GetElectionData")
	defer span.End()

	key := fmt.Sprintf("%selection:%s", r.prefix, electionID)
	val, err := r.client.Get(ctx, key)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return nil, errors.New("election data not in cache")
		}
		return nil, err
	}

	return val, nil
}

func (r *RedisCache) SetElectionData(ctx context.Context, electionID string, data []byte, expiration time.Duration) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.SetElectionData")
	defer span.End()

	key := fmt.Sprintf("%selection:%s", r.prefix, electionID)
	err := r.client.Set(ctx, key, data, int(expiration.Seconds()))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"electionID": electionID,
		}).ErrorWithCtx(ctx, "[RedisCache.SetElectionData] Failed to set election data")
		return err
	}

	return nil
}

func (r *RedisCache) IncrementVoteCounter(ctx context.Context, counterKey string) (int64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.IncrementVoteCounter")
	defer span.End()

	key := fmt.Sprintf("%scounter:%s", r.prefix, counterKey)
	val, err := r.client.Increment(ctx, key, 1)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).ErrorWithCtx(ctx, "[RedisCache.IncrementVoteCounter] Failed to increment vote counter")
		return 0, err
	}

	return val, nil
}

func (r *RedisCache) GetVoteCounter(ctx context.Context, counterKey string) (int64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.GetVoteCounter")
	defer span.End()

	key := fmt.Sprintf("%scounter:%s", r.prefix, counterKey)
	val, err := r.client.GetInt(ctx, key)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return 0, nil
		}
		return 0, err
	}

	return val, nil
}

func (r *RedisCache) InvalidateVoterCache(ctx context.Context, ktp string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.InvalidateVoterCache")
	defer span.End()

	key := fmt.Sprintf("%seligible:%s", r.prefix, ktp)
	err := r.client.Delete(ctx, key)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ktp":   ktp,
		}).ErrorWithCtx(ctx, "[RedisCache.InvalidateVoterCache] Failed to invalidate voter cache")
		return err
	}

	voterIDKey := fmt.Sprintf("%svoterid:%s", r.prefix, ktp)
	err = r.client.Delete(ctx, voterIDKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ktp":   ktp,
		}).ErrorWithCtx(ctx, "[RedisCache.InvalidateVoterCache] Failed to invalidate voter ID cache")
		return err
	}
	return nil
}

func (r *RedisCache) InvalidateElectionCache(ctx context.Context, electionID string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "RedisCache.InvalidateElectionCache")
	defer span.End()

	key := fmt.Sprintf("%selection:%s", r.prefix, electionID)
	err := r.client.Delete(ctx, key)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"electionID": electionID,
		}).ErrorWithCtx(ctx, "[RedisCache.InvalidateElectionCache] Failed to invalidate election cache")
		return err
	}

	return nil
}
