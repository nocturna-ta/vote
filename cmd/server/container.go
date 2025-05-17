package server

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/txmanager"
	txSql "github.com/nocturna-ta/golib/txmanager/sql"
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
	Cfg    *config.MainConfig
	DB     *sql.Store
	Client *ethclient.Client
}

func newContainer(opts *options) *container {
	voteRepo, err := dao.NewBlockchainRepository(&dao.OptsVoteRepository{
		Client:          opts.Client,
		ContractAddress: common.HexToAddress(opts.Cfg.Blockchain.ContractAddress),
	})
	if err != nil {
		log.Fatal("Failed to initiate vote repository")
	}

	_, err = txmanager.New(context.Background(), &txmanager.DriverConfig{
		Type: "sql",
		Config: txSql.Config{
			DB: opts.DB,
		},
	})
	if err != nil {
		log.Fatal("Failed to instantiate transaction manager ")
	}

	voteUc := vote.New(&vote.Opts{
		VoteRepo: voteRepo,
	})
	return &container{
		Cfg:    *opts.Cfg,
		VoteUc: voteUc,
	}
}
