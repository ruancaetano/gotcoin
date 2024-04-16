package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/network"
)

var bc *core.BlockChain
var node *network.Node
var wallet = core.Wallet1
var walletTwo = core.Wallet2

var cliOptions = []string{"Create transaction", "Get Balance"}

func main() {
	ctx := context.Background()
	host, err := network.InitHost(false, 0)
	if err != nil {
		panic(err)
	}

	bc = core.NewEmptyBlockChain()
	eh := core.NewEventHandler(bc)

	node = network.NewNode(ctx, host, bc, eh)
	for !bc.Synced {
	}

	fmt.Println("Blockchain synced")

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
		transaction := core.NewTransaction(wallet.PublicKey, walletTwo.PublicKey, 1)
		transaction.Sign(wallet.PrivateKey)
		node.SendBroadcastEvent(core.SendNewTransactionEvent(transaction))

	case 2:
		fmt.Println("Balance 1: ", bc.GetBalance(wallet.PublicKey))
		fmt.Println("Balance 2: ", bc.GetBalance(walletTwo.PublicKey))
	}
}
