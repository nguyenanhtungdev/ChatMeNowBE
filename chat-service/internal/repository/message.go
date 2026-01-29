package repository

import (
	"context"

	"github.com/chatmenow/chat-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(ctx context.Context, msg *model.Message) error
	FindByConversation(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]model.Message, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Message, error)
	Update(ctx context.Context, msg *model.Message) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, msg *model.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *messageRepository) FindByConversation(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]model.Message, error) {
	var messages []model.Message

	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *messageRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	var message model.Message
	err := r.db.WithContext(ctx).First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) Update(ctx context.Context, msg *model.Message) error {
	return r.db.WithContext(ctx).Save(msg).Error
}

func (r *messageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Message{}, "id = ?", id).Error
}
