package request

import (
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"strings"
)

type CastVoteRequest struct {
	VoterID           string `json:"voter_id"`
	ElectionPairID    string `json:"election_pair_id"`
	Region            string `json:"region"`
	SignedTransaction string `json:"signed_transaction"`
	OTPToken          string `json:"otp_token"`
}

func (req *CastVoteRequest) Validate() error {
	if strings.TrimSpace(req.VoterID) == "" {
		return &custerr.ErrChain{
			Message: "voter_id is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if strings.TrimSpace(req.ElectionPairID) == "" {
		return &custerr.ErrChain{
			Message: "election_pair_id is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if strings.TrimSpace(req.Region) == "" {
		return &custerr.ErrChain{
			Message: "region is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if strings.TrimSpace(req.SignedTransaction) == "" {
		return &custerr.ErrChain{
			Message: "signed_transaction is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if strings.TrimSpace(req.OTPToken) == "" {
		return &custerr.ErrChain{
			Message: "otp_token is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}
