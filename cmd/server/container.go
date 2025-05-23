package server

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/txmanager"
	txSql "github.com/nocturna-ta/golib/txmanager/sql"
	"github.com/nocturna-ta/golib/utils/encryption"
	"github.com/nocturna-ta/vote/config"
	"github.com/nocturna-ta/vote/internal/interfaces/dao"
	"github.com/nocturna-ta/vote/internal/usecases"
	"github.com/nocturna-ta/vote/internal/usecases/vote"
	"log"
)

type container struct {
	Cfg    config.MainConfig
	VoteUc usecases.VoteUseCases
}

type options struct {
	Cfg       *config.MainConfig
	DB        *sql.Store
	Client    ethereum.Client
	Publisher event.MessagePublisher
}

func newContainer(opts *options) *container {

	encryptor, err := encryption.NewEncryption(opts.Cfg.Encryption.Key)
	if err != nil {
		log.Fatal("Failed to instantiate encryption service", err)
	}

	voteRepo := dao.NewVoteRepository(&dao.OptsVoteRepository{
		Client:          opts.Client,
		ContractAddress: common.HexToAddress(opts.Cfg.Blockchain.ElectionManagerAddress),
		DB:              opts.DB,
		Encryptor:       encryptor,
	})

	txMgr, err := txmanager.New(context.Background(), &txmanager.DriverConfig{
		Type: "sql",
		Config: txSql.Config{
			DB: opts.DB,
		},
	})
	if err != nil {
		log.Fatal("Failed to instantiate transaction manager ")
	}

	voteUc := vote.New(&vote.Opts{
		VoteRepo:  voteRepo,
		TxMgr:     txMgr,
		Publisher: opts.Publisher,
		Topics:    opts.Cfg.Kafka.Topics,
		Encryptor: encryptor,
	})
	return &container{
		Cfg:    *opts.Cfg,
		VoteUc: voteUc,
	}
}
