package sms

import "context"

type SMSProvider interface {
	SendSMS(ctx context.Context, to, message string) error
	ValidatePhoneNumber(phoneNumber string) error
	GetProviderName() string
}

type SMSMessage struct {
	To      string
	Message string
	From    string
}

type SMSResponse struct {
	Success   bool
	MessageID string
	Error     string
	Cost      float64
	Provider  string
}
