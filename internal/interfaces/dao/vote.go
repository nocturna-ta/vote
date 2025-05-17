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

func NewBlockchainRepository(opts *OptsVoteRepository) (repository.VoteRepository, error) {
	var contractInterface interfaces.ElectionManagerInterface
	contract, err := electionManager.NewElectionManager(opts.ContractAddress, opts.Client.GetEthClient())
	if err != nil {
		return nil, err
	}
	contractInterface = contract
	return &VoteRepository{
		client:   opts.Client,
		contract: contractInterface,
		db:       opts.DB,
	}, nil
}

const (
	insertVoteQuery = `INSERT INTO votes (id, election_pair_id, voted_at, status, transaction_hash,
    region, created_at, updated_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	selectVoteQuery       = `SELECT %s FROM votes %s WHERE TRUE %s`
	updateVoteStatusQuery = `UPDATE votes SET %s WHERE TRUE %s`
)

func (v *VoteRepository) InsertVote(ctx context.Context, vote *model.Vote) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteRepository.InsertVote")
	defer span.End()

	var err error

	sqlTrx := utils.GetSqlTx(ctx)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, insertVoteQuery, vote.ID, vote.ElectionPairID, vote.VotedAt,
			vote.Status, vote.TransactionHash, vote.Region, vote.CreatedAt, vote.UpdatedAt, vote.IsDeleted)
	} else {
		_, err = v.db.GetMaster().ExecContext(ctx, insertVoteQuery, vote.ID, vote.ElectionPairID, vote.VotedAt,
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

func (v *VoteRepository) GetVoteByID(ctx context.Context, id uuid.UUID) (*model.Vote, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteRepository.GetVoteByID")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		vote model.Vote
		err  error
		args []any
	)

	selectQuery := `SELECT id, election_pair_id, voted_at, status, transaction_hash, region, created_at, updated_at, is_deleted`
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

	selectQuery := `SELECT id, election_pair_id, voted_at, status, transaction_hash, region, created_at, updated_at, is_deleted`
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
