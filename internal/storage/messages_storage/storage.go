package storage

import (
	"context"

	"github.com/itmo-lite-chat/messages_svc/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	collection *mongo.Collection
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		collection: db.Collection("messages"), // Как на схеме m_messages
	}
}

// Метод должен называться Save, чтобы соответствовать интерфейсу из service.go
func (s *Storage) Save(ctx context.Context, msg *domain.Message) error {
	_, err := s.collection.InsertOne(ctx, msg)
	return err
}
