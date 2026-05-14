package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/itmo-lite-chat/messages_svc/internal/api/grpc_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/itmo-lite-chat/proto-registry/gen/services/messages_service/messages/v1"
)

type server struct {
	addr   string
	server *grpc.Server
}

func NewServer(addr string, messagesAPI *grpc_api.API) *server {
	s := grpc.NewServer()

	pb.RegisterMessagesServiceServer(s, messagesAPI)

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
