package blockchainsvc

import (
	"github.com/ruancaetano/gotcoin/core"
)

func (bs *blockchainServiceImpl) GetBlocks() []*core.Block {
	return bs.Blockchain.Blocks
}
