package request

type CastVoteRequest struct {
	VoterID           string `json:"voter_id"`
	ElectionPairID    string `json:"election_pair_id"`
	Region            string `json:"region"`
	SignedTransaction string `json:"signed_transaction"`
}
