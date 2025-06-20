package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/vote/internal/usecases/request"
	"github.com/nocturna-ta/vote/internal/usecases/response"
)

type ElectionTimeUseCases interface {
	CreateElectionTime(ctx context.Context, req *request.CreateElectionTimeRequest) (*response.ElectionTimeResponse, error)
	UpdateElectionTime(ctx context.Context, id uuid.UUID, req *request.UpdateElectionTimeRequest) (*response.ElectionTimeResponse, error)
	DeleteElectionTime(ctx context.Context, id uuid.UUID) error

	GetCurrentElectionStatus(ctx context.Context) (*response.ElectionStatusResponse, error)
	GetElectionTimeByID(ctx context.Context, id uuid.UUID) (*response.ElectionTimeResponse, error)

	SyncElectionStatuses(ctx context.Context) error
	ActivateElection(ctx context.Context, id uuid.UUID) (*response.ElectionTimeResponse, error)
}
