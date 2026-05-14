package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/itmo-lite-chat/messages_svc/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	collection *mongo.Collection
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		collection: db.Collection("messages"),
	}
}

func (s *Storage) Save(ctx context.Context, msg *domain.Message) error {
	msg.CreatedAt = time.Now()
	msg.Timestamp = msg.CreatedAt

	_, err := s.collection.InsertOne(ctx, msg)
	return err
}

func (s *Storage) GetChatHistory(ctx context.Context, chatID uuid.UUID, limit int64, beforeMessageID *int64) ([]*domain.Message, error) {
	filter := bson.M{
		"chat_id":    chatID,
		"deleted_at": nil,
	}
	if beforeMessageID != nil {
		filter["_id"] = bson.M{"$lt": *beforeMessageID}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(limit)

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*domain.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	for left, right := 0, len(messages)-1; left < right; left, right = left+1, right-1 {
		messages[left], messages[right] = messages[right], messages[left]
	}

	return messages, nil
}

func (s *Storage) GetLastMessages(ctx context.Context, chatIDs []uuid.UUID) (map[uuid.UUID]*domain.Message, error) {
	result := make(map[uuid.UUID]*domain.Message, len(chatIDs))
	for _, chatID := range chatIDs {
		filter := bson.M{
			"chat_id":    chatID,
			"deleted_at": nil,
		}
		opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})

		var msg domain.Message
		err := s.collection.FindOne(ctx, filter, opts).Decode(&msg)
		if err == mongo.ErrNoDocuments {
			continue
		}
		if err != nil {
			return nil, err
		}

		result[chatID] = &msg
	}

	return result, nil
}

func (s *Storage) GetByID(ctx context.Context, msgID int64) (*domain.Message, error) {
	var msg domain.Message
	err := s.collection.FindOne(ctx, bson.M{"_id": msgID}).Decode(&msg)
	return &msg, err
}

func (s *Storage) Update(ctx context.Context, msgID int64, update bson.M) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": msgID}, update)
	return err
}

func (s *Storage) Count(ctx context.Context, filter bson.M) (int32, error) {
	count, err := s.collection.CountDocuments(ctx, filter)
	return int32(count), err
}
