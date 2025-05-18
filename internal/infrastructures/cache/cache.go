package cache

import (
	"context"
	"time"
)

type VoteCache interface {
	IsVoterEligible(ctx context.Context, ktp string) (bool, error)
	SetVoterEligibility(ctx context.Context, ktp string, eligible bool, expiration time.Duration) error

	GetVoteStatus(ctx context.Context, voterID string) (string, error)
	SetVoteStatus(ctx context.Context, voterID string, status string, expiration time.Duration) error

	GetVoterID(ctx context.Context, ktp string) (string, error)
	SetVoterID(ctx context.Context, ktp string, voterID string, expiration time.Duration) error

	GetElectionData(ctx context.Context, electionID string) ([]byte, error)
	SetElectionData(ctx context.Context, electionID string, data []byte, expiration time.Duration) error

	IncrementVoteCounter(ctx context.Context, counterKey string) (int64, error)
	GetVoteCounter(ctx context.Context, counterKey string) (int64, error)

	InvalidateVoterCache(ctx context.Context, ktp string) error
	InvalidateElectionCache(ctx context.Context, electionID string) error
}
