package service

import (
	"context"

	"github.com/chatmenow/chat-service/internal/repository"
)

type PresenceService struct {
	redis *repository.RedisClient
}

func NewPresenceService(redis *repository.RedisClient) *PresenceService {
	return &PresenceService{redis: redis}
}

func (s *PresenceService) SetOnline(ctx context.Context, userID string) error {
	return s.redis.SetUserOnline(ctx, userID)
}

func (s *PresenceService) SetOffline(ctx context.Context, userID string) error {
	return s.redis.SetUserOffline(ctx, userID)
}

func (s *PresenceService) IsOnline(ctx context.Context, userID string) (bool, error) {
	return s.redis.IsUserOnline(ctx, userID)
}

func (s *PresenceService) StartTyping(ctx context.Context, conversationID, userID string) error {
	return s.redis.AddTyping(ctx, conversationID, userID)
}

func (s *PresenceService) StopTyping(ctx context.Context, conversationID, userID string) error {
	return s.redis.RemoveTyping(ctx, conversationID, userID)
}

func (s *PresenceService) GetTypingUsers(ctx context.Context, conversationID string) ([]string, error) {
	return s.redis.GetTypingUsers(ctx, conversationID)
}
