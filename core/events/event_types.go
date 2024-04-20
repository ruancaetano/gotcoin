package events

import (
	"github.com/google/uuid"

	"github.com/ruancaetano/gotcoin/core"
)

const RequestBlockChainSyncEventType = "request_blockchain_sync"
const ResponseBlockChainSyncEventType = "response_blockchain_sync"
const NewBlockEventType = "new_block"
const NewTransactionEventType = "new_transaction"

func RequestBlockChainSyncEvent() EventData {
	return EventData{
		ID:   uuid.NewString(),
		Type: RequestBlockChainSyncEventType,
	}
}

func ResponseBlockChainSyncEvent(block core.Block, blockCount int) EventData {
	return EventData{
		ID:         uuid.NewString(),
		Type:       ResponseBlockChainSyncEventType,
		Block:      &block,
		BlockCount: &blockCount,
	}
}

func SendNewBlockEvent(block *core.Block, newDifficult *int) EventData {
	return EventData{
		ID:           uuid.NewString(),
		Type:         NewBlockEventType,
		Block:        block,
		NewDifficult: newDifficult,
		Metadata: EventMetadata{
			MustPropagate: true,
		},
	}
}

func SendNewTransactionEvent(transaction *core.Transaction) EventData {
	return EventData{
		ID:          uuid.NewString(),
		Type:        NewTransactionEventType,
		Transaction: transaction,
		Metadata: EventMetadata{
			MustPropagate: true,
		},
	}
}
