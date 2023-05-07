package handlers

import (
	"github.com/gorilla/mux"
	"github.com/powerslider/ethereum-block-scanner/pkg/configs"
	"github.com/powerslider/ethereum-block-scanner/pkg/sdk"
)

// InitializeHandlers registers HTTP routes and wires dependencies for HTTP handlers.
func InitializeHandlers(
	config *configs.Config,
	router *mux.Router,
	parser sdk.Parser,
) *mux.Router {
	handler := NewBlockHandler(parser)

	registerHTTPRoutes(config, router, handler)

	return router
}
