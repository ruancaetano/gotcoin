package util

import (
	"github.com/ruancaetano/gotcoin/core"
)

func CompareTransactionFunc(a, b *core.Transaction) bool {
	return a.Hash == b.Hash
}
