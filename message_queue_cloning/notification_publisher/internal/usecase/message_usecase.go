package usecase

import (
	"context"
	"notification_publisher/internal/entity"
	"notification_publisher/internal/repository"
)

type MessageUseCase struct {
	repo *repository.RabbitMQRepository
}

func NewMessageUseCase(repo *repository.RabbitMQRepository) *MessageUseCase {
	return &MessageUseCase{repo: repo}
}

func (u *MessageUseCase) PublishMessage(ctx context.Context, routingKey string, message entity.Message) error {
	return u.repo.PublishMessage(ctx, routingKey, message)
}
