package consumer

import (
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/vote/config"
	"github.com/nocturna-ta/vote/internal/domain/repository"
	"github.com/nocturna-ta/vote/internal/usecases"
)

type Module struct {
	voteRepo       repository.VoteRepository
	publisher      event.MessagePublisher
	topics         config.KafkaTopics
	maxRetries     int
	dlqTopic       string
	processedTopic string
}

type Options struct {
	VoteRepo       repository.VoteRepository
	Publisher      event.MessagePublisher
	Topics         config.KafkaTopics
	MaxRetries     int
	DLQTopic       string
	ProcessedTopic string
}

func New(opts *Options) usecases.Consumer {
	maxRetries := 3
	if opts.MaxRetries > 0 {
		maxRetries = opts.MaxRetries
	}

	return &Module{
		voteRepo:       opts.VoteRepo,
		publisher:      opts.Publisher,
		topics:         opts.Topics,
		maxRetries:     maxRetries,
		dlqTopic:       opts.DLQTopic,
		processedTopic: opts.ProcessedTopic,
	}
}
