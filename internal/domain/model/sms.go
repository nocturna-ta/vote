package model

import (
	"fmt"
	"time"
)

type SMSLog struct {
	ID          string    `json:"id" db:"id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Message     string    `json:"message" db:"message"`
	Provider    string    `json:"provider" db:"provider"`
	Status      string    `json:"status" db:"status"` // sent, failed, pending
	MessageID   string    `json:"message_id" db:"message_id"`
	Cost        float64   `json:"cost" db:"cost"`
	Error       string    `json:"error" db:"error"`
	SentAt      time.Time `json:"sent_at" db:"sent_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type PhoneNumber struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
	Full        string `json:"full"`
}

func ParsePhoneNumber(phoneNumber string) (*PhoneNumber, error) {
	if len(phoneNumber) < 3 || phoneNumber[0] != '+' {
		return nil, fmt.Errorf("phone number must start with + and country code")
	}

	return &PhoneNumber{
		Full: phoneNumber,
	}, nil
}
