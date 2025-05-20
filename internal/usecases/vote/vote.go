package vote

import (
	"context"
	"errors"
	"github.com/nocturna-ta/golib/custerr"
	response2 "github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/vote/internal/domain/model"
	"github.com/nocturna-ta/vote/internal/interfaces/dao"
	"github.com/nocturna-ta/vote/internal/usecases/request"
	"github.com/nocturna-ta/vote/internal/usecases/response"
)

func (m *Module) CastVote(ctx context.Context, req *request.CastVoteRequest) (*response.CastVoterResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteUseCases.CastVote")
	defer span.End()

	var (
		vote   *model.Vote
		txHash string
	)

	transaction := func(txCtx context.Context) (any, error) {
		vote = model.ConstructCastVote(req)

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
		}

		txHash, errTx = m.voteRepo.InsertVoteBlockchain(txCtx, req.SignedTransaction)
		if errTx != nil {
			return nil, errTx
		}
		// example publish event
		//errTx = m.publisher.Publish(txCtx, m.topics.MasterDataParty.Value, updatedParty.ID.String(), updatedParty.ToMessageModel(), map[string]any{
		//	constants.MetaDataOperation: constants.Update,
		//})

		return nil, nil
	}

	_, err := m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		return nil, err
	}

	return &response.CastVoterResponse{
		ID:      vote.ID.String(),
		VotedAt: vote.VotedAt.String(),
		Status:  model.ToStringStatus(model.VoteStatuConfirmed),
		TxHash:  txHash,
	}, nil
}
