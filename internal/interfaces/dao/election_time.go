package dao

import (
	"context"
	sql2 "database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/golib/txmanager/utils"
	"github.com/nocturna-ta/vote/internal/domain/model"
	"github.com/nocturna-ta/vote/internal/domain/repository"
	"time"
)

type ElectionTimeRepository struct {
	client ethereum.Client
	db     *sql.Store
}

type OptsElectionTimeRepository struct {
	Client ethereum.Client
	DB     *sql.Store
}

func NewElectionTimeRepository(opts *OptsElectionTimeRepository) repository.ElectionTimeRepository {
	return &ElectionTimeRepository{
		client: opts.Client,
		db:     opts.DB,
	}
}

const (
	insertElectionTimeQuery = `INSERT INTO election_times (id, title, description, start_time, end_time, status, is_active, created_at, updated_at, is_deleted) 
							   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	selectElectionTimeQuery = `SELECT %s FROM election_times %s WHERE TRUE %s`

	updateElectionTimeQuery = `UPDATE election_times SET %s WHERE TRUE %s`
)

func (e *ElectionTimeRepository) CreateElectionTime(ctx context.Context, electionTime *model.ElectionTime) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.CreateElectionTime")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		err error
	)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, insertElectionTimeQuery,
			electionTime.ID, electionTime.Title, electionTime.Description,
			electionTime.StartTime, electionTime.EndTime, electionTime.Status,
			electionTime.IsActive, electionTime.CreatedAt, electionTime.UpdatedAt,
			electionTime.IsDeleted)
	} else {
		_, err = e.db.GetMaster().ExecContext(ctx, insertElectionTimeQuery,
			electionTime.ID, electionTime.Title, electionTime.Description,
			electionTime.StartTime, electionTime.EndTime, electionTime.Status,
			electionTime.IsActive, electionTime.CreatedAt, electionTime.UpdatedAt,
			electionTime.IsDeleted)
	}

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				log.WithFields(log.Fields{
					"error":         err,
					"election_time": electionTime,
				}).ErrorWithCtx(ctx, "[ElectionTimeRepository.CreateElectionTime] Duplicate entry for election time")
				return ErrDuplicate
			}
		}
		log.WithFields(log.Fields{
			"error":         err,
			"election_time": electionTime,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.CreateElectionTime] failed to create election time")
		return err
	}

	return nil
}

func (e *ElectionTimeRepository) UpdateElectionTime(ctx context.Context, electionTime *model.ElectionTime) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.UpdateElectionTime")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		err    error
		result sql2.Result
		args   []any
	)

	setQuery := `title = $1, description = $2, start_time = $3, end_time = $4, status = $5, is_active = $6, updated_at = $7`
	whereQuery := `WHERE id = $8 AND is_deleted = false`
	args = append(args, electionTime.Title, electionTime.Description, electionTime.StartTime, electionTime.EndTime,
		electionTime.Status, electionTime.IsActive, time.Now(), electionTime.ID)

	query := fmt.Sprintf(updateElectionTimeQuery, setQuery, whereQuery)

	if sqlTrx != nil {
		result, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		result, err = e.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":         err,
			"election_time": electionTime,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.UpdateElectionTime] failed to update election time")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithFields(log.Fields{
			"error":         err,
			"election_time": electionTime,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.UpdateElectionTime] failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		return ErrNoUpdateHappened
	}

	return nil
}

func (e *ElectionTimeRepository) UpdateElectionTimeStatus(ctx context.Context, id uuid.UUID, status model.ElectionStatus, isActive bool) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.UpdateElectionTimeStatus")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		err    error
		result sql2.Result
		args   []any
	)

	setQuery := `status = $1, is_active = $2, updated_at = $3`
	whereQuery := ` AND id = $4 AND is_deleted = false`
	args = append(args, status, isActive, time.Now(), id)

	query := fmt.Sprintf(updateElectionTimeQuery, setQuery, whereQuery)

	if sqlTrx != nil {
		result, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		result, err = e.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"id":        id,
			"status":    status,
			"is_active": isActive,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.UpdateElectionTimeStatus] failed to update election time status")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"id":        id,
			"status":    status,
			"is_active": isActive,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.UpdateElectionTimeStatus] failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		return ErrNoUpdateHappened
	}

	return nil
}

func (e *ElectionTimeRepository) GetElectionTimeByID(ctx context.Context, id uuid.UUID) (*model.ElectionTime, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.GetElectionTimeByID")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		electionTime model.ElectionTime
		err          error
		args         []any
	)

	selectQuery := `id, title, description, start_time, end_time, status, is_active, created_at, updated_at, is_deleted`
	whereQuery := `AND id = $1 AND is_deleted = false`
	args = append(args, id)

	query := fmt.Sprintf(selectElectionTimeQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &electionTime, query, args...)
	} else {
		err = e.db.GetMaster().GetContext(ctx, &electionTime, query, args...)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			log.WithFields(log.Fields{
				"error": err,
				"id":    id,
			}).ErrorWithCtx(ctx, "[ElectionTimeRepository.GetElectionTimeByID] election time not found")
			return nil, ErrNoResult
		}
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.GetElectionTimeByID] failed to get election time by id")
		return nil, err
	}

	return &electionTime, nil
}

func (e *ElectionTimeRepository) GetActiveElectionTime(ctx context.Context) (*model.ElectionTime, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.GetActiveElectionTime")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		electionTime model.ElectionTime
		err          error
	)

	selectQuery := `id, title, description, start_time, end_time, status, is_active, created_at, updated_at, is_deleted`
	whereQuery := `AND is_active = true AND is_deleted = false`

	query := fmt.Sprintf(selectElectionTimeQuery, selectQuery, "", whereQuery)
	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &electionTime, query)
	} else {
		err = e.db.GetMaster().GetContext(ctx, &electionTime, query)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionTimeRepository.GetActiveElectionTime] no active election time found")
			return nil, ErrNoResult
		}
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.GetActiveElectionTime] failed to get active election time")
		return nil, err
	}

	return &electionTime, nil
}

func (e *ElectionTimeRepository) GetCurrentElectionTime(ctx context.Context) (*model.ElectionTime, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.GetCurrentElectionTime")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		electionTime model.ElectionTime
		err          error
		args         []any
	)

	selectQuery := `id, title, description, start_time, end_time, status, is_active, created_at, updated_at, is_deleted`
	whereQuery := ` AND start_time <= $1 AND end_time >= $1 AND is_deleted = false ORDER BY start_time DESC LIMIT 1`
	args = append(args, time.Now())

	query := fmt.Sprintf(selectElectionTimeQuery, selectQuery, "", whereQuery)
	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &electionTime, query, args...)
	} else {
		err = e.db.GetMaster().GetContext(ctx, &electionTime, query, args...)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionTimeRepository.GetCurrentElectionTime] no current election time found")
			return nil, ErrNoResult
		}
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.GetCurrentElectionTime] failed to get current election time")
		return nil, err
	}

	return &electionTime, nil
}

func (e *ElectionTimeRepository) GetAllElectionTimes(ctx context.Context, limit, offset int) ([]model.ElectionTime, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.GetAllElectionTimes")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		electionTimes []model.ElectionTime
		err           error
		args          []any
	)

	selectQuery := `id, title, description, start_time, end_time, status, is_active, created_at, updated_at, is_deleted`
	whereQuery := `AND is_deleted = false ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	args = append(args, limit, offset)

	query := fmt.Sprintf(selectElectionTimeQuery, selectQuery, "", whereQuery)
	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &electionTimes, query, args...)
	} else {
		err = e.db.GetMaster().SelectContext(ctx, &electionTimes, query, args...)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			log.WithFields(log.Fields{
				"error":  err,
				"limit":  limit,
				"offset": offset,
			}).ErrorWithCtx(ctx, "[ElectionTimeRepository.GetAllElectionTimes] no election times found")
			return nil, ErrNoResult
		}
		log.WithFields(log.Fields{
			"error":  err,
			"limit":  limit,
			"offset": offset,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.GetAllElectionTimes] failed to get all election times")
		return nil, err
	}

	return electionTimes, nil
}

