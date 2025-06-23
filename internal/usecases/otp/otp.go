package otp

import (
	"context"
	"fmt"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/log"
	response2 "github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/vote/internal/domain/model"
	"github.com/nocturna-ta/vote/internal/usecases/request"
	"github.com/nocturna-ta/vote/internal/usecases/response"
	"time"
)

func (m *Module) GenerateOTP(ctx context.Context, req *request.GenerateOTPRequest) (*response.GenerateOTPResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "OTPUseCases.GenerateOTP")
	defer span.End()

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if !m.otpConfig.Enabled {
		return nil, &custerr.ErrChain{
			Message: "OTP service is disabled",
			Code:    503,
			Type:    response2.ErrInternalServerError,
		}
	}

	if m.smsConfig.Enabled && req.PhoneNumber != "" {
		if err := m.smsProvider.ValidatePhoneNumber(req.PhoneNumber); err != nil {
			return nil, &custerr.ErrChain{
				Message: fmt.Sprintf("Invalid phone number: %s", err.Error()),
				Code:    400,
				Type:    response2.ErrBadRequest,
			}
		}
	}

	lockKey := fmt.Sprintf("otp:generate:%s:%s", req.VoterID, req.Purpose)

	err := m.redLock.AcquireLockWithTTL(ctx, lockKey, 30*time.Second, 1)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"voter_id": req.VoterID,
			"purpose":  req.Purpose,
		}).WarnWithCtx(ctx, "[GenerateOTP] Failed to acquire lock")
		return nil, &custerr.ErrChain{
			Message: "OTP generation in progress, please wait",
			Code:    429,
			Type:    response2.ErrConflict,
		}
	}
	defer m.redLock.ReleaseLock(ctx, lockKey)

	// Check if there's an existing valid OTP
	existingOTP, err := m.getOTPFromRedis(ctx, req.VoterID, req.Purpose)
	if err == nil && existingOTP.IsValid() {
		remainingTime := existingOTP.GetTimeRemaining()
		return &response.GenerateOTPResponse{
			VoterID:           req.VoterID,
			Purpose:           req.Purpose,
			ExpiresAt:         existingOTP.ExpiresAt,
			TimeRemaining:     remainingTime.String(),
			Message:           "OTP already exists and is still valid",
			MaxAttempts:       existingOTP.MaxAttempts,
			RemainingAttempts: existingOTP.MaxAttempts - existingOTP.AttemptCount,
		}, nil
	}

	// Generate new OTP
	otp, err := model.NewOTP(req.VoterID, req.Purpose, m.otpConfig.TTL, m.otpConfig.MaxRetries)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"voter_id": req.VoterID,
			"purpose":  req.Purpose,
		}).ErrorWithCtx(ctx, "[GenerateOTP] Failed to create new OTP")
		return nil, &custerr.ErrChain{
			Message: "Failed to generate OTP",
			Code:    500,
			Type:    response2.ErrInternalServerError,
		}
	}

	// Store in Redis
	err = m.storeOTPInRedis(ctx, otp)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"voter_id": req.VoterID,
			"purpose":  req.Purpose,
		}).ErrorWithCtx(ctx, "[GenerateOTP] Failed to store OTP in Redis")
		return nil, &custerr.ErrChain{
			Message: "Failed to store OTP",
			Code:    500,
			Type:    response2.ErrInternalServerError,
		}
	}

	if m.smsConfig.Enabled && req.PhoneNumber != "" {
		if err = m.sendOTPSMS(ctx, req.PhoneNumber, otp.Code); err != nil {
			log.WithFields(log.Fields{
				"error":        err,
				"voter_id":     req.VoterID,
				"purpose":      req.Purpose,
				"phone_number": req.PhoneNumber,
			}).ErrorWithCtx(ctx, "[GenerateOTP] Failed to send OTP via SMS")
		} else {
			log.WithFields(log.Fields{
				"voter_id":     req.VoterID,
				"purpose":      req.Purpose,
				"phone_number": req.PhoneNumber,
			}).InfoWithCtx(ctx, "[GenerateOTP] OTP sent via SMS successfully")
		}
	}

	log.WithFields(log.Fields{
		"voter_id":   req.VoterID,
		"purpose":    req.Purpose,
		"code":       otp.Code, // Remove this in production
		"expires_at": otp.ExpiresAt,
	}).InfoWithCtx(ctx, "[GenerateOTP] OTP generated successfully")

	timeremaining := otp.GetTimeRemaining()

	return &response.GenerateOTPResponse{
		VoterID:           req.VoterID,
		Purpose:           req.Purpose,
		ExpiresAt:         otp.ExpiresAt,
		TimeRemaining:     timeremaining.String(),
		Message:           "OTP generated successfully",
		MaxAttempts:       otp.MaxAttempts,
		RemainingAttempts: otp.MaxAttempts,
	}, nil
}

