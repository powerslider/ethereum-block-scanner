package sdk

import (
	"context"

	"github.com/powerslider/ethereum-block-scanner/pkg/transport/client/jsonrpc"

	"github.com/powerslider/ethereum-block-scanner/pkg/blocks"
)

// Parser is a port interface that defines SDK operations on the Ethereum blockchain.
type Parser interface {
	// GetCurrentBlock last parsed block.
	GetCurrentBlock(ctx context.Context) (int, error)

	// Subscribe adds an address to be observed for new transactions.
	Subscribe(address string) bool

	// GetBlockTransactions returns all transactions contained is a block.
	GetBlockTransactions(ctx context.Context, blockNum int) ([]blocks.Transaction, error)

	// GetTransactionsPerSubscriber lists observed transactions given a registered subscriber address.
	GetTransactionsPerSubscriber(address string) []blocks.Transaction

	// GetTransactionsForBlockRange lists inbound or outbound transactions for an address for a given block range
	// from latest to a specified one.
	GetTransactionsForBlockRange(ctx context.Context, address string, blockRange int) ([]blocks.Transaction, error)
}

// SubscriptionsStore is a port interface for storage operations related to address subscriptions.
type SubscriptionsStore interface {
	// GetAllSubscriptions returns all address subscriptions.
	GetAllSubscriptions() []string

	// InsertObservedTransaction inserts a new transaction that involves a subscribed address.
	InsertObservedTransaction(address string, tx blocks.Transaction)

	// InsertSubscriberAddress inserts a new address as a subscriber to be observed for new transactions.
	InsertSubscriberAddress(address string)

	// GetObservedTransactionsPerAddress returns all observed transactions per subscribed address.
	GetObservedTransactionsPerAddress(address string) []blocks.Transaction
}

// TransactionHistoryStore is a port interface for storage operations on transaction history for a given address.
type TransactionHistoryStore interface {
	// Insert inserts a new blocks.Transaction entity.
	Insert(address string, blockNumber int, tx blocks.Transaction, isInbound bool)

	// GetLatestBlockNumberPerAddress returns the latest block containing transactions to/from a given address.
	GetLatestBlockNumberPerAddress(address string) int

	// GetAllTransactionsPerAddress returns all inbound transactions per a given address.
	GetAllTransactionsPerAddress(address string) []blocks.Transaction
}

// RPCClient is a port interface defining JSON-RPC methods.
type RPCClient interface {
	// Call calls a JSON-RPC method with optional params.
	Call(ctx context.Context, method string, params ...any) (*jsonrpc.RPCResponse, error)

	// CallFor calls a JSON-RPC method and deserializes the response in a specified response object.
	CallFor(ctx context.Context, out any, method string, params ...any) error
}
