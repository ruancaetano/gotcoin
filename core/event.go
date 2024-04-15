package core

const RequestBlockChainSyncEventType = "request_blockchain_sync"
const ResponseBlockChainSyncEventType = "response_blockchain_sync"
const NewBlockEventType = "new_block"
const NewTransactionEventType = "new_transaction"

type EventData struct {
	Type         string       `json:"type"`
	Block        *Block       `json:"block"`
	BlockCount   *int         `json:"block_count"`
	NewDifficult *int         `json:"new_difficult"`
	Transaction  *Transaction `json:"transaction"`
}

func RequestBlockChainSyncEvent() EventData {
	return EventData{
		Type: RequestBlockChainSyncEventType,
	}
}

func ResponseBlockChainSyncEvent(block Block, blockCount int) EventData {
	return EventData{
		Type:       ResponseBlockChainSyncEventType,
		Block:      &block,
		BlockCount: &blockCount,
	}
}

func SendNewBlockEvent(block *Block, newDifficult *int) EventData {
	return EventData{
		Type:         NewBlockEventType,
		Block:        block,
		NewDifficult: newDifficult,
	}
}

func SendNewTransactionEvent(transaction *Transaction) EventData {
	return EventData{
		Type:        NewTransactionEventType,
		Transaction: transaction,
	}
}
