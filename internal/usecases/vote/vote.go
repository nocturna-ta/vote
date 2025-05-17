package vote

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nocturna-ta/golib/tracing"
	"math/big"
)

func (m *Module) CastVote(ctx context.Context, ktp string, candidateID uint64) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteUseCases.CastVote")
	defer span.End()

	privateKey, _ := crypto.HexToECDSA("")
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(5777))

	candidateIDNew := new(big.Int).SetUint64(candidateID)

	tx, err := m.contractAddr.Vote(auth, ktp, candidateIDNew)

}
