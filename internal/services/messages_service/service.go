package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/itmo-lite-chat/messages_svc/internal/domain"
	"github.com/itmo-lite-chat/messages_svc/internal/services/messages_service/clients"
	"go.mongodb.org/mongo-driver/bson"
)

type MessageService struct {
	storage clients.MessageStorage
	node    *snowflake.Node
}

func NewService(st clients.MessageStorage) *MessageService {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic("failed to create snowflake node in service: " + err.Error())
	}

	return &MessageService{
		storage: st,
		node:    node,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error) {
	if msg.MessageID == 0 {
		msg.MessageID = s.node.Generate().Int64()
	}

	if err := s.storage.Save(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageService) GetChatHistory(ctx context.Context, chatID uuid.UUID, limit int32, beforeMessageID *int64) ([]*domain.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	return s.storage.GetChatHistory(ctx, chatID, int64(limit), beforeMessageID)
}

func (s *MessageService) GetLastMessages(ctx context.Context, chatIDs []uuid.UUID) (map[uuid.UUID]*domain.Message, error) {
	return s.storage.GetLastMessages(ctx, chatIDs)
}

func (s *MessageService) EditMessage(ctx context.Context, msgID int64, senderID uuid.UUID, newContent domain.Content) (*domain.Message, error) {
	msg, err := s.storage.GetByID(ctx, msgID)
	if err != nil {
		return nil, err
	}

	if msg.SenderID != senderID {
		return nil, fmt.Errorf("permission denied")
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"content":    newContent,
			"updated_at": now,
		},
	}

	if err := s.storage.Update(ctx, msgID, update); err != nil {
		return nil, err
	}

	msg.Content = newContent
	msg.UpdatedAt = &now
	return msg, nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, msgID int64, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	msg, err := s.storage.GetByID(ctx, msgID)
	if err != nil {
		return err
	}
	if msg.SenderID != userUUID {
		return fmt.Errorf("permission denied")
	}

	update := bson.M{
		"$set": bson.M{"deleted_at": time.Now()},
	}
	return s.storage.Update(ctx, msgID, update)
}

func (s *MessageService) GetUnreadCount(ctx context.Context, chatID uuid.UUID, lastReadID int64) (int32, error) {
	filter := bson.M{
		"chat_id":    chatID,
		"_id":        bson.M{"$gt": lastReadID},
		"deleted_at": nil,
	}
	return s.storage.Count(ctx, filter)
}
