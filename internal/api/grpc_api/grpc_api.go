package grpc_api

import (
	"context"

	"github.com/google/uuid"
	"github.com/itmo-lite-chat/messages_svc/internal/domain"
	pb "github.com/itmo-lite-chat/proto-registry/gen/services/messages_service/messages/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API struct {
	// Встраиваем обязательную структуру для gRPC серверов
	pb.UnimplementedMessagesServiceServer
	service Service
}

func NewAPI(service Service) *API {
	return &API{service: service}
}

// SendMessage — реализация метода из proto контракта
func (a *API) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid chatID")
	}
	// 1. Превращаем gRPC запрос в доменную модель (как на твоей схеме m_messages)
	msg := &domain.Message{
		ChatID: chatID,
		Content: domain.Content{
			Type: externalContentTypeToInternal(req.Content.Type),
			Body: req.Content.Body,
		},
		ReplyToID: req.ReplyToId,
	}

	// 2. Вызываем сервис
	if err := a.service.SendMessage(ctx, msg); err != nil {
		return nil, err
	}

	// 3. Возвращаем ответ согласно контракту
	return &pb.SendMessageResponse{}, nil
}
