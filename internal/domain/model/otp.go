package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type OTPStatus string

const (
	OTPStatusPending  OTPStatus = "pending"
	OTPStatusVerified OTPStatus = "verified"
	OTPStatusExpired  OTPStatus = "expired"
	OTPStatusUsed     OTPStatus = "used"
)

type OTP struct {
	Code         string     `json:"code"`
	VoterID      string     `json:"voter_id"`
	Purpose      string     `json:"purpose"` // e.g., "vote_cast"
	Status       OTPStatus  `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiresAt    time.Time  `json:"expires_at"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	AttemptCount int        `json:"attempt_count"`
	MaxAttempts  int        `json:"max_attempts"`
}

func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		length = 4
	}

	digits := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		digits[i] = byte(num.Int64() + '0')
	}

	return string(digits), nil
}

func NewOTP(voterID, purpose string, ttl time.Duration, maxAttempts int) (*OTP, error) {
	code, err := GenerateOTP(4)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta") // WIB timezone
	now := time.Now().In(loc)

	return &OTP{
		Code:         code,
		VoterID:      voterID,
		Purpose:      purpose,
		Status:       OTPStatusPending,
		CreatedAt:    now,
		ExpiresAt:    now.Add(ttl),
		AttemptCount: 0,
		MaxAttempts:  maxAttempts,
	}, nil
}

func (o *OTP) IsExpired() bool {
	loc, _ := time.LoadLocation("Asia/Jakarta") // WIB timezone

	return time.Now().In(loc).After(o.ExpiresAt)
}

func (o *OTP) IsValid() bool {
	return o.Status == OTPStatusPending && !o.IsExpired() && o.AttemptCount < o.MaxAttempts
}

func (o *OTP) IncrementAttempt() {
	o.AttemptCount++
}

func (o *OTP) MarkAsVerified() {
	now := time.Now()
	o.Status = OTPStatusVerified
	o.VerifiedAt = &now
}

func (o *OTP) MarkAsUsed() {
	o.Status = OTPStatusUsed
}

func (o *OTP) GetRedisKey() string {
	return fmt.Sprintf("otp:%s:%s", o.VoterID, o.Purpose)
}

func (o *OTP) GetTimeRemaining() time.Duration {
	if o.IsExpired() {
		return 0
	}
	return time.Until(o.ExpiresAt)
}
