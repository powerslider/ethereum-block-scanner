package handlers

import (
	"github.com/gorilla/mux"
	"github.com/powerslider/ethereum-block-scanner/pkg/configs"
	"github.com/powerslider/ethereum-block-scanner/pkg/sdk"
	"github.com/powerslider/ethereum-block-scanner/pkg/storage/memory"
	"github.com/powerslider/ethereum-block-scanner/pkg/transport/client/jsonrpc"
)

// InitializeHandlers registers HTTP routes and wires dependencies for HTTP handlers.
func InitializeHandlers(
	config *configs.Config, router *mux.Router, client *jsonrpc.RPCClient) *mux.Router {
	txStore := memory.NewTransactionsRepository()
	parser := sdk.NewBlockParser(client, txStore)
	handler := NewBlockHandler(parser)

	registerHTTPRoutes(config, router, handler)

	return router
}
