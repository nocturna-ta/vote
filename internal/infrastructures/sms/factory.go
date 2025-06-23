package sms

import (
	"fmt"
	"github.com/nocturna-ta/vote/config"
)

func NewSMSProvider(cfg config.SMSConfig) (SMSProvider, error) {

	if !cfg.Enabled {
		return NewMockProvider(false), nil
	}
	switch cfg.Provider {
	case "twilio":
		return NewTwilioProvider(TwilioConfig{
			AccountSID: cfg.Twilio.AccountSID,
			AuthToken:  cfg.Twilio.AuthToken,
			FromNumber: cfg.Twilio.FromNumber,
		}), nil

	case "mock":
		return NewMockProvider(false), nil
	default:
		return nil, fmt.Errorf("unsupported SMS provider: %s", cfg.Provider)
	}
}
