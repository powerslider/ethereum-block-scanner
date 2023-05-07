package memory

import (
	"sync"

	"github.com/powerslider/ethereum-block-scanner/pkg/blocks"
)

// SubscriptionsRepository holds the CRUD db operations for CasinoRoundBet.
type SubscriptionsRepository struct {
	sync.RWMutex
	subsStore       map[string]int
	observedTxStore *MultiMap[string, blocks.Transaction]
}

// NewSubscriptionsRepository is a constructor function for SubscriptionsRepository.
func NewSubscriptionsRepository() *SubscriptionsRepository {
	return &SubscriptionsRepository{
		subsStore:       make(map[string]int, 0),
		observedTxStore: New[string, blocks.Transaction](),
	}
}

// InsertSubscriberAddress inserts a new address as a subscriber to be observed for new transactions.
func (r *SubscriptionsRepository) InsertSubscriberAddress(address string) {
	r.Lock()
	r.subsStore[address] = -1
	r.Unlock()
}

// InsertObservedTransaction inserts a new transaction that involves a subscribed address.
func (r *SubscriptionsRepository) InsertObservedTransaction(address string, tx blocks.Transaction) {
	r.observedTxStore.Put(address, tx)
}

// GetAllSubscriptions returns all address subscriptions.
func (r *SubscriptionsRepository) GetAllSubscriptions() []string {
	addresses := make([]string, len(r.subsStore))

	var i int

	r.RLock()
	for k := range r.subsStore {
		addresses[i] = k
		i++
	}
	r.RUnlock()

	return addresses
}

// GetLastCheckedBlockNumberPerAddress returns the last check block number a given subscribed address.
func (r *SubscriptionsRepository) GetLastCheckedBlockNumberPerAddress(address string) int {
	r.RLock()
	blockNum, found := r.subsStore[address]
	r.RUnlock()

	if !found {
		return -1
	}

	return blockNum
}

// GetObservedTransactionsPerAddress returns all observed transactions per subscribed address.
func (r *SubscriptionsRepository) GetObservedTransactionsPerAddress(address string) []blocks.Transaction {
	txs, found := r.observedTxStore.Get(address)
	if !found {
		return make([]blocks.Transaction, 0)
	}

	return txs
}
