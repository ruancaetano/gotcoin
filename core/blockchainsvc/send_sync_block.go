package blockchainsvc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ruancaetano/gotcoin/core/events"
)

func (bs *blockchainServiceImpl) SendSyncBlocks(targetPeerID string) {
	stream := bs.Node.GetPeerStream(targetPeerID)
	if stream == nil {
		log.Println("Peer not found")
		return
	}

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	for _, block := range bs.Blockchain.Blocks {
		// Send block to peer
		event := events.ResponseBlockChainSyncEvent(*block, len(bs.Blockchain.Blocks))
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
