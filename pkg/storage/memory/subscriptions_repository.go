package memory

import (
	"sync"
)

// SubscriptionsRepository holds the CRUD db operations for CasinoRoundBet.
type SubscriptionsRepository struct {
	sync.RWMutex
	subsStore map[string]int
}

// NewSubscriptionsRepository is a constructor function for SubscriptionsRepository.
func NewSubscriptionsRepository() *SubscriptionsRepository {
	return &SubscriptionsRepository{
		subsStore: make(map[string]int, 0),
	}
}

// Insert inserts a new blocks.Transaction entity.
func (r *SubscriptionsRepository) Insert(address string) {
	r.Lock()
	r.subsStore[address] = -1
	r.Unlock()
}

// GetAllSubscriptions returns all address subscriptions.
func (r *SubscriptionsRepository) GetAllSubscriptions() []string {
	addresses := make([]string, len(r.subsStore))

	r.RLock()
	for k, i := range r.subsStore {
		addresses[i] = k
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
