package usecases

import (
	"context"
	"github.com/nocturna-ta/vote/internal/usecases/request"
	"github.com/nocturna-ta/vote/internal/usecases/response"
)

type OTPUseCases interface {
	GenerateOTP(ctx context.Context, req *request.GenerateOTPRequest) (*response.GenerateOTPResponse, error)
	VerifyOTP(ctx context.Context, req *request.VerifyOTPRequest) (*response.VerifyOTPResponse, error)
	ResendOTP(ctx context.Context, req *request.ResendOTPRequest) (*response.ResendOTPResponse, error)
	GetOTPStatus(ctx context.Context, voterID, purpose string) (*response.OTPStatusResponse, error)
	ValidateOTPToken(ctx context.Context, voterID, purpose, token string) (bool, error)
}
