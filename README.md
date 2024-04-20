# GotCoin

## Description
GotCoin is an experiment used to learn about blockchain technology. The project is written in Golang, and contains elements like data transfer on a p2p network.

### Features
- [x] Proof of Work
- [x] Transactions
- [x] Wallets
- [x] Mining
- [x] P2P Network
- [x] CLI to send transactions
- [ ] Consensus Algorithm

## How to enjoy

Run the first node with the following command:
It`s a genesis node, so it will create the first block of the blockchain.
```
    go ./cmd -genesis=true
```

Run the other nodes with the following command. It will connect to the first node and receive the blockchain data.
```
    go ./cmd
```

To send a transaction, use the simple cli and follow the instructions.
```
    go ./cli
```

