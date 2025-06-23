package sms

import (
	"context"
	"fmt"
	"github.com/nocturna-ta/golib/log"
	"regexp"
)

type MockProvider struct {
	shouldFail bool
}

func NewMockProvider(shouldFail bool) SMSProvider {
	return &MockProvider{
		shouldFail: shouldFail,
	}
}

func (m *MockProvider) SendSMS(ctx context.Context, to, message string) error {
	if m.shouldFail {
		return fmt.Errorf("mock SMS provider: simulated failure")
	}

	// Log the SMS instead of sending it (for development)
	log.WithFields(log.Fields{
		"to":       to,
		"message":  message,
		"provider": "mock",
	}).InfoWithCtx(ctx, "[MockSMS] SMS would be sent")

	return nil
}

func (m *MockProvider) ValidatePhoneNumber(phoneNumber string) error {
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phoneNumber) {
		return fmt.Errorf("invalid phone number format: %s", phoneNumber)
	}
	return nil
}

func (m *MockProvider) GetProviderName() string {
	return "mock"
}
