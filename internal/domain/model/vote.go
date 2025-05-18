package model

import (
	"github.com/google/uuid"
	"time"
)

type VoteStatus string

const (
	VoteStatusPending  VoteStatus = "pending"
	VoteStatuConfirmed VoteStatus = "confirmed"
	VoteStatusRejected VoteStatus = "rejected"
	VoteStatusError    VoteStatus = "error"
)

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

type Voter struct {
	BaseModel
	ID      uuid.UUID `db:"id"`
	KTP     string    `db:"ktp"`
	Name    string    `db:"name"`
	Address string    `db:"address"`
	Region  string    `db:"region"`
}
