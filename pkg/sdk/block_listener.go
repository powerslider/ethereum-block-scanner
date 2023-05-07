package sdk

import (
	"context"
	"strings"
	"time"

	"github.com/powerslider/ethereum-block-scanner/pkg/storage/memory"
)

// BlockListener implements SDK operations on the Ethereum blockchain.
type BlockListener struct {
	BlockParser *BlockParser
	TxStore     *memory.TransactionsRepository
	SubsStore   *memory.SubscriptionsRepository
}

// NewBlockListener is a constructor function for BlockParser.
func NewBlockListener(
	blockParser *BlockParser,
	txStore *memory.TransactionsRepository,
	subsStore *memory.SubscriptionsRepository,
) *BlockListener {
	return &BlockListener{
		BlockParser: blockParser,
		TxStore:     txStore,
		SubsStore:   subsStore,
	}
}

// ListenForNewTransactions implements polling for the latest block and matching if the subscribed addresses
// have inbound or outbound transactions.
func (p *BlockListener) ListenForNewTransactions(ctx context.Context, errCh chan error) {
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
				if a == strings.ToLower(tx.To) {
					p.TxStore.Insert(a, latestBlockNum, tx, true)
				} else if a == strings.ToLower(tx.From) {
					p.TxStore.Insert(a, latestBlockNum, tx, false)
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}
