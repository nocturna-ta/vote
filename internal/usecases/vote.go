package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/vote/internal/usecases/request"
	"github.com/nocturna-ta/vote/internal/usecases/response"
)

type VoteUseCases interface {
	CastVote(ctx context.Context, req *request.CastVoteRequest) (*response.CastVoterResponse, error)
	GetVoteStatus(ctx context.Context, voteID uuid.UUID) (*response.VoteStatusResponse, error)
}
