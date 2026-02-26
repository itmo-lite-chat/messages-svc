package grpcserver

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	addr   string
	server *grpc.Server
}

func NewServer(addr string) *server {
	s := grpc.NewServer()

	reflection.Register(s)

	return &server{
		addr:   addr,
		server: s,
	}
}

func (s *server) Run(_ context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc listen error: %w", err)
	}

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("grpc server shutdown: %w", err)
	}
	return nil
}

func (s *server) Shutdown(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}
