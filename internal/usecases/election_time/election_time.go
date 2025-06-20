package election_time

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/log"
	response2 "github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/vote/internal/domain/model"
	"github.com/nocturna-ta/vote/internal/interfaces/dao"
	"github.com/nocturna-ta/vote/internal/usecases/request"
	"github.com/nocturna-ta/vote/internal/usecases/response"
	"github.com/nocturna-ta/vote/pkg/utils"
	"time"
)

func (m *Module) CreateElectionTime(ctx context.Context, req *request.CreateElectionTimeRequest) (*response.ElectionTimeResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeUseCases.CreateElectionTime")
	defer span.End()

	var electionTime *model.ElectionTime
	now := time.Now()
	transaction := func(txCtx context.Context) (any, error) {
		currentElection, err := m.electionTimeRepo.GetCurrentElectionTime(txCtx)
		if err != nil && !errors.Is(err, dao.ErrNoResult) {
			return nil, err
		}
		if currentElection != nil {
			if req.StartTime.Before(currentElection.EndTime) && req.EndTime.After(currentElection.StartTime) {
				return nil, &custerr.ErrChain{
					Message: "Election time overlaps with current election time",
					Code:    409,
					Type:    response2.ErrConflict,
				}
			}
		}

		electionTime = model.ConstructCreateElection(req)

		if req.StartTime.Before(now) || req.StartTime.Equal(now) {
			if req.EndTime.After(now) {
				electionTime.Status = model.ElectionStatusActive
				electionTime.IsActive = true

				log.WithFields(log.Fields{
					"election_id": electionTime.ID,
				}).InfoWithCtx(ctx, "[CreateElectionTime] Election is active immediately")
			}
		}

		err = m.electionTimeRepo.CreateElectionTime(txCtx, electionTime)
		if err != nil {
			if errors.Is(err, dao.ErrNoResult) {
				return nil, &custerr.ErrChain{
					Message: "Failed to create election time",
					Code:    500,
					Type:    response2.ErrInternalServerError,
				}
			}
			return nil, err
		}

		return nil, nil
	}

	_, err := m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		return nil, err
	}

	if electionTime.IsActive {
		log.WithFields(log.Fields{
			"election_id": electionTime.ID,
		}).InfoWithCtx(ctx, "[CreateElectionTime] Election is active immediately")
	} else {
		log.WithFields(log.Fields{
			"election_id": electionTime.ID,
			"start_time":  electionTime.StartTime,
		}).InfoWithCtx(ctx, "[CreateElectionTime] Election is scheduled to start at %s", electionTime.StartTime)
	}

	return &response.ElectionTimeResponse{
		ID:          electionTime.ID.String(),
		Title:       electionTime.Title,
		Description: electionTime.Description,
		StartTime:   electionTime.StartTime,
		EndTime:     electionTime.EndTime,
		Status:      model.ToStringElectionStatus(electionTime.Status),
		IsActive:    electionTime.IsActive,
		CreatedAt:   electionTime.CreatedAt,
		UpdatedAt:   electionTime.UpdatedAt,
	}, nil
}

func (m *Module) UpdateElectionTime(ctx context.Context, id uuid.UUID, req *request.UpdateElectionTimeRequest) (*response.ElectionTimeResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeUseCases.UpdateElectionTime")
	defer span.End()

	var electionTime *model.ElectionTime

	transaction := func(txCtx context.Context) (any, error) {
		existing, err := m.electionTimeRepo.GetElectionTimeByID(txCtx, id)
		if err != nil {
			if errors.Is(err, dao.ErrNoResult) {
				return nil, &custerr.ErrChain{
					Message: "Election time not found",
					Code:    404,
					Type:    response2.ErrNotFound,
				}
			}
			return nil, err
		}

		if existing.IsActive && time.Now().After(existing.StartTime) {
			return nil, &custerr.ErrChain{
				Message: "Cannot update active election that has already started",
				Code:    400,
				Type:    response2.ErrBadRequest,
			}
		}

		electionTime = existing
		electionTime.Title = req.Title
		electionTime.Description = req.Description
		electionTime.StartTime = req.StartTime
		electionTime.EndTime = req.EndTime
		electionTime.UpdatedAt = time.Now()

		electionTime.UpdateStatusBasedOnTime()

		err = m.electionTimeRepo.UpdateElectionTime(txCtx, electionTime)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err := m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"id":      id,
			"request": req,
		}).ErrorWithCtx(ctx, "[UpdateElectionTime] Failed to update election time")
		return nil, err
	}

	log.WithFields(log.Fields{
		"election_id": electionTime.ID,
		"title":       electionTime.Title,
	}).InfoWithCtx(ctx, "[UpdateElectionTime] Election time updated successfully")

	return &response.ElectionTimeResponse{
		ID:          electionTime.ID.String(),
		Title:       electionTime.Title,
		Description: electionTime.Description,
		StartTime:   electionTime.StartTime,
		EndTime:     electionTime.EndTime,
		Status:      model.ToStringElectionStatus(electionTime.Status),
		IsActive:    electionTime.IsActive,
		CreatedAt:   electionTime.CreatedAt,
		UpdatedAt:   electionTime.UpdatedAt,
	}, nil
}

