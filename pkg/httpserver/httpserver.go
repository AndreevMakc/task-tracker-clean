package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	server *http.Server
}

func New(router *chi.Mux, opts ...Option) *Server {
	s := &Server{
		server: &http.Server{
			Handler:      router,
			ReadTimeout:  _defaultReadTimeout,
			WriteTimeout: _defaultWriteTimeout,
			IdleTimeout:  _defaultIdleTimeout,
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

const (
	_defaultReadTimeout  = 5 * time.Second
	_defaultWriteTimeout = 10 * time.Second
	_defaultIdleTimeout  = 60 * time.Second
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = ":" + port
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func IdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.IdleTimeout = timeout
	}
}

func (s *Server) Start() error {
	if s.server.Addr == "" {
		return fmt.Errorf("server address not set")
	}
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
