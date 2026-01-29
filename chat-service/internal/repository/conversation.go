package repository

import (
	"context"

	"github.com/chatmenow/chat-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	Create(ctx context.Context, conv *model.Conversation, memberIDs []uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Conversation, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error)
	GetMembers(ctx context.Context, conversationID uuid.UUID) ([]model.ConversationMember, error)
	AddMember(ctx context.Context, member *model.ConversationMember) error
	RemoveMember(ctx context.Context, conversationID, userID uuid.UUID) error
	Update(ctx context.Context, conv *model.Conversation) error
}

type conversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) Create(ctx context.Context, conv *model.Conversation, memberIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create conversation
		if err := tx.Create(conv).Error; err != nil {
			return err
		}

		// Add members
		for i, memberID := range memberIDs {
			member := &model.ConversationMember{
				ConversationID: conv.ID,
				UserID:         memberID,
				Role:           "member",
			}

			// First member (creator) is admin
			if i == 0 || memberID == conv.CreatedBy {
				member.Role = "admin"
			}

			if err := tx.Create(member).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *conversationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.db.WithContext(ctx).
		Preload("Members").
		First(&conversation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *conversationRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error) {
	var conversations []model.Conversation

	err := r.db.WithContext(ctx).
		Joins("JOIN conversation_members ON conversation_members.conversation_id = conversations.id").
		Where("conversation_members.user_id = ? AND conversation_members.deleted_at IS NULL", userID).
		Preload("Members").
		Order("conversations.updated_at DESC").
		Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	return conversations, nil
}

func (r *conversationRepository) GetMembers(ctx context.Context, conversationID uuid.UUID) ([]model.ConversationMember, error) {
	var members []model.ConversationMember
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *conversationRepository) AddMember(ctx context.Context, member *model.ConversationMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *conversationRepository) RemoveMember(ctx context.Context, conversationID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Delete(&model.ConversationMember{}).Error
}

func (r *conversationRepository) Update(ctx context.Context, conv *model.Conversation) error {
	return r.db.WithContext(ctx).Save(conv).Error
}
