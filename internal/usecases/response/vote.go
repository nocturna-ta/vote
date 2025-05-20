package response

type CastVoterResponse struct {
	ID      string `json:"id"`
	VotedAt string `json:"voted_at"`
	Status  string `json:"status"`
	TxHash  string `json:"tx_hash"`
}
