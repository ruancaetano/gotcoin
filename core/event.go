package core

import (
	"github.com/google/uuid"
)

const RequestBlockChainSyncEventType = "request_blockchain_sync"
const ResponseBlockChainSyncEventType = "response_blockchain_sync"
const NewBlockEventType = "new_block"
const NewTransactionEventType = "new_transaction"

type EventMetadata struct {
	OriginPeerID  string `json:"origin_peer_id"`
	FromPeerID    string `json:"from_peer_id"`
	MustPropagate bool   `json:"must_propagate"`
}

type EventData struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Block        *Block       `json:"block"`
	BlockCount   *int         `json:"block_count"`
	NewDifficult *int         `json:"new_difficult"`
	Transaction  *Transaction `json:"transaction"`
	Metadata     EventMetadata
}

func RequestBlockChainSyncEvent() EventData {
	return EventData{
		ID:   uuid.NewString(),
		Type: RequestBlockChainSyncEventType,
	}
}

func ResponseBlockChainSyncEvent(block Block, blockCount int) EventData {
	return EventData{
		ID:         uuid.NewString(),
		Type:       ResponseBlockChainSyncEventType,
		Block:      &block,
		BlockCount: &blockCount,
	}
}

func SendNewBlockEvent(block *Block, newDifficult *int) EventData {
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

func SendNewTransactionEvent(transaction *Transaction) EventData {
	return EventData{
		ID:          uuid.NewString(),
		Type:        NewTransactionEventType,
		Transaction: transaction,
		Metadata: EventMetadata{
			MustPropagate: true,
		},
	}
}
