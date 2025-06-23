package otp

import (
	"github.com/nocturna-ta/golib/cache"
	"github.com/nocturna-ta/golib/cache/redlock"
	"github.com/nocturna-ta/vote/config"
	"github.com/nocturna-ta/vote/internal/infrastructures/sms"
	"github.com/nocturna-ta/vote/internal/usecases"
	"time"
)

type Module struct {
	redis       cache.Cache
	redLock     redlock.RedLock
	otpConfig   config.OTPConfig
	smsProvider sms.SMSProvider
	smsConfig   config.SMSConfig
}

type Options struct {
	Redis       cache.Cache
	RedLock     redlock.RedLock
	OtpConfig   config.OTPConfig
	SMSProvider sms.SMSProvider
	SMSConfig   config.SMSConfig
}

func New(opts *Options) usecases.OTPUseCases {
	if opts.OtpConfig.Length == 0 {
		opts.OtpConfig.Length = 6
	}
	if opts.OtpConfig.TTL == 0 {
		opts.OtpConfig.TTL = 5 * time.Minute
	}
	if opts.OtpConfig.MaxRetries == 0 {
		opts.OtpConfig.MaxRetries = 3
	}

	return &Module{
		redis:       opts.Redis,
		redLock:     opts.RedLock,
		otpConfig:   opts.OtpConfig,
		smsProvider: opts.SMSProvider,
		smsConfig:   opts.SMSConfig,
	}
}
