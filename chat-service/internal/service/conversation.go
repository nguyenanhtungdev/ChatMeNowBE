package service

import (
	"context"
	"fmt"

	"github.com/chatmenow/chat-service/internal/model"
	"github.com/chatmenow/chat-service/internal/repository"
	"github.com/google/uuid"
)

type ConversationService struct {
	repo repository.ConversationRepository
}

func NewConversationService(repo repository.ConversationRepository) *ConversationService {
	return &ConversationService{repo: repo}
}

func (s *ConversationService) Create(ctx context.Context, req *model.CreateConversationRequest, createdBy uuid.UUID) (*model.Conversation, error) {
	// Validate conversation type
	if req.Type != "direct" && req.Type != "group" {
		return nil, fmt.Errorf("invalid conversation type: %s", req.Type)
	}

	// Add creator to members if not present
	found := false
	for _, id := range req.MemberIDs {
		if id == createdBy {
			found = true
			break
		}
	}
	if !found {
		req.MemberIDs = append(req.MemberIDs, createdBy)
	}

	// For direct chat, must have exactly 2 members (after adding creator)
	if req.Type == "direct" && len(req.MemberIDs) != 2 {
		return nil, fmt.Errorf("direct conversation must have exactly 2 members (you + 1 other person)")
	}

	conv := &model.Conversation{
		Name:      req.Name,
		Type:      req.Type,
		CreatedBy: createdBy,
	}

	if err := s.repo.Create(ctx, conv, req.MemberIDs); err != nil {
		return nil, err
	}

	return conv, nil
}

func (s *ConversationService) GetByID(ctx context.Context, id uuid.UUID) (*model.Conversation, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ConversationService) GetByUser(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error) {
	return s.repo.FindByUser(ctx, userID)
}

func (s *ConversationService) GetMembers(ctx context.Context, conversationID uuid.UUID) ([]model.ConversationMember, error) {
	return s.repo.GetMembers(ctx, conversationID)
}

func (s *ConversationService) AddMember(ctx context.Context, conversationID, userID uuid.UUID, role string) error {
	member := &model.ConversationMember{
		ConversationID: conversationID,
		UserID:         userID,
		Role:           role,
	}
	return s.repo.AddMember(ctx, member)
}

func (s *ConversationService) RemoveMember(ctx context.Context, conversationID, userID uuid.UUID) error {
	return s.repo.RemoveMember(ctx, conversationID, userID)
}

func (s *ConversationService) Update(ctx context.Context, conv *model.Conversation) error {
	return s.repo.Update(ctx, conv)
}
