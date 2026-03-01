package service

import (
	"context"

	"github.com/itmo-lite-chat/messages_svc/internal/domain"
)

// Оставляем интерфейс "на вырост"
type MessageStorage interface {
	Save(ctx context.Context, msg *domain.Message) error
}

type MessageService struct {
	storage MessageStorage
	// centrifugo RealTimeClient <- пока закомментировано
}

// Конструктор пока принимает только сторадж
func NewService(st MessageStorage) *MessageService {
	return &MessageService{
		storage: st,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, msg *domain.Message) error {
	// 1. Сохраняем в БД
	if err := s.storage.Save(ctx, msg); err != nil {
		return err
	}

	// TODO: Добавить отправку в Centrifugo, когда будет готов rtClient

	return nil
}
