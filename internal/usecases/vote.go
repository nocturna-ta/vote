package usecases

import (
	"context"
	"github.com/google/uuid"
)

type VoteUseCases interface {
	CastVote(ctx context.Context, ktp string, candidateID uuid.UUID) error
	GetVoteStatus(ctx context.Context, voterID uuid.UUID)
}
