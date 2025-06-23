package request

import (
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"strings"
)

type GenerateOTPRequest struct {
	VoterID     string `json:"voter_id" validate:"required"`
	Purpose     string `json:"purpose" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type VerifyOTPRequest struct {
	VoterID string `json:"voter_id" validate:"required"`
	Purpose string `json:"purpose" validate:"required"`
	Code    string `json:"code" validate:"required"`
}

type ResendOTPRequest struct {
	VoterID     string `json:"voter_id" validate:"required"`
	Purpose     string `json:"purpose" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

func (req *GenerateOTPRequest) Validate() error {
	if strings.TrimSpace(req.VoterID) == "" {
		return &custerr.ErrChain{
			Message: "voter_id is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if strings.TrimSpace(req.Purpose) == "" {
		return &custerr.ErrChain{
			Message: "purpose is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	validPurposes := map[string]bool{
		"vote_cast": true,
	}

	if !validPurposes[req.Purpose] {
		return &custerr.ErrChain{
			Message: "invalid purpose",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if req.PhoneNumber != "" {
		if !strings.HasPrefix(req.PhoneNumber, "+") {
			return &custerr.ErrChain{
				Message: "phone number must be in international format (e.g., +628123456789)",
				Code:    400,
				Type:    response.ErrBadRequest,
			}
		}
	}

	return nil
}

func (req *VerifyOTPRequest) Validate() error {
	if strings.TrimSpace(req.VoterID) == "" {
		return &custerr.ErrChain{
			Message: "voter_id is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if strings.TrimSpace(req.Purpose) == "" {
		return &custerr.ErrChain{
			Message: "purpose is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if strings.TrimSpace(req.Code) == "" {
		return &custerr.ErrChain{
			Message: "code is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if len(req.Code) < 4 || len(req.Code) > 8 {
		return &custerr.ErrChain{
			Message: "invalid OTP code format",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

func (req *ResendOTPRequest) Validate() error {
	if strings.TrimSpace(req.VoterID) == "" {
		return &custerr.ErrChain{
			Message: "voter_id is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if strings.TrimSpace(req.Purpose) == "" {
		return &custerr.ErrChain{
			Message: "purpose is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}
