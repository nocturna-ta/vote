package model

import (
	"github.com/google/uuid"
	"github.com/nocturna-ta/common-model/models/event"
	"github.com/nocturna-ta/golib/utils/encryption"
	"github.com/nocturna-ta/vote/internal/usecases/request"
	"time"
)

type VoteStatus string

const (
	VoteStatusPending  VoteStatus = "pending"
	VoteStatuConfirmed VoteStatus = "confirmed"
	VoteStatusRejected VoteStatus = "rejected"
	VoteStatusError    VoteStatus = "error"
	VoteStatusQueued   VoteStatus = "queued"
	VoteStatusRetrying VoteStatus = "retrying"
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
	case VoteStatusQueued:
		return "queued"
	case VoteStatusRetrying:
		return "retrying"
	default:
		return "unknown"
	}
}

type Vote struct {
	BaseModel
	ID               uuid.UUID  `db:"id"`
	VoterID          uuid.UUID  `db:"-"`
	EncryptedVoterID string     `db:"voter_id"`
	ElectionPairID   uuid.UUID  `db:"election_pair_id"`
	VotedAt          time.Time  `db:"voted_at"`
	Status           VoteStatus `db:"status"`
	TransactionHash  string     `db:"transaction_hash"`
	Region           string     `db:"region"`
	ErrorMessage     string     `db:"error_message"`
	RetryCount       int        `db:"retry_count"`
	ProcessedAt      *time.Time `db:"processed_at"`
}

func (v *Vote) EncryptVoterID(encryptor *encryption.Encryption) error {
	if v.VoterID == uuid.Nil {
		return nil
	}

	encrypted, err := encryptor.Encrypt(v.VoterID.String())
	if err != nil {
		return err
	}

	v.EncryptedVoterID = encrypted
	return nil
}

func (v *Vote) DecryptVoterID(encryptor *encryption.Encryption) error {
	if v.EncryptedVoterID == "" {
		return nil
	}

	decrypted, err := encryptor.Decrypt(v.EncryptedVoterID)
	if err != nil {
		return err
	}

	voterID, err := uuid.Parse(decrypted)
	if err != nil {
		return err
	}

	v.VoterID = voterID
	return nil
}

func ConstructCastVote(req *request.CastVoteRequest, encryptor *encryption.Encryption) (*Vote, error) {
	voterID, err := uuid.Parse(req.VoterID)
	if err != nil {
		return nil, err
	}
	vote := &Vote{
		BaseModel: BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:             uuid.New(),
		VoterID:        voterID,
		ElectionPairID: uuid.MustParse(req.ElectionPairID),
		VotedAt:        time.Now(),
		Status:         VoteStatusQueued,
		Region:         req.Region,
		RetryCount:     0,
	}

	err = vote.EncryptVoterID(encryptor)
	if err != nil {
		return nil, err
	}
	return vote, nil
}

func (v *Vote) ToSubmitMessageModel(signedTx string) *event.VoteSubmitMessage {
	return &event.VoteSubmitMessage{
		BaseModelMessage: event.BaseModelMessage{
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			IsDeleted: v.IsDeleted,
		},
		VoteID:            v.ID.String(),
		VoterID:           v.EncryptedVoterID,
		ElectionPairID:    v.ElectionPairID.String(),
		Region:            v.Region,
		SignedTransaction: signedTx,
		SubmittedAt:       time.Now(),
	}
}

func (v *Vote) ToProcessedMessageModel(errorMsg string) *event.VoteProcessedMessage {
	return &event.VoteProcessedMessage{
		VoteID:          v.ID.String(),
		VoterID:         v.EncryptedVoterID,
		Status:          ToStringStatus(v.Status),
		TransactionHash: v.TransactionHash,
		ErrorMessage:    errorMsg,
		ProcessedAt:     time.Now(),
	}
}
