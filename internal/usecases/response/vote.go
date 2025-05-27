package response

import "time"

type CastVoterResponse struct {
	ID      string `json:"id"`
	VotedAt string `json:"voted_at"`
	Status  string `json:"status"`
	TxHash  string `json:"tx_hash,omitempty"`
}

type VoteStatusResponse struct {
	ID           string     `json:"id"`
	Status       string     `json:"status"`
	TxHash       string     `json:"tx_hash,omitempty"`
	VotedAt      string     `json:"voted_at,omitempty"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
}
