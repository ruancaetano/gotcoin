package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ruancaetano/gotcoin/adapters"
	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/core/blockchainsvc"
	"github.com/ruancaetano/gotcoin/core/events"
	"github.com/ruancaetano/gotcoin/core/protocols"
	"github.com/ruancaetano/gotcoin/infra/network"
	"github.com/ruancaetano/gotcoin/util"
)

var bc *core.BlockChain
var bs protocols.BlockchainService
var node protocols.Node

var cliOptions = []string{"Create transaction", "Get Balance"}

func main() {
	ctx := context.Background()
	host, err := network.InitHost(false, 0)
	if err != nil {
		panic(err)
	}

	node = adapters.NewNodeAdapter(ctx, host)

	bc = core.NewEmptyBlockChain()
	bs = blockchainsvc.NewBlockchainServiceImpl(bc, node)
	eh := adapters.NewEventHandlerAdapter(bs)
	node.Setup(ctx, eh)

	time.Sleep(1 * time.Second)
	runCli()
}

func runCli() {
	for {
		// Print options
		for i, option := range cliOptions {
			fmt.Printf("%d. %s\n", i+1, option)
		}

		// Prompt user for input
		fmt.Print("Enter the number of your choice: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		// Convert input to integer
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 1 || choice > len(cliOptions) {
			fmt.Println("Invalid choice.")
		}
		handleSelectedOption(choice)
	}
}

func handleSelectedOption(choice int) {
	fmt.Printf("You selected: %s\n", cliOptions[choice-1])
	switch choice {
	case 1:
		transaction := core.NewTransaction(util.Wallet1.PublicKey, util.Wallet2.PublicKey, 1)
		transaction.Sign(util.Wallet1.PrivateKey)
		bs.AddTransaction(transaction)
		node.SendBroadcastEvent(events.SendNewTransactionEvent(transaction))

	case 2:
		fmt.Println("Balance 1: ", bs.GetBalance(util.Wallet1.PublicKey))
		fmt.Println("Balance 2: ", bs.GetBalance(util.Wallet2.PublicKey))
	}
}
