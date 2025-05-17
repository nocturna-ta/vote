package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/vote/config"
)

func GetEthereumClient(cfg *config.BlockchainConfig) (ethereum.Client, error) {
	return ethereum.New(&ethereum.Options{
		URL: cfg.GanacheURL,
	})
}
func GetElectionContractAddress(cfg *config.BlockchainConfig) common.Address {
	return common.HexToAddress(cfg.ElectionManagerAddress)
}
