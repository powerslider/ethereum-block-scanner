package sdk

import (
	"context"
	"strings"
	"time"
)

// BlockObserver implements SDK operations on the Ethereum blockchain.
type BlockObserver struct {
	BlockParser Parser
	SubsStore   SubscriptionsStore
}

// NewBlockObserver is a constructor function for BlockObserver.
func NewBlockObserver(
	blockParser Parser,
	subsStore SubscriptionsStore,
) *BlockObserver {
	return &BlockObserver{
		BlockParser: blockParser,
		SubsStore:   subsStore,
	}
}

// ListenForNewTransactions implements polling for the latest block and matching if the subscribed addresses
// have inbound or outbound transactions contained in it.
func (p *BlockObserver) ListenForNewTransactions(ctx context.Context, errCh chan error) {
	for {
		addresses := p.SubsStore.GetAllSubscriptions()
		if len(addresses) == 0 {
			time.Sleep(5 * time.Second)

			continue
		}

		latestBlockNum, err := p.BlockParser.GetCurrentBlock(ctx)
		if err != nil {
			errCh <- err
		}

		blockTransactions, err := p.BlockParser.GetBlockTransactions(ctx, latestBlockNum)
		if err != nil {
			errCh <- err
		}

		for _, a := range addresses {
			for _, tx := range blockTransactions {
				if a == strings.ToLower(tx.To) || a == strings.ToLower(tx.From) {
					p.SubsStore.InsertObservedTransaction(a, tx)
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}
