package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
)

type EventHandler struct {
	BlockChain *BlockChain
}

func NewEventHandler(bc *BlockChain) *EventHandler {
	return &EventHandler{bc}
}

func (eh *EventHandler) HandleEvent(rw *bufio.ReadWriter, event EventData) {
	switch event.Type {
	case RequestBlockChainSyncEventType:
		eh.handleBlockChainSyncRequest(rw)
		break
	case ResponseBlockChainSyncEventType:
		eh.handleBlockChainSyncResponse(event)
	}
}

func (eh *EventHandler) handleBlockChainSyncRequest(rw *bufio.ReadWriter) {
	for _, block := range eh.BlockChain.Blocks {
		// Send block to peer
		event := ResponseBlockChainSyncEvent(*block, len(eh.BlockChain.Blocks))
		eventBytes, err := json.Marshal(event)
		if err != nil {
			log.Println(err)
		}

		if _, err = rw.WriteString(fmt.Sprintf("%s\n", eventBytes)); err != nil {
			log.Println(err)
		}

		if err = rw.Flush(); err != nil {
			log.Println(err)
		}
	}
}

func (eh *EventHandler) handleBlockChainSyncResponse(event EventData) {
	if event.Block == nil {
		return
	}

	if event.Block.IsValid() {
		eh.BlockChain.AddBlock(event.Block)
	}

	if event.BlockCount != nil && len(eh.BlockChain.Blocks) < *event.BlockCount {
		log.Println("Sync percent complete: ", float64(len(eh.BlockChain.Blocks))/float64(*event.BlockCount)*100.0)
		return
	}

	if event.BlockCount != nil && len(eh.BlockChain.Blocks) == *event.BlockCount {
		eh.BlockChain.SortBlocks()
		eh.BlockChain.Synced = true
		log.Println("Block chain sync completed")
		log.Println("Blocks count: ", len(eh.BlockChain.Blocks))
		log.Println("Block chain valid: ", eh.BlockChain.IsChainValid())
		return
	}
}
