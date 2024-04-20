package protocols

import (
	"context"

	"github.com/libp2p/go-libp2p/core/network"

	"github.com/ruancaetano/gotcoin/core/events"
)

type Node interface {
	GetID() string
	GetAddr() string
	Setup(ctx context.Context, eh EventHandler)
	HandleNewStream(s network.Stream)
	InitSync() error
	ReadEvent(s network.Stream)
	SendEvent(s network.Stream, event events.EventData) error
	SendBroadcastEvent(event events.EventData)
	PropagateEvent(eventData events.EventData)
	GetPeerStream(string) network.Stream
}