func (m *Module) DeleteElectionTime(ctx context.Context, id uuid.UUID) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeUseCases.DeleteElectionTime")
	defer span.End()

	transaction := func(txCtx context.Context) (any, error) {
		existing, err := m.electionTimeRepo.GetElectionTimeByID(txCtx, id)
		if err != nil {
			if errors.Is(err, dao.ErrNoResult) {
				return nil, &custerr.ErrChain{
					Message: "Election time not found",
					Code:    404,
					Type:    response2.ErrNotFound,
				}
			}
			return nil, err
		}

		if existing.IsActive {
			return nil, &custerr.ErrChain{
				Message: "Cannot delete active election time",
				Code:    400,
				Type:    response2.ErrBadRequest,
			}
		}

		err = m.electionTimeRepo.DeleteElectionTime(txCtx, id)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[DeleteElectionTime] Failed to delete election time")
		return err
	}

	log.WithFields(log.Fields{
		"election_id": id,
	}).InfoWithCtx(ctx, "[DeleteElectionTime] Election time deleted successfully")

	return nil
}

func (m *Module) GetCurrentElectionStatus(ctx context.Context) (*response.ElectionStatusResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeUseCases.GetCurrentElectionStatus")
	defer span.End()

	activeElection, err := m.electionTimeRepo.GetActiveElectionTime(ctx)
	if err != nil && !errors.Is(err, dao.ErrNoResult) {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[GetCurrentElectionStatus] Failed to get active election")
		return nil, err
	}

	now := time.Now()

	if activeElection != nil {
		fmt.Println(activeElection.StartTime)
		fmt.Println(activeElection.EndTime)
		fmt.Println(now)
		fmt.Println(now.After(activeElection.StartTime) && now.Before(activeElection.EndTime) && activeElection.IsDeleted == false)
		if activeElection.IsElectionPeriodActive() {
			timeUntilEnd := int64(activeElection.EndTime.Sub(now).Seconds())
			remainingTime := utils.FormatDuration(activeElection.EndTime.Sub(now))

			return &response.ElectionStatusResponse{
				IsElectionActive: true,
				CurrentElection: &response.ElectionTimeResponse{
					ID:          activeElection.ID.String(),
					Title:       activeElection.Title,
					Description: activeElection.Description,
					StartTime:   activeElection.StartTime,
					EndTime:     activeElection.EndTime,
					Status:      model.ToStringElectionStatus(activeElection.Status),
					IsActive:    activeElection.IsActive,
					CreatedAt:   activeElection.CreatedAt,
					UpdatedAt:   activeElection.UpdatedAt,
				},
				Message:             "Election is currently active. Voting is open!",
				TimeUntilEnd:        &timeUntilEnd,
				RemainingVotingTime: remainingTime,
			}, nil
		}
	}

	allElections, err := m.electionTimeRepo.GetAllElectionTimes(ctx, 10, 0)
	if err != nil && !errors.Is(err, dao.ErrNoResult) {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[GetCurrentElectionStatus] Failed to get elections")
		return nil, err
	}

	var nextElection *model.ElectionTime
	for _, election := range allElections {
		if election.StartTime.After(now) && election.Status == model.ElectionStatusScheduled {
			if nextElection == nil || election.StartTime.Before(nextElection.StartTime) {
				nextElection = &election
			}
		}
	}

	if nextElection != nil {
		timeUntilStart := int64(nextElection.StartTime.Sub(now).Seconds())
		return &response.ElectionStatusResponse{
			IsElectionActive: false,
			CurrentElection: &response.ElectionTimeResponse{
				ID:          nextElection.ID.String(),
				Title:       nextElection.Title,
				Description: nextElection.Description,
				StartTime:   nextElection.StartTime,
				EndTime:     nextElection.EndTime,
				Status:      model.ToStringElectionStatus(nextElection.Status),
				IsActive:    nextElection.IsActive,
				CreatedAt:   nextElection.CreatedAt,
				UpdatedAt:   nextElection.UpdatedAt,
			},
			Message:        fmt.Sprintf("Election is scheduled to start at %s", nextElection.StartTime.Format("2006-01-02 15:04:05")),
			TimeUntilStart: &timeUntilStart,
		}, nil
	}

	return &response.ElectionStatusResponse{
		IsElectionActive: false,
		Message:          "No election is currently active or scheduled",
	}, nil
}

