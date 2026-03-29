package usecase

import (
	"context"
	"notification_consumer/internal/repository"
)

type MessageUseCase struct {
	repo *repository.RabbitMQRepository
}

func NewMessageUseCase(repo *repository.RabbitMQRepository) *MessageUseCase {
	return &MessageUseCase{repo: repo}
}

func (u *MessageUseCase) ConsumeMessagesSms(ctx context.Context, serviceName string) error {
	return u.repo.ConsumeMessagesSms(ctx, serviceName)
}

func (u *MessageUseCase) ConsumeMessagesFcm(ctx context.Context, serviceName string) error {
	return u.repo.ConsumeMessagesFcm(ctx, serviceName)
}

func (u *MessageUseCase) ConsumeMessagesEmail(ctx context.Context, serviceName string) error {
	return u.repo.ConsumeMessagesEmail(ctx, serviceName)
}
