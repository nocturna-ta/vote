package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/vote/internal/domain/model"
)

type ElectionTimeRepository interface {
	CreateElectionTime(ctx context.Context, electionTime *model.ElectionTime) error
	UpdateElectionTime(ctx context.Context, electionTime *model.ElectionTime) error
	UpdateElectionTimeStatus(ctx context.Context, id uuid.UUID, status model.ElectionStatus, isActive bool) error

	GetElectionTimeByID(ctx context.Context, id uuid.UUID) (*model.ElectionTime, error)
	GetActiveElectionTime(ctx context.Context) (*model.ElectionTime, error)
	GetCurrentElectionTime(ctx context.Context) (*model.ElectionTime, error)
	GetAllElectionTimes(ctx context.Context, limit, offset int) ([]model.ElectionTime, error)

	DeactivateAllElections(ctx context.Context) error
	ActivateElectionTime(ctx context.Context, id uuid.UUID) error

	GetExpiredActiveElections(ctx context.Context) ([]*model.ElectionTime, error)
	GetElectionsToActivate(ctx context.Context) ([]*model.ElectionTime, error)

	DeleteElectionTime(ctx context.Context, id uuid.UUID) error
}
