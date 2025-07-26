package http

import (
	"context"
	"time"

	"github.com/valyala/fasthttp"
)

// Server wraps fasthttp.Server with additional functionality
type Server struct {
	server *fasthttp.Server
	router *Router
}

// NewServer creates a new HTTP server
func NewServer(router *Router) *Server {
	server := &fasthttp.Server{
		Handler:      router.Handler(),
		Name:         "Banner-API/1.0",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{
		server: server,
		router: router,
	}
}

// Start starts the HTTP server
func (s *Server) Start(address string) error {
	return s.server.ListenAndServe(address)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.ShutdownWithContext(ctx)
}
