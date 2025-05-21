package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/vote/internal/domain/model"
)

type VoteRepository interface {
	InsertVote(ctx context.Context, vote *model.Vote) error
	InsertVoteBlockchain(ctx context.Context, signedTransaction string) (string, error)
	GetVoteByID(ctx context.Context, id uuid.UUID) (*model.Vote, error)
	GetVoteByElectionPairID(ctx context.Context, electionPairID uuid.UUID) (*model.Vote, error)
	UpdateVoteStatus(ctx context.Context, id uuid.UUID, status model.VoteStatus) error
	UpdateVote(ctx context.Context, vote *model.Vote) error
	GetPendingVotes(ctx context.Context, limit, offset int) ([]*model.Vote, error)
	GetFailedVotes(ctx context.Context, limit, offset int) ([]*model.Vote, error)
}
