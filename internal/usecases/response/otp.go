package response

import "time"

type GenerateOTPResponse struct {
	VoterID           string    `json:"voter_id"`
	Purpose           string    `json:"purpose"`
	ExpiresAt         time.Time `json:"expires_at"`
	TimeRemaining     string    `json:"time_remaining_seconds"`
	Message           string    `json:"message"`
	MaxAttempts       int       `json:"max_attempts"`
	RemainingAttempts int       `json:"remaining_attempts"`
}

type VerifyOTPResponse struct {
	VoterID     string `json:"voter_id"`
	Purpose     string `json:"purpose"`
	IsValid     bool   `json:"is_valid"`
	Message     string `json:"message"`
	VerifiedAt  string `json:"verified_at,omitempty"`
	OTPToken    string `json:"otp_token,omitempty"`
	TokenExpiry string `json:"token_expiry,omitempty"`
}

type ResendOTPResponse struct {
	VoterID           string    `json:"voter_id"`
	Purpose           string    `json:"purpose"`
	ExpiresAt         time.Time `json:"expires_at"`
	TimeRemaining     string    `json:"time_remaining_seconds"`
	Message           string    `json:"message"`
	MaxAttempts       int       `json:"max_attempts"`
	RemainingAttempts int       `json:"remaining_attempts"`
}

type OTPStatusResponse struct {
	VoterID           string     `json:"voter_id"`
	Purpose           string     `json:"purpose"`
	Status            string     `json:"status"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	TimeRemaining     *string    `json:"time_remaining_seconds,omitempty"`
	AttemptCount      int        `json:"attempt_count"`
	MaxAttempts       int        `json:"max_attempts"`
	RemainingAttempts int        `json:"remaining_attempts"`
	CanResend         bool       `json:"can_resend"`
}
