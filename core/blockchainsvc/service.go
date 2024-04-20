package blockchainsvc

import (
	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/core/protocols"
)

type blockchainServiceImpl struct {
	Blockchain *core.BlockChain
	Node       protocols.Node
}

func NewBlockchainServiceImpl(bc *core.BlockChain, node protocols.Node) protocols.BlockchainService {
	return &blockchainServiceImpl{
		Blockchain: bc,
		Node:       node,
	}
}
