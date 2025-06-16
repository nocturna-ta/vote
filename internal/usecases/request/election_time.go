package request

import (
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"time"
)

type CreateElectionTimeRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type UpdateElectionTimeRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type GetElectionTimeRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ActivateElectionTimeRequest struct {
	ID string `json:"id"`
}

func (req *CreateElectionTimeRequest) Validate() error {
	if req.EndTime.Before(req.StartTime) || req.EndTime.Equal(req.StartTime) {
		return &custerr.ErrChain{
			Message: "End time must be after start time",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if req.StartTime.Before(time.Now()) {
		return &custerr.ErrChain{
			Message: "Start time cannot be in the past",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if req.Title == "" {
		return &custerr.ErrChain{
			Message: "Title is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

func (req *UpdateElectionTimeRequest) Validate() error {
	if req.EndTime.Before(req.StartTime) || req.EndTime.Equal(req.StartTime) {
		return &custerr.ErrChain{
			Message: "End time must be after start time",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}