func (e *ElectionTimeRepository) ActivateElectionTime(ctx context.Context, id uuid.UUID) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.ActivateElectionTime")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		err    error
		result sql2.Result
		args   []any
	)

	setQuery := `is_active = true, status = $1, updated_at = $2`
	whereQuery := ` AND id = $3 AND is_deleted = false`
	args = append(args, model.ElectionStatusActive, time.Now(), id)

	query := fmt.Sprintf(updateElectionTimeQuery, setQuery, whereQuery)

	if sqlTrx != nil {
		result, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		result, err = e.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.ActivateElectionTime] failed to activate election time")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.ActivateElectionTime] failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		return ErrNoUpdateHappened
	}

	return nil
}

func (e *ElectionTimeRepository) DeactivateAllElections(ctx context.Context) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.DeactivateAllElections")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		err  error
		args []any
	)

	setQuery := `is_active = false, updated_at = $1`
	whereQuery := `AND is_active = true AND is_deleted = false`

	args = append(args, time.Now())
	query := fmt.Sprintf(updateElectionTimeQuery, setQuery, whereQuery)
	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		_, err = e.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.DeactivateAllElections] failed to deactivate all elections")
		return err
	}

	return nil
}

func (e *ElectionTimeRepository) GetExpiredActiveElections(ctx context.Context) ([]*model.ElectionTime, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.GetExpiredActiveElections")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		electionTimes []*model.ElectionTime
		err           error
		args          []any
	)

	selectQuery := `id, title, description, start_time, end_time, status, is_active, created_at, updated_at, is_deleted`
	whereQuery := `AND is_active = true AND end_time < $1 AND is_deleted = false`
	args = append(args, time.Now())

	query := fmt.Sprintf(selectElectionTimeQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &electionTimes, query, args...)
	} else {
		err = e.db.GetMaster().SelectContext(ctx, &electionTimes, query, args...)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			return []*model.ElectionTime{}, nil
		}
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.GetExpiredActiveElections] failed to get expired active elections")
		return nil, err
	}

	return electionTimes, nil
}