func (m *Module) VerifyOTP(ctx context.Context, req *request.VerifyOTPRequest) (*response.VerifyOTPResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "OTPUseCases.VerifyOTP")
	defer span.End()

	if err := req.Validate(); err != nil {
		return nil, err
	}

	if !m.otpConfig.Enabled {
		return &response.VerifyOTPResponse{
			VoterID: req.VoterID,
			Purpose: req.Purpose,
			IsValid: true,
			Message: "OTP verification skipped (OTP disabled)",
		}, nil
	}

	otp, err := m.getOTPFromRedis(ctx, req.VoterID, req.Purpose)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"voter_id": req.VoterID,
			"purpose":  req.Purpose,
		}).WarnWithCtx(ctx, "[VerifyOTP] OTP not found")
		return &response.VerifyOTPResponse{
			VoterID: req.VoterID,
			Purpose: req.Purpose,
			IsValid: false,
			Message: "OTP not found or expired",
		}, nil
	}

	if !otp.IsValid() {
		var message string
		if otp.IsExpired() {
			message = "OTP has expired"
		} else if otp.AttemptCount >= otp.MaxAttempts {
			message = "Maximum verification attempts exceeded"
		} else {
			message = "OTP is no longer valid"
		}

		return &response.VerifyOTPResponse{
			VoterID: req.VoterID,
			Purpose: req.Purpose,
			IsValid: false,
			Message: message,
		}, nil
	}
	otp.IncrementAttempt()
	isCodeValid := otp.Code == req.Code

	if isCodeValid {
		otp.MarkAsVerified()
		token, err := m.generateOTPToken(req.VoterID, req.Purpose)
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"voter_id": req.VoterID,
				"purpose":  req.Purpose,
			}).ErrorWithCtx(ctx, "[VerifyOTP] Failed to generate OTP token")
		}

		if token != "" {
			tokenKey := fmt.Sprintf("otp:token:%s:%s", req.VoterID, req.Purpose)
			tokenExpiry := 10 * time.Minute
			err = m.redis.Set(ctx, tokenKey, token, int(tokenExpiry))
			if err != nil {
				log.WithFields(log.Fields{
					"error":    err,
					"voter_id": req.VoterID,
					"purpose":  req.Purpose,
				}).WarnWithCtx(ctx, "[VerifyOTP] Failed to store OTP token")
			}
		}

		log.WithFields(log.Fields{
			"voter_id": req.VoterID,
			"purpose":  req.Purpose,
		}).InfoWithCtx(ctx, "[VerifyOTP] OTP verified successfully")

		return &response.VerifyOTPResponse{
			VoterID:     req.VoterID,
			Purpose:     req.Purpose,
			IsValid:     true,
			Message:     "OTP verified successfully",
			VerifiedAt:  otp.VerifiedAt.Format(time.RFC3339),
			OTPToken:    token,
			TokenExpiry: time.Now().Add(10 * time.Minute).Format(time.RFC3339),
		}, nil
	}

	err = m.storeOTPInRedis(ctx, otp)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"voter_id": req.VoterID,
			"purpose":  req.Purpose,
		}).WarnWithCtx(ctx, "[VerifyOTP] Failed to update OTP attempt count")
	}

	log.WithFields(log.Fields{
		"voter_id":      req.VoterID,
		"purpose":       req.Purpose,
		"attempt_count": otp.AttemptCount,
		"max_attempts":  otp.MaxAttempts,
	}).WarnWithCtx(ctx, "[VerifyOTP] Invalid OTP code provided")

	return &response.VerifyOTPResponse{
		VoterID: req.VoterID,
		Purpose: req.Purpose,
		IsValid: false,
		Message: fmt.Sprintf("Invalid OTP code. %d attempts remaining", otp.MaxAttempts-otp.AttemptCount),
	}, nil
}

