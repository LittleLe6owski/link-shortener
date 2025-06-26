package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
	logger zerolog.Logger
	config *GRPCServerConfig
}

type GRPCServerConfig struct {
	Host            string
	Port            string
	ShutdownTimeout time.Duration
	Options         []grpc.ServerOption
}

func NewGRPCServer(config *GRPCServerConfig) *GRPCServer {
	server := grpc.NewServer(config.Options...)
	reflection.Register(server)
	return &GRPCServer{
		server: server,
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
		config: config,
	}
}

func (s *GRPCServer) Start() error {
	addr := s.config.Host + ":" + s.config.Port
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %q: %w", addr, err)
	}

	go func() {
		s.logger.Info().Msgf("grpc server started at %s:%s", s.config.Host, s.config.Port)
		_ = s.server.Serve(l)
	}()

	return nil
}

func (s *GRPCServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	s.logger.Info().Msg("gracefully shutting down grpc server...")

	serverStoppedCh := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(serverStoppedCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-serverStoppedCh:
		return nil
	}
}

func (s *GRPCServer) SetLogger(l zerolog.Logger) {
	s.logger = l
}

func (s *GRPCServer) RegisterService(serviceDesc *grpc.ServiceDesc, implementation interface{}) {
	s.server.RegisterService(serviceDesc, implementation)
}
