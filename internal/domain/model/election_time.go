package model

import (
	"github.com/google/uuid"
	"github.com/nocturna-ta/vote/internal/usecases/request"
	"time"
)

type ElectionStatus string

const (
	ElectionStatusScheduled ElectionStatus = "scheduled"
	ElectionStatusActive    ElectionStatus = "active"
	ElectionStatusEnded     ElectionStatus = "ended"
)

func ToStringElectionStatus(status ElectionStatus) string {
	switch status {
	case ElectionStatusScheduled:
		return "scheduled"
	case ElectionStatusActive:
		return "active"
	case ElectionStatusEnded:
		return "ended"
	default:
		return "unknown"
	}
}

func ParseElectionStatus(status string) ElectionStatus {
	switch status {
	case "scheduled":
		return ElectionStatusScheduled
	case "active":
		return ElectionStatusActive
	case "ended":
		return ElectionStatusEnded
	default:
		return ElectionStatusScheduled
	}
}

type ElectionTime struct {
	BaseModel
	ID          uuid.UUID      `db:"id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	StartTime   time.Time      `db:"start_time"`
	EndTime     time.Time      `db:"end_time"`
	Status      ElectionStatus `db:"status"`
	IsActive    bool           `db:"is_active"`
}

func ConstructCreateElection(req *request.CreateElectionTimeRequest) *ElectionTime {
	now := time.Now()

	electionTime := &ElectionTime{
		BaseModel: BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
			IsDeleted: false,
		},
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      ElectionStatusScheduled,
		IsActive:    false,
	}

	return electionTime
}

func (e *ElectionTime) IsElectionPeriodActive() bool {
	now := time.Now()
	return now.After(e.StartTime) && now.Before(e.EndTime) && e.IsDeleted == false
}

func (e *ElectionTime) ShouldBeActive() bool {
	now := time.Now()
	return now.After(e.StartTime) && now.Before(e.EndTime)
}

func (e *ElectionTime) ShouldBeEnded() bool {
	now := time.Now()
	return now.After(e.EndTime)
}

func (e *ElectionTime) UpdateStatusBasedOnTime() {
	now := time.Now()

	if now.Before(e.StartTime) {
		e.Status = ElectionStatusScheduled
		e.IsActive = false
	} else if now.After(e.StartTime) && now.Before(e.EndTime) {
		e.Status = ElectionStatusActive
		e.IsActive = true
	} else {
		e.Status = ElectionStatusEnded
		e.IsActive = false
	}
}
