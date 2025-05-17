package vote

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nocturna-ta/vote/internal/domain/repository"
	"github.com/nocturna-ta/vote/internal/usecases"
	"github.com/nocturna-ta/vote/pkg/binding"
)

type Module struct {
	voteRepo     repository.VoteRepository
	ethConn      *ethclient.Client
	contractAddr binding.Votechain
}

type Opts struct {
	VoteRepo     repository.VoteRepository
	EthConn      *ethclient.Client
	ContractAddr common.Address
}

func New(opts *Opts) usecases.VoteUseCases {
	contract, err := binding.NewVotechain(opts.ContractAddr, opts.EthConn)
	if err != nil {
		return nil
	}

	return &Module{
		voteRepo:     opts.VoteRepo,
		ethConn:      opts.EthConn,
		contractAddr: contract,
	}
}
