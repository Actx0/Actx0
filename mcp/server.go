// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package mcp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
)

// Options configures the MCP proxy server.
type Options struct {
	Name    string
	Version string
	Mgt     string
	Port    int
	Path    string
	Timeout time.Duration
}

// Run starts the actx0 MCP server over streamable HTTP.
func Run(ctx context.Context, opts Options) error {
	server := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    opts.Name,
		Version: opts.Version,
	}, nil)

	api := NewAPIClient(opts.Mgt, opts.Timeout)

	server.AddReceivingMiddleware(ScopeMiddleware)
	RegisterMemoryTools(server, api)

	handler := mcpsdk.NewStreamableHTTPHandler(func(*http.Request) *mcpsdk.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle(opts.Path, HTTPScope(handler))

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", opts.Port),
		Handler:           mux,
		MaxHeaderBytes:    230 * 1024,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadTimeout:       opts.Timeout + 5*time.Second,
		WriteTimeout:      opts.Timeout + 5*time.Second,
	}

	serr := make(chan error, 1)

	go func() {
		log.Info().
			Str("mgt", opts.Mgt).
			Int("port", opts.Port).
			Str("path", opts.Path).
			Str("name", opts.Name).
			Str("version", opts.Version).
			Msg("Starting MCP HTTP server")

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serr:
		return fmt.Errorf("mcp server error: %w", err)
	case sig := <-quit:
		log.Info().
			Str("signal", sig.String()).
			Msg("Received MCP shutdown signal")

		shutdownTimeout := 30 * time.Second
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		log.Info().
			Dur("timeout", shutdownTimeout).
			Msg("Gracefully shutting down MCP server")

		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			return fmt.Errorf("mcp server forced to shutdown: %w", err)
		}

		log.Info().Msg("MCP server shutdown complete")
		return nil
	}
}
