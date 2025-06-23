package sms

import (
	"context"
	"fmt"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"regexp"
)

type TwilioProvider struct {
	client     *twilio.RestClient
	fromNumber string
}

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

func NewTwilioProvider(config TwilioConfig) SMSProvider {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	return &TwilioProvider{
		client:     client,
		fromNumber: config.FromNumber,
	}
}

func (t *TwilioProvider) SendSMS(ctx context.Context, to, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetBody(message)

	params.SetFrom(t.fromNumber)
	params.SetTo(to)

	resp, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	if resp != nil && *resp.Status == "failed" {
		return err
	}
	return nil
}

func (t *TwilioProvider) ValidatePhoneNumber(phoneNumber string) error {
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phoneNumber) {
		return fmt.Errorf("invalid phone number format: %s (must be in international format, e.g., +628123456789)", phoneNumber)
	}
	return nil
}

func (t *TwilioProvider) GetProviderName() string {
	return "twilio"
}
