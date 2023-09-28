package api

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RouteHandler struct {
	fiber  *fiber.App
	logger *slog.Logger
}

func NewRouteHandler(logger *slog.Logger) *RouteHandler {
	app := fiber.New()

	app.Get("/", getHandler(logger))

	app.Use(notFoundHandler(logger))

	return &RouteHandler{
		logger: logger,
		fiber:  app,
	}
}

func (r *RouteHandler) Start(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)

	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = ":3000"
	}

	go func() {
		defer cancel()
		<-ctx.Done()
		err := r.Stop()
		if err != nil {
			r.logger.Error("error shutting down server", "error", err)
		}
	}()

	err := r.fiber.Listen(port)
	if err != nil {
		r.logger.Error("listener returned error", "error", err)
		return err
	}
	return nil
}

func (r *RouteHandler) Stop() error {
	return r.fiber.ShutdownWithTimeout(time.Second * 15)
}

func getHandler(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Debug("executing get handler", "headers", c.GetReqHeaders())
		return c.SendString("yup")
	}

}

func notFoundHandler(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Debug("executing nothing was found handler", "headers", c.GetReqHeaders())
		return c.Status(fiber.StatusNotFound).SendString("404...these are not the droids you are looking for.")
	}
}