func (m *Module) GetElectionTimeByID(ctx context.Context, id uuid.UUID) (*response.ElectionTimeResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeUseCases.GetElectionTimeByID")
	defer span.End()

	electionTime, err := m.electionTimeRepo.GetElectionTimeByID(ctx, id)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "Election time not found",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[GetElectionTimeByID] Failed to get election time")
		return nil, err
	}

	return &response.ElectionTimeResponse{
		ID:          electionTime.ID.String(),
		Title:       electionTime.Title,
		Description: electionTime.Description,
		StartTime:   electionTime.StartTime,
		EndTime:     electionTime.EndTime,
		Status:      model.ToStringElectionStatus(electionTime.Status),
		IsActive:    electionTime.IsActive,
		CreatedAt:   electionTime.CreatedAt,
		UpdatedAt:   electionTime.UpdatedAt,
	}, nil

}

func (m *Module) ActivateElection(ctx context.Context, id uuid.UUID) (*response.ElectionTimeResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeUseCases.ActivateElection")
	defer span.End()

	var electionTime *model.ElectionTime

	transaction := func(txCtx context.Context) (any, error) {
		election, err := m.electionTimeRepo.GetElectionTimeByID(txCtx, id)
		if err != nil {
			if errors.Is(err, dao.ErrNoResult) {
				return nil, &custerr.ErrChain{
					Message: "Election time not found",
					Code:    404,
					Type:    response2.ErrNotFound,
				}
			}
			return nil, err
		}

		now := time.Now()
		if now.Before(election.StartTime) {
			return nil, &custerr.ErrChain{
				Message: "Cannot activate election before its start time",
				Code:    400,
				Type:    response2.ErrBadRequest,
			}
		}

		if now.After(election.EndTime) {
			return nil, &custerr.ErrChain{
				Message: "Cannot activate election after its end time",
				Code:    400,
				Type:    response2.ErrBadRequest,
			}
		}

		if election.IsActive {
			return nil, &custerr.ErrChain{
				Message: "Election is already active",
				Code:    400,
				Type:    response2.ErrBadRequest,
			}
		}

		err = m.electionTimeRepo.DeactivateAllElections(txCtx)
		if err != nil {
			return nil, err
		}

		err = m.electionTimeRepo.ActivateElectionTime(txCtx, id)
		if err != nil {
			return nil, err
		}

		electionTime, err = m.electionTimeRepo.GetElectionTimeByID(txCtx, id)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err := m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ActivateElection] Failed to activate election")
		return nil, err
	}

	log.WithFields(log.Fields{
		"election_id": electionTime.ID,
		"title":       electionTime.Title,
	}).InfoWithCtx(ctx, "[ActivateElection] Election activated successfully")

	return &response.ElectionTimeResponse{
		ID:          electionTime.ID.String(),
		Title:       electionTime.Title,
		Description: electionTime.Description,
		StartTime:   electionTime.StartTime,
		EndTime:     electionTime.EndTime,
		Status:      model.ToStringElectionStatus(electionTime.Status),
		IsActive:    electionTime.IsActive,
		CreatedAt:   electionTime.CreatedAt,
		UpdatedAt:   electionTime.UpdatedAt,
	}, nil
}

func (m *Module) SyncElectionStatuses(ctx context.Context) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeUseCases.SyncElectionStatuses")
	defer span.End()

	transaction := func(txCtx context.Context) (any, error) {
		expiredElections, err := m.electionTimeRepo.GetExpiredActiveElections(txCtx)
		if err != nil {
			return nil, err
		}

		for _, election := range expiredElections {
			err = m.electionTimeRepo.UpdateElectionTimeStatus(txCtx, election.ID, model.ElectionStatusEnded, false)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"id":    election.ID,
				}).ErrorWithCtx(ctx, "[SyncElectionStatuses] Failed to update election status")
				continue
			}
			log.WithFields(log.Fields{
				"election_id": election.ID,
				"title":       election.Title,
			}).InfoWithCtx(ctx, "[SyncElectionStatuses] Election status updated to ended")
		}

		electionToActivate, err := m.electionTimeRepo.GetElectionsToActivate(txCtx)
		if err != nil {
			return nil, err
		}

		for _, election := range electionToActivate {
			err := m.electionTimeRepo.DeactivateAllElections(txCtx)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).ErrorWithCtx(ctx, "[SyncElectionStatuses] Failed to deactivate all elections")
				continue
			}

			err = m.electionTimeRepo.ActivateElectionTime(txCtx, election.ID)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"id":    election.ID,
				}).ErrorWithCtx(ctx, "[SyncElectionStatuses] Failed to activate election")
				continue
			}
			log.WithFields(log.Fields{
				"election_id": election.ID,
				"title":       election.Title,
			}).InfoWithCtx(ctx, "[SyncElectionStatuses] Election activated successfully")
		}

		return nil, nil
	}

	_, err := m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[SyncElectionStatuses] Failed to sync election statuses")
		return err
	}

	return nil
}
