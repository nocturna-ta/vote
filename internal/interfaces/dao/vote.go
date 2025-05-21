package dao

import (
	"context"
	sql2 "database/sql"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/golib/txmanager/utils"
	"github.com/nocturna-ta/vote/internal/domain/model"
	"github.com/nocturna-ta/vote/internal/domain/repository"
	utils2 "github.com/nocturna-ta/vote/pkg/utils"
	"github.com/nocturna-ta/votechain-contract/binding/electionManager"
	"github.com/nocturna-ta/votechain-contract/interfaces"
	"time"
)

type VoteRepository struct {
	client   ethereum.Client
	contract interfaces.ElectionManagerInterface
	db       *sql.Store
}

type OptsVoteRepository struct {
	Client          ethereum.Client
	ContractAddress common.Address
	Contract        interfaces.ElectionManagerInterface
	DB              *sql.Store
}

func NewVoteRepository(opts *OptsVoteRepository) repository.VoteRepository {
	var contractInterface interfaces.ElectionManagerInterface
	contract, err := electionManager.NewElectionManager(opts.ContractAddress, opts.Client.GetEthClient())
	if err != nil {
		return nil
	}
	contractInterface = contract
	return &VoteRepository{
		client:   opts.Client,
		contract: contractInterface,
		db:       opts.DB,
	}
}