func (m *Module) ResendOTP(ctx context.Context, req *request.ResendOTPRequest) (*response.ResendOTPResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "OTPUseCases.ResendOTP")
	defer span.End()

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if !m.otpConfig.Enabled {
		return nil, &custerr.ErrChain{
			Message: "OTP service is disabled",
			Code:    503,
			Type:    response2.ErrInternalServerError,
		}
	}

	genReq := &request.GenerateOTPRequest{
		VoterID:     req.VoterID,
		Purpose:     req.Purpose,
		PhoneNumber: req.PhoneNumber,
	}

	_ = m.deleteOTPFromRedis(ctx, req.VoterID, req.Purpose)

	genResp, err := m.GenerateOTP(ctx, genReq)
	if err != nil {
		return nil, err
	}

	return &response.ResendOTPResponse{
		VoterID:           genResp.VoterID,
		Purpose:           genResp.Purpose,
		ExpiresAt:         genResp.ExpiresAt,
		TimeRemaining:     genResp.TimeRemaining,
		Message:           "OTP resent successfully",
		MaxAttempts:       genResp.MaxAttempts,
		RemainingAttempts: genResp.RemainingAttempts,
	}, nil
}

func (m *Module) GetOTPStatus(ctx context.Context, voterID, purpose string) (*response.OTPStatusResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "OTPUseCases.GetOTPStatus")
	defer span.End()

	if !m.otpConfig.Enabled {
		return &response.OTPStatusResponse{
			VoterID:           voterID,
			Purpose:           purpose,
			Status:            "disabled",
			AttemptCount:      0,
			MaxAttempts:       0,
			RemainingAttempts: 0,
			CanResend:         false,
		}, nil
	}

	otp, err := m.getOTPFromRedis(ctx, voterID, purpose)
	if err != nil {
		return &response.OTPStatusResponse{
			VoterID:           voterID,
			Purpose:           purpose,
			Status:            "not_found",
			AttemptCount:      0,
			MaxAttempts:       m.otpConfig.MaxRetries,
			RemainingAttempts: m.otpConfig.MaxRetries,
			CanResend:         true,
		}, nil
	}

	timeRemaining := otp.GetTimeRemaining()
	timeRemaining2 := timeRemaining.String()

	return &response.OTPStatusResponse{
		VoterID:           voterID,
		Purpose:           purpose,
		Status:            string(otp.Status),
		ExpiresAt:         &otp.ExpiresAt,
		TimeRemaining:     &timeRemaining2,
		AttemptCount:      otp.AttemptCount,
		MaxAttempts:       otp.MaxAttempts,
		RemainingAttempts: otp.MaxAttempts - otp.AttemptCount,
		CanResend:         otp.IsExpired() || otp.AttemptCount >= otp.MaxAttempts,
	}, nil
}

func (m *Module) ValidateOTPToken(ctx context.Context, voterID, purpose, token string) (bool, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "OTPUseCases.ValidateOTPToken")
	defer span.End()

	if !m.otpConfig.Enabled {
		return true, nil
	}

	tokenKey := fmt.Sprintf("otp:token:%s:%s", voterID, purpose)

	var storedToken string
	err := m.redis.GetObject(ctx, tokenKey, &storedToken)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"voter_id": voterID,
			"purpose":  purpose,
		}).WarnWithCtx(ctx, "[ValidateOTPToken] Token not found or expired")
		return false, nil
	}

	isValid := storedToken == token

	if isValid {
		_ = m.redis.Delete(ctx, tokenKey)

		log.WithFields(log.Fields{
			"voter_id": voterID,
			"purpose":  purpose,
		}).InfoWithCtx(ctx, "[ValidateOTPToken] Token validated successfully")
	}

	return isValid, nil
}
