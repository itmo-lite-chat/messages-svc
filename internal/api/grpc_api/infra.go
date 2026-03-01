package grpc_api

import (
	"context"

	"github.com/itmo-lite-chat/messages_svc/internal/domain"
)

// Service — описываем интерфейс того, что нам нужно от бизнес-логики.
// Это позволит API не зависеть от конкретной реализации сервиса.
type Service interface {
	SendMessage(ctx context.Context, msg *domain.Message) error
}
