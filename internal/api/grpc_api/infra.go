package grpc_api

import (
	"context"

	"github.com/google/uuid"
	"github.com/itmo-lite-chat/messages_svc/internal/domain"
)

type Service interface {
	SendMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error)
	GetChatHistory(ctx context.Context, chatID uuid.UUID, limit int32, beforeMessageID *int64) ([]*domain.Message, error)
	GetLastMessages(ctx context.Context, chatIDs []uuid.UUID) (map[uuid.UUID]*domain.Message, error)
	EditMessage(ctx context.Context, msgID int64, senderID uuid.UUID, newContent domain.Content) (*domain.Message, error)
	DeleteMessage(ctx context.Context, msgID int64, userID string) error
	GetUnreadCount(ctx context.Context, chatID uuid.UUID, lastReadID int64) (int32, error)
}
