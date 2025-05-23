package vote

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/log"
	response2 "github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/vote/internal/domain/model"
	"github.com/nocturna-ta/vote/internal/interfaces/dao"
	"github.com/nocturna-ta/vote/internal/usecases/request"
	"github.com/nocturna-ta/vote/internal/usecases/response"
	"github.com/nocturna-ta/vote/pkg/constants"
)

func (m *Module) CastVote(ctx context.Context, req *request.CastVoteRequest) (*response.CastVoterResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteUseCases.CastVote")
	defer span.End()

	var (
		vote *model.Vote
	)

	transaction := func(txCtx context.Context) (any, error) {
		vote, err := model.ConstructCastVote(req, m.encryptor)
		if err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"request_id": ctx.Value("request_id"),
				"voter_id":   req.VoterID,
			}).ErrorWithCtx(ctx, "[CastVote] Failed to construct vote with encryption")
			return nil, err
		}

		errTx := m.voteRepo.InsertVote(txCtx, vote)
		if errTx != nil {
			if errors.Is(errTx, dao.ErrDuplicate) {
				return nil, &custerr.ErrChain{
					Message: "vote duplicated",
					Cause:   errTx,
					Code:    500,
					Type:    response2.ErrInternalServerError,
				}
			}
			return nil, errTx
		}

		return vote, nil
	}

	result, err := m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"request_id": ctx.Value("request_id"),
			"voter_id":   req.VoterID,
		}).ErrorWithCtx(ctx, "[CastVote] Failed to create vote")
		return nil, err
	}

	vote = result.(*model.Vote)

	voteMessage := vote.ToSubmitMessageModel(req.SignedTransaction)

	err = m.publisher.Publish(ctx, m.topics.VoteSubmitData.Value, vote.ID.String(), voteMessage, map[string]any{
		constants.MetaDataOperation: constants.Create,
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"vote_id":  vote.ID,
			"voter_id": req.VoterID,
		}).ErrorWithCtx(ctx, "[CastVote] Failed to publish vote message")
	}

	return &response.CastVoterResponse{
		ID:      vote.ID.String(),
		VotedAt: vote.VotedAt.String(),
		Status:  model.ToStringStatus(model.VoteStatuConfirmed),
	}, nil
}

func (m *Module) GetVoteStatus(ctx context.Context, voteID uuid.UUID) (*response.VoteStatusResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteUseCases.GetVoteStatus")
	defer span.End()

	vote, err := m.voteRepo.GetVoteByID(ctx, voteID)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "vote not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	return &response.VoteStatusResponse{
		ID:           vote.ID.String(),
		Status:       model.ToStringStatus(vote.Status),
		TxHash:       vote.TransactionHash,
		VotedAt:      vote.VotedAt.String(),
		ProcessedAt:  vote.ProcessedAt,
		ErrorMessage: vote.ErrorMessage,
	}, nil
}
