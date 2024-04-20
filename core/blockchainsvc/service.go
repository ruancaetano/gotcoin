package blockchainsvc

import (
	"go.uber.org/zap"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/core/protocols"
	"github.com/ruancaetano/gotcoin/infra"
)

type blockchainServiceImpl struct {
	Blockchain *core.BlockChain
	Node       protocols.Node
	Logger     *zap.Logger
}

func NewBlockchainServiceImpl(bc *core.BlockChain, node protocols.Node) protocols.BlockchainService {
	return &blockchainServiceImpl{
		Blockchain: bc,
		Node:       node,
		Logger:     infra.GetLogger("BlockchainService"),
	}
}