func (e *ElectionTimeRepository) GetElectionsToActivate(ctx context.Context) ([]*model.ElectionTime, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.GetElectionsToActivate")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		electionTimes []*model.ElectionTime
		err           error
		args          []any
	)

	selectQuery := `id, title, description, start_time, end_time, status, is_active, created_at, updated_at, is_deleted`
	whereQuery := `AND is_active = false AND start_time <= $1 AND end_time >= $1 AND status = $2 AND is_deleted = false`
	args = append(args, time.Now(), model.ElectionStatusScheduled)
	query := fmt.Sprintf(selectElectionTimeQuery, selectQuery, "", whereQuery)
	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &electionTimes, query, args...)
	} else {
		err = e.db.GetMaster().SelectContext(ctx, &electionTimes, query, args...)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			return []*model.ElectionTime{}, nil
		}
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.GetElectionsToActivate] failed to get elections to activate")
		return nil, err
	}

	return electionTimes, nil
}

func (e *ElectionTimeRepository) DeleteElectionTime(ctx context.Context, id uuid.UUID) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeRepository.DeleteElectionTime")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		err    error
		result sql2.Result
		args   []any
	)

	setQuery := `is_deleted = true, updated_at = $1`
	whereQuery := `AND id = $2 AND is_deleted = false`
	args = append(args, time.Now(), id)

	query := fmt.Sprintf(updateElectionTimeQuery, setQuery, whereQuery)

	if sqlTrx != nil {
		result, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		result, err = e.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.DeleteElectionTime] failed to delete election time")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionTimeRepository.DeleteElectionTime] failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		return ErrNoUpdateHappened
	}

	return nil
}
