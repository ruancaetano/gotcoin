package core

func (bc *BlockChain) CreateGenesisBlock() *Block {
	block, _ := NewBlock(0, "", []*Transaction{})
	return block
}
