package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	event2 "github.com/nocturna-ta/common-model/models/event"
	libCtx "github.com/nocturna-ta/golib/context"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/vote/internal/domain/model"
	"github.com/nocturna-ta/vote/internal/interfaces/dao"
	"github.com/nocturna-ta/vote/pkg/constants"
	"time"
)

func (m *Module) ConsumeVoteSubmit(ctx context.Context, message *event.EventConsumeMessage) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteConsumer.ConsumeVoteSubmit")
	defer span.End()

	requestId := libCtx.ReadRequestId(ctx)
	log.WithFields(log.Fields{
		"request_id": requestId,
		"topic":      message.Topic,
		"key":        message.Key,
	}).InfoWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] consume message")

	var voteMessage event2.VoteSubmitMessage
	err := json.Unmarshal(message.Data, &voteMessage)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"topic":      message.Topic,
			"data":       string(message.Data),
			"request_id": requestId,
		}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] failed to unmarshal message")
		return nil
	}

	operation, ok := message.Metadata[constants.MetaDataOperation].(string)
	if !ok {
		log.WithFields(log.Fields{
			"topic":      message.Topic,
			"request_id": requestId,
			"metadata":   message.Metadata,
		}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] failed to get operation from metadata")
		operation = constants.Create
	}
	
	if operation != constants.Create {
		log.WithFields(log.Fields{
			"topic":      message.Topic,
			"request_id": requestId,
			"operation":  operation,
		}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] invalid operation")
		return nil
	}

	voteID, err := uuid.Parse(voteMessage.VoteID)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"request_id": requestId,
			"vote_id":    voteMessage.VoteID,
		}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] failed to parse vote id")
		return nil
	}

	vote, err := m.voteRepo.GetVoteByID(ctx, voteID)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			log.WithFields(log.Fields{
				"request_id": requestId,
				"vote_id":    voteID,
			}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] vote not found")
			return nil
		}
		log.WithFields(log.Fields{
			"error":      err,
			"request_id": requestId,
			"vote_id":    voteID,
		}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] failed to get vote by id")
		return err
	}

	if vote.Status != model.VoteStatusQueued && vote.Status != model.VoteStatusRetrying {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"vote_id":    voteID,
			"status":     vote.Status,
		}).InfoWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] vote already processed, skipping")
		return nil
	}

	err = m.voteRepo.UpdateVoteStatus(ctx, voteID, model.VoteStatusPending)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"request_id": requestId,
			"vote_id":    voteID,
		}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] failed to update vote status")
		return err
	}

	txHash, err := m.processVoteTx(ctx, voteMessage.SignedTransaction)
	now := time.Now()

	if err != nil {
		retryCount := vote.RetryCount + 1
		errorMsg := err.Error()

		updateErr := m.voteRepo.UpdateVote(ctx, &model.Vote{
			ID:           voteID,
			Status:       model.VoteStatusError,
			ErrorMessage: errorMsg,
			RetryCount:   retryCount,
			ProcessedAt:  &now,
		})

		if updateErr != nil {
			log.WithFields(log.Fields{
				"error":      updateErr,
				"request_id": requestId,
				"vote_id":    voteID,
			}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] failed to update vote status")
		}

		if retryCount <= m.maxRetries {
			log.WithFields(log.Fields{
				"request_id": requestId,
				"vote_id":    voteID,
				"retry":      retryCount,
				"max_retry":  m.maxRetries,
			}).InfoWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] retrying vote transaction")

			_ = m.voteRepo.UpdateVoteStatus(ctx, voteID, model.VoteStatusRetrying)

			return err
		}

		log.WithFields(log.Fields{
			"request_id": requestId,
			"vote_id":    voteID,
			"error":      err,
			"retry":      retryCount,
			"max_retry":  m.maxRetries,
		}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] failed to process vote transaction")

		if m.dlqTopic != "" {
			dlqMessage := vote.ToProcessedMessageModel(errorMsg)
			publishErr := m.publisher.Publish(ctx, m.dlqTopic, vote.ID.String(), dlqMessage, nil)
			if publishErr != nil {
				log.WithFields(log.Fields{
					"error":      publishErr,
					"request_id": requestId,
					"vote_id":    voteID,
				}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] failed to publish DLQ message")
			}
		}
		return nil
	}

	updateErr := m.voteRepo.UpdateVote(ctx, &model.Vote{
		ID:              voteID,
		Status:          model.VoteStatuConfirmed,
		TransactionHash: txHash,
		ProcessedAt:     &now,
		RetryCount:      vote.RetryCount,
	})

	if updateErr != nil {
		log.WithFields(log.Fields{
			"error":      updateErr,
			"request_id": requestId,
			"vote_id":    voteID,
		}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] failed to update vote status")
		return updateErr
	}

	if m.processedTopic != "" {
		processedMessage := vote.ToProcessedMessageModel("")
		processedMessage.TransactionHash = txHash
		processedMessage.Status = model.ToStringStatus(model.VoteStatuConfirmed)

		publishErr := m.publisher.Publish(ctx, m.processedTopic, vote.ID.String(), processedMessage, nil)
		if publishErr != nil {
			log.WithFields(log.Fields{
				"error":      publishErr,
				"request_id": requestId,
				"vote_id":    voteID,
			}).ErrorWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] failed to publish processed message")
		}
	}

	log.WithFields(log.Fields{
		"request_id": requestId,
		"vote_id":    voteID,
		"tx_hash":    txHash,
	}).InfoWithCtx(ctx, "[VoteConsumer.ConsumeVoteSubmit] vote transaction processed successfully")
	return nil
}

func (m *Module) processVoteTx(ctx context.Context, signedTransaction string) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteConsumer.processVoteTx")
	defer span.End()

	txHash, err := m.voteRepo.InsertVoteBlockchain(ctx, signedTransaction)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[VoteConsumer.processVoteTx] failed to insert vote blockchain")
		return "", err
	}

	return txHash, nil
}
