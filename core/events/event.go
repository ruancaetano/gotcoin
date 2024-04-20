package events

import (
	"github.com/ruancaetano/gotcoin/core"
)

type EventMetadata struct {
	OriginPeerID  string `json:"origin_peer_id"`
	FromPeerID    string `json:"from_peer_id"`
	MustPropagate bool   `json:"must_propagate"`
}

type EventData struct {
	ID           string            `json:"id"`
	Type         string            `json:"type"`
	Block        *core.Block       `json:"block"`
	BlockCount   *int              `json:"block_count"`
	NewDifficult *int              `json:"new_difficult"`
	Transaction  *core.Transaction `json:"transaction"`
	Metadata     EventMetadata
}
