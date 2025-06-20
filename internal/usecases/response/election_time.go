package response

import "time"

type ElectionTimeResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ElectionStatusResponse struct {
	IsElectionActive    bool                  `json:"is_election_active"`
	CurrentElection     *ElectionTimeResponse `json:"current_election,omitempty"`
	Message             string                `json:"message,omitempty"`
	TimeUntilStart      *int64                `json:"time_until_start,omitempty"`
	TimeUntilEnd        *int64                `json:"time_until_end,omitempty"`
	RemainingVotingTime string                `json:"remaining_voting_time,omitempty"`
}
