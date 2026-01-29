package service

import (
	"context"

	"github.com/chatmenow/chat-service/internal/model"
	"github.com/chatmenow/chat-service/internal/repository"
	"github.com/google/uuid"
)

type MessageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) Create(ctx context.Context, msg *model.Message) error {
	if msg.Type == "" {
		msg.Type = "text"
	}
	return s.repo.Create(ctx, msg)
}

func (s *MessageService) GetMessages(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]model.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.FindByConversation(ctx, conversationID, limit, offset)
}

func (s *MessageService) GetByID(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *MessageService) Update(ctx context.Context, msg *model.Message) error {
	return s.repo.Update(ctx, msg)
}

func (s *MessageService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
