package usecases

import "context"

type VoteUseCases interface {
	CastVote(ctx context.Context, ktp string, candidateID uint64) error
}
