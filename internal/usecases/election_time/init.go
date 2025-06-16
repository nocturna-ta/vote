package election_time

import (
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/txmanager"
	"github.com/nocturna-ta/vote/config"
	"github.com/nocturna-ta/vote/internal/domain/repository"
	"github.com/nocturna-ta/vote/internal/usecases"
)

type Module struct {
	electionTimeRepo repository.ElectionTimeRepository
	txMgr            txmanager.TxManager
	publisher        event.Publisher
	topics           config.KafkaTopics
}

type Opts struct {
	ElectionTimeRepo repository.ElectionTimeRepository
	TxMgr            txmanager.TxManager
	Publisher        event.Publisher
	Topics           config.KafkaTopics
}

func New(opts *Opts) usecases.ElectionTimeUseCases {
	return &Module{
		electionTimeRepo: opts.ElectionTimeRepo,
		txMgr:            opts.TxMgr,
		publisher:        opts.Publisher,
		topics:           opts.Topics,
	}
}
