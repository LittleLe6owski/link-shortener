package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type HTTPServer struct {
	server *http.Server
	router *http.ServeMux
	logger zerolog.Logger
	config *HTTPServerConfig
}

type HTTPServerConfig struct {
	Host            string
	Port            string
	ShutdownTimeout time.Duration
}

func NewHTTPServer(config *HTTPServerConfig) *HTTPServer {
	router := http.NewServeMux()

	return &HTTPServer{
		server: &http.Server{
			Handler: router,
		},
		router: router,
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
		config: config,
	}
}

func (s *HTTPServer) Start() error {
	addr := s.config.Host + ":" + s.config.Port
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %q: %w", addr, err)
	}

	go func() {
		s.logger.Info().Msgf("http server started at %s:%s", s.config.Host, s.config.Port)
		_ = s.server.Serve(l)
	}()

	return nil
}

func (s *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	s.logger.Info().Msg("gracefully shutting down http server...")
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) SetLogger(l zerolog.Logger) {
	s.logger = l
}

func (s *HTTPServer) Router() *http.ServeMux {
	return s.router
}
