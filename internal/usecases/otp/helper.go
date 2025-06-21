package otp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/nocturna-ta/vote/internal/domain/model"
	"time"
)

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
