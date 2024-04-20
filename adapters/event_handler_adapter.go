package adapters

import (
	"log"

	"go.uber.org/zap"

	"github.com/ruancaetano/gotcoin/core/events"
	"github.com/ruancaetano/gotcoin/core/protocols"
	"github.com/ruancaetano/gotcoin/infra"
)

type EventHandlerAdapter struct {
	blockchainService protocols.BlockchainService
	logger            *zap.Logger
}

func NewEventHandlerAdapter(blockchainService protocols.BlockchainService) protocols.EventHandler {
	return &EventHandlerAdapter{
		blockchainService: blockchainService,
		logger:            infra.GetLogger("EventHandlerAdapter"),
	}
}

func (eh *EventHandlerAdapter) Handle(event events.EventData) {
	eh.logger.Debug("Event received", zap.String("type", event.Type))
	switch event.Type {
	case events.RequestBlockChainSyncEventType:
		eh.handleBlockChainSyncRequest(event)
		break
	case events.ResponseBlockChainSyncEventType:
		eh.handleBlockChainSyncResponse(event)
		break
	case events.NewTransactionEventType:
		eh.handleNewTransaction(event)
	case events.NewBlockEventType:
		eh.handleNewBlock(event)
	default:
		log.Println("Unknown event type")
	}
}

func (eh *EventHandlerAdapter) handleBlockChainSyncRequest(event events.EventData) {
	if event.Metadata.OriginPeerID == "" {
		log.Println("OriginPeerID is required")
		return
	}
	eh.blockchainService.SendSyncBlocks(event.Metadata.OriginPeerID)
}

func (eh *EventHandlerAdapter) handleBlockChainSyncResponse(event events.EventData) {
	if event.Block == nil {
		log.Println("Block is required")
		return
	}

	eh.blockchainService.ReceiveSyncBlock(event.Block, *event.BlockCount, *event.NewDifficult)
}

func (eh *EventHandlerAdapter) handleNewTransaction(event events.EventData) {
	if event.Transaction == nil {
		log.Println("Transaction is required")
		return
	}

	if err := eh.blockchainService.AddTransaction(event.Transaction); err != nil {
		log.Println("Skipping transaction: ", err.Error())
	}
}

func (eh *EventHandlerAdapter) handleNewBlock(event events.EventData) {
	if event.Block == nil {
		log.Println("Block is required")
		return
	}

	eh.blockchainService.AddBlock(event.Block, *event.NewDifficult)
}
