package consumer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/event/handler"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/utils/encryption"
	"github.com/nocturna-ta/vote/config"
	"github.com/nocturna-ta/vote/internal/interfaces/dao"
	"github.com/nocturna-ta/vote/internal/usecases"
	"github.com/nocturna-ta/vote/internal/usecases/consumer"
)

type container struct {
	Cfg          config.MainConfig
	ConsumerUc   usecases.Consumer
	EventHandler handler.EventHandler
}

type options struct {
	Cfg       *config.MainConfig
	Publisher event.MessagePublisher
	DB        *sql.Store
	Client    ethereum.Client
}

func newContainer(opts *options) *container {
	encryptor, err := encryption.NewEncryption(opts.Cfg.Encryption.Key)
	if err != nil {
		log.Fatal("failed to create encryption: %v", err)
	}

	voteRepo := dao.NewVoteRepository(&dao.OptsVoteRepository{
		Client:          opts.Client,
		ContractAddress: common.HexToAddress(opts.Cfg.Blockchain.ElectionManagerAddress),
		DB:              opts.DB,
		Encryptor:       encryptor,
	})

	consumerUc := consumer.New(&consumer.Options{
		VoteRepo:       voteRepo,
		Publisher:      opts.Publisher,
		Topics:         opts.Cfg.Kafka.Topics,
		MaxRetries:     opts.Cfg.Kafka.Consumer.MaxRetries,
		DLQTopic:       opts.Cfg.Kafka.Topics.VoteDLQ.Value,
		ProcessedTopic: opts.Cfg.Kafka.Topics.VoteProcessed.Value,
	})

	eventHandler := handler.New(&handler.Options{
		RetryConfig: handler.RetryConfig{
			MaxRetry:          opts.Cfg.Kafka.Consumer.Retry.MaxRetry,
			RetryInitialDelay: opts.Cfg.Kafka.Consumer.Retry.RetryInitialDelay,
			MaxJitter:         opts.Cfg.Kafka.Consumer.Retry.MaxJitter,
			HandlerTimeout:    opts.Cfg.Kafka.Consumer.Retry.HandlerTimeout,
			BackOffConfig:     opts.Cfg.Kafka.Consumer.Retry.BackOffConfig,
		},
		Publisher:   opts.Publisher,
		DlqTopic:    opts.Cfg.Kafka.Topics.VoteDLQ.Value,
		ServiceName: "vote-service",
	})

	return &container{
		Cfg:          *opts.Cfg,
		ConsumerUc:   consumerUc,
		EventHandler: eventHandler,
	}

}
