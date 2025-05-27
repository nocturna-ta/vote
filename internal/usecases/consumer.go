package usecases

import (
	"context"
	"github.com/nocturna-ta/golib/event"
)

type Consumer interface {
	ConsumeVoteSubmit(ctx context.Context, message *event.EventConsumeMessage) error
}
