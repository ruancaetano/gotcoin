package protocols

import (
	"github.com/ruancaetano/gotcoin/core/events"
)

type EventHandler interface {
	Handle(event events.EventData)
}
