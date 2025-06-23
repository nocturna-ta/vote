package otp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/nocturna-ta/vote/internal/domain/model"
	"strings"
	"time"
)

func (m *Module) sendOTPSMS(ctx context.Context, phoneNumber, otpCode string) error {
	message := m.formatOTPMessage(otpCode)

	return m.smsProvider.SendSMS(ctx, phoneNumber, message)
}

func (m *Module) formatOTPMessage(otpCode string) string {
	template := m.smsConfig.Templates.OTPMessage
	if template == "" {
		template = "Kode OTP untuk voting Anda adalah: {{.Code}}. Berlaku selama {{.Minutes}} menit. Jangan bagikan kode ini kepada siapa pun."
	}

	message := strings.ReplaceAll(template, "{{.Code}}", otpCode)
	message = strings.ReplaceAll(message, "{{.Minutes}}", fmt.Sprintf("%.0f", m.otpConfig.TTL.Minutes()))

	return message
}

func (m *Module) getOTPGenerationMessage(phoneNumber string) string {
	if m.smsConfig.Enabled && phoneNumber != "" {
		return fmt.Sprintf("OTP generated and sent to %s", m.maskPhoneNumber(phoneNumber))
	}
	return "OTP generated successfully"
}

func (m *Module) maskPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) < 4 {
		return "****"
	}

	if len(phoneNumber) > 8 {
		return phoneNumber[:3] + "****" + phoneNumber[len(phoneNumber)-4:]
	}

	return phoneNumber[:2] + "****"
}

func (m *Module) storeOTPInRedis(ctx context.Context, otp *model.OTP) error {
	key := otp.GetRedisKey()
	ttl := time.Until(otp.ExpiresAt)
	return m.redis.Set(ctx, key, otp, int(ttl))
}

func (m *Module) getOTPFromRedis(ctx context.Context, voterID, purpose string) (*model.OTP, error) {
	key := fmt.Sprintf("otp:%s:%s", voterID, purpose)
	var otp model.OTP
	err := m.redis.GetObject(ctx, key, &otp)
	return &otp, err
}

func (m *Module) deleteOTPFromRedis(ctx context.Context, voterID, purpose string) error {
	key := fmt.Sprintf("otp:%s:%s", voterID, purpose)
	return m.redis.Delete(ctx, key)
}

func (m *Module) generateOTPToken(voterID, purpose string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
