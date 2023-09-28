package server

import (
	"context"
	"log/slog"
	"os"

	"github.com/PorgGopher/fiber-plus-htmx/internal/api"
)

type Server struct {
	logger *slog.Logger
	server *api.RouteHandler
}

func NewServer(logger *slog.Logger, server *api.RouteHandler) *Server {
	return &Server{
		logger: logger,
		server: server,
	}
}

func (s *Server) Start(parentCtx context.Context) {

	ctx, cancel := context.WithCancel(parentCtx)

	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = ":3000"
	}

	go func() {
		defer cancel()
		<-ctx.Done()
		err := s.Stop()
		if err != nil {
			s.logger.Error("error shutting down server", "error", err)
		}
	}()

	err := s.server.Start(ctx)
	if err != nil {
		s.logger.Error("listener returned error", "error", err)
	}

}

func (s *Server) Stop() error {
	return s.server.Stop()
}
