package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/powerslider/ethereum-block-scanner/pkg/configs"
)

// Server represents an HTTP server.
type Server struct {
	server *http.Server
}

// NewServer constructs new HTTP server with the provided muxer.
func NewServer(
	config *configs.Config,
	muxer *mux.Router,
) *Server {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: muxer,
	}

	return &Server{
		server: server,
	}
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context, errChan chan error) {
	log.Printf("[Start] HTTP server is starting on %s:\n", s.server.Addr)

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		errChan <- errors.WithStack(err)
	}
}

// Stop stops the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	log.Println("[Shutdown] HTTP server is shutting down...")

	shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.server.Shutdown(shutdownCtx)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
