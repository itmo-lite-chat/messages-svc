package clients

import (
	"context"

	"github.com/google/uuid"
	"github.com/itmo-lite-chat/messages_svc/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type MessageStorage interface {
	Save(ctx context.Context, msg *domain.Message) error
	GetChatHistory(ctx context.Context, chatID uuid.UUID, limit int64, beforeMessageID *int64) ([]*domain.Message, error)
	GetLastMessages(ctx context.Context, chatIDs []uuid.UUID) (map[uuid.UUID]*domain.Message, error)
	Update(ctx context.Context, msgID int64, update bson.M) error
	GetByID(ctx context.Context, msgID int64) (*domain.Message, error)
	Count(ctx context.Context, filter bson.M) (int32, error)
}
