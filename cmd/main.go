package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	server "github.com/PorgGopher/fiber-plus-htmx/internal"
	"github.com/PorgGopher/fiber-plus-htmx/internal/api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	app := api.NewRouteHandler(logger)
	server := server.NewServer(logger, app)

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGQUIT, syscall.SIGINT)

	go func() {
		s := <-sigCh
		logger.Info("caught shutdown signal", "code", s)
		cancel()
	}()

	server.Start(ctx)
}