const (
	insertVoteQuery = `INSERT INTO votes (id, voter_id, election_pair_id, voted_at, status, transaction_hash,
    region, created_at, updated_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	selectVoteQuery       = `SELECT %s FROM votes %s WHERE TRUE %s`
	updateVoteStatusQuery = `UPDATE votes SET %s WHERE TRUE %s`
)

func (v *VoteRepository) InsertVote(ctx context.Context, vote *model.Vote) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteRepository.InsertVote")
	defer span.End()

	var err error

	sqlTrx := utils.GetSqlTx(ctx)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, insertVoteQuery, vote.ID, vote.VoterID, vote.ElectionPairID, vote.VotedAt,
			vote.Status, vote.TransactionHash, vote.Region, vote.CreatedAt, vote.UpdatedAt, vote.IsDeleted)
	} else {
		_, err = v.db.GetMaster().ExecContext(ctx, insertVoteQuery, vote.ID, vote.VoterID, vote.ElectionPairID, vote.VotedAt,
			vote.Status, vote.TransactionHash, vote.Region, vote.CreatedAt, vote.UpdatedAt, vote.IsDeleted)
	}

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				log.WithFields(log.Fields{
					"error": err,
					"vote":  vote,
				}).ErrorWithCtx(ctx, "[VoteRepository.InsertVote] Duplicate vote")
				return ErrDuplicate
			}
		}
		log.WithFields(log.Fields{
			"error": err,
			"vote":  vote,
		}).ErrorWithCtx(ctx, "[VoteRepository.InsertVote] Failed to insert vote")
		return err
	}

	return nil
}

func (v *VoteRepository) InsertVoteBlockchain(ctx context.Context, signedTransaction string) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteRepository.InsertVoteBlockchain")
	defer span.End()

	tx, err := utils2.StringToTx(signedTransaction)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[VoteRepository.InsertVoteBlockchain] Failed to convert signed transaction")
		return "", err
	}

	txHash, err := v.client.SendTransaction(ctx, tx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"tx":    tx,
		}).ErrorWithCtx(ctx, "[VoteRepository.InsertVoteBlockchain] Failed to send transaction")
		return "", err
	}

	return txHash, nil
}

func (v *VoteRepository) GetVoteByID(ctx context.Context, id uuid.UUID) (*model.Vote, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteRepository.GetVoteByID")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		vote model.Vote
		err  error
		args []any
	)

	selectQuery := `id, election_pair_id, voted_at, status, transaction_hash, region, created_at, updated_at, is_deleted`
	whereClause := `AND id = $1 AND is_deleted = false`
	joinQuery := ``
	args = append(args, id)

	query := fmt.Sprintf(selectVoteQuery, selectQuery, joinQuery, whereClause)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &vote, query, args...)
	} else {
		err = v.db.GetMaster().GetContext(ctx, &vote, query, args...)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			log.WithFields(log.Fields{
				"error": err,
				"id":    id,
			}).ErrorWithCtx(ctx, "[VoteRepository.GetVoteByID] No vote found")
			return nil, ErrNoResult
		}
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[VoteRepository.GetVoteByID] Failed to get vote by ID")
		return nil, err
	}

	return &vote, nil
}

func (v *VoteRepository) GetVoteByElectionPairID(ctx context.Context, electionPairID uuid.UUID) (*model.Vote, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteRepository.GetVoteByElectionPairID")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		vote model.Vote
		err  error
		args []any
	)

	selectQuery := `id, election_pair_id, voted_at, status, transaction_hash, region, created_at, updated_at, is_deleted`
	whereClause := `AND election_pair_id = $1 AND is_deleted = false`
	joinQuery := ``
	args = append(args, electionPairID)

	query := fmt.Sprintf(selectVoteQuery, selectQuery, joinQuery, whereClause)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &vote, query, args...)
	} else {
		err = v.db.GetMaster().GetContext(ctx, &vote, query, args...)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			log.WithFields(log.Fields{
				"error": err,
				"id":    electionPairID,
			}).ErrorWithCtx(ctx, "[VoteRepository.GetVoteByElectionPairID] No vote found")
			return nil, ErrNoResult
		}
		log.WithFields(log.Fields{
			"error": err,
			"id":    electionPairID,
		}).ErrorWithCtx(ctx, "[VoteRepository.GetVoteByElectionPairID] Failed to get vote by election pair ID")
		return nil, err
	}

	return &vote, nil
}

func (v *VoteRepository) UpdateVoteStatus(ctx context.Context, id uuid.UUID, status model.VoteStatus) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteRepository.UpdateVoteStatus")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		err    error
		result sql2.Result
		args   []any
	)

	setQuery := `status = $1, updated_at = $2`
	whereClause := `AND id = $3 AND is_deleted = false`
	args = append(args, status, time.Now(), id)

	query := fmt.Sprintf(updateVoteStatusQuery, setQuery, whereClause)

	if sqlTrx != nil {
		result, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		result, err = v.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"id":     id,
			"status": status,
		}).ErrorWithCtx(ctx, "[VoteRepository.UpdateVoteStatus] Failed to update vote status")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[VoteRepository.UpdateVoteStatus] Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		return ErrNoUpdateHappened
	}

	return nil
}
func (v *VoteRepository) UpdateVote(ctx context.Context, vote *model.Vote) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteRepository.UpdateVote")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		err    error
		args   []any
		result sql2.Result
	)

	setQuery := `status = $1, transaction_hash = $2, error_message = $3, retry_count = $4, processed_at = $5, updated_at = $6`
	whereQuery := `AND id = $7 AND is_deleted = false`
	args = append(args, vote.Status, vote.TransactionHash, vote.ErrorMessage, vote.RetryCount, vote.ProcessedAt, time.Now(), vote.ID)
	query := fmt.Sprintf(updateVoteStatusQuery, setQuery, whereQuery)

	if sqlTrx != nil {
		result, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		result, err = v.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"vote":  vote,
		}).ErrorWithCtx(ctx, "[VoteRepository.UpdateVote] Failed to update vote")
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"vote":  vote,
		}).ErrorWithCtx(ctx, "[VoteRepository.UpdateVote] Failed to get rows affected")
		return err
	}

	if rowAffected == 0 {
		return ErrNoUpdateHappened
	}

	return nil
}

func (v *VoteRepository) GetPendingVotes(ctx context.Context, limit, offset int) ([]*model.Vote, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteRepository.GetPendingVotes")
	defer span.End()

	var (
		err   error
		votes []*model.Vote
		args  []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `id, voter_id, election_pair_id, voted_at, status, transaction_hash, region, error_message, retry_count, processed_at, created_at, updated_at, is_deleted`
	whereQuery := ` AND status IN ($1, $2) AND is_deleted = false ORDER BY created_at ASC LIMIT $3 OFFSET $4`
	args = append(args, model.VoteStatusQueued, model.VoteStatusRetrying, limit, offset)

	query := fmt.Sprintf(selectVoteQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &votes, query, args...)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &votes, query, args...)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			log.WithFields(log.Fields{
				"error":  err,
				"limit":  limit,
				"offset": offset,
			}).ErrorWithCtx(ctx, "[VoteRepository.GetPendingVotes] No pending votes found")
			return nil, ErrNoResult
		}
		log.WithFields(log.Fields{
			"error":  err,
			"limit":  limit,
			"offset": offset,
		}).ErrorWithCtx(ctx, "[VoteRepository.GetPendingVotes] Failed to get pending votes")
		return nil, err
	}

	return votes, nil
}

func (v *VoteRepository) GetFailedVotes(ctx context.Context, limit, offset int) ([]*model.Vote, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteRepository.GetFailedVotes")
	defer span.End()

	var (
		err   error
		votes []*model.Vote
		args  []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `id, voter_id, election_pair_id, voted_at, status, transaction_hash, region, error_message, retry_count, processed_at, created_at, updated_at, is_deleted`
	whereQuery := ` AND status = ? AND is_deleted = false ORDER BY created_at ASC LIMIT ? OFFSET ?`
	args = append(args, model.VoteStatusError, limit, offset)

	query := fmt.Sprintf(selectVoteQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &votes, query, args...)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &votes, query, args...)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			log.WithFields(log.Fields{
				"error":  err,
				"limit":  limit,
				"offset": offset,
			}).ErrorWithCtx(ctx, "[VoteRepository.GetFailedVotes] No failed votes found")
			return nil, ErrNoResult
		}
		log.WithFields(log.Fields{
			"error":  err,
			"limit":  limit,
			"offset": offset,
		}).ErrorWithCtx(ctx, "[VoteRepository.GetFailedVotes] Failed to get failed votes")
		return nil, err
	}

	return votes, nil
}
