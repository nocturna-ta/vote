package model

import (
	"github.com/google/uuid"
	"github.com/nocturna-ta/vote/internal/usecases/request"
	"time"
)

type VoteStatus string

const (
	VoteStatusPending  VoteStatus = "pending"
	VoteStatuConfirmed VoteStatus = "confirmed"
	VoteStatusRejected VoteStatus = "rejected"
	VoteStatusError    VoteStatus = "error"
)

func ToStringStatus(status VoteStatus) string {
	switch status {
	case VoteStatusPending:
		return "pending"
	case VoteStatuConfirmed:
		return "confirmed"
	case VoteStatusRejected:
		return "rejected"
	case VoteStatusError:
		return "error"
	default:
		return "unknown"
	}
}

type Vote struct {
	BaseModel
	ID              uuid.UUID  `db:"id"`
	VoterID         uuid.UUID  `db:"voter_id"`
	ElectionPairID  uuid.UUID  `db:"election_pair_id"`
	VotedAt         time.Time  `db:"voted_at"`
	Status          VoteStatus `db:"status"`
	TransactionHash string     `db:"transaction_hash"`
	Region          string     `db:"region"`
}

func ConstructCastVote(req *request.CastVoteRequest) *Vote {
	vote := &Vote{
		BaseModel: BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:             uuid.New(),
		VoterID:        uuid.MustParse(req.VoterID),
		ElectionPairID: uuid.MustParse(req.ElectionPairID),
		VotedAt:        time.Now(),
		Status:         VoteStatusPending,
		Region:         req.Region,
	}

	return vote
}
