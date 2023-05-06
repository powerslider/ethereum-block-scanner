package memory

import (
	"sync"

	"github.com/powerslider/ethereum-block-scanner/pkg/blocks"
)

// TransactionsRepository holds the CRUD db operations for CasinoRoundBet.
type TransactionsRepository struct {
	inboundStore     *MultiMap[string, blocks.Transaction]
	outboundStore    *MultiMap[string, blocks.Transaction]
	latestBlockStore sync.Map
}

// NewTransactionsRepository is a constructor function for TransactionsRepository.
func NewTransactionsRepository() *TransactionsRepository {
	return &TransactionsRepository{
		inboundStore:     New[string, blocks.Transaction](),
		outboundStore:    New[string, blocks.Transaction](),
		latestBlockStore: sync.Map{},
	}
}

// Insert inserts a new blocks.Transaction entity.
func (r *TransactionsRepository) Insert(address string, blockNumber int, tx blocks.Transaction, isInbound bool) {
	r.latestBlockStore.Store(address, blockNumber)

	if isInbound {
		r.inboundStore.Put(address, tx)
	} else {
		r.outboundStore.Put(address, tx)
	}
}

// GetLatestBlockNumberPerAddress returns the latest block containing transactions to/from a given address.
func (r *TransactionsRepository) GetLatestBlockNumberPerAddress(address string) int {
	blockNum, found := r.latestBlockStore.Load(address)
	if !found {
		return -1
	}

	return blockNum.(int)
}

// GetInboundTransactionsPerAddress returns all inbound transactions per a given address.
func (r *TransactionsRepository) GetInboundTransactionsPerAddress(address string) []blocks.Transaction {
	txs, found := r.inboundStore.Get(address)
	if !found {
		return nil
	}

	return txs
}

// GetOutboundTransactionsPerAddress returns all outbound transactions per a given address.
func (r *TransactionsRepository) GetOutboundTransactionsPerAddress(address string) []blocks.Transaction {
	txs, found := r.outboundStore.Get(address)
	if !found {
		return nil
	}

	return txs
}

// GetAllTransactionsPerAddress returns all inbound transactions per a given address.
func (r *TransactionsRepository) GetAllTransactionsPerAddress(address string) []blocks.Transaction {
	txs := make([]blocks.Transaction, 0)

	txs = append(txs, r.GetInboundTransactionsPerAddress(address)...)
	txs = append(txs, r.GetOutboundTransactionsPerAddress(address)...)

	return txs
}
