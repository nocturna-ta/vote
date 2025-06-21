package vote

import (
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/txmanager"
	"github.com/nocturna-ta/golib/utils/encryption"
	"github.com/nocturna-ta/vote/config"
	"github.com/nocturna-ta/vote/internal/domain/repository"
	"github.com/nocturna-ta/vote/internal/usecases"
)

type Module struct {
	voteRepo  repository.VoteRepository
	txMgr     txmanager.TxManager
	publisher event.MessagePublisher
	topics    config.KafkaTopics
	encryptor *encryption.Encryption
	otpUc     usecases.OTPUseCases
}

type Opts struct {
	VoteRepo    repository.VoteRepository
	TxMgr       txmanager.TxManager
	Publisher   event.MessagePublisher
	Topics      config.KafkaTopics
	Encryptor   *encryption.Encryption
	OTPUseCases usecases.OTPUseCases
}

func New(opts *Opts) usecases.VoteUseCases {
	return &Module{
		voteRepo:  opts.VoteRepo,
		txMgr:     opts.TxMgr,
		publisher: opts.Publisher,
		topics:    opts.Topics,
		encryptor: opts.Encryptor,
		otpUc:     opts.OTPUseCases,
	}
}
