package core

const RequestBlockChainSyncEventType = "request_blockchain_sync"
const ResponseBlockChainSyncEventType = "response_blockchain_sync"

type EventData struct {
	Type       string `json:"type"`
	Block      *Block `json:"block"`
	BlockCount *int   `json:"block_count"`
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
