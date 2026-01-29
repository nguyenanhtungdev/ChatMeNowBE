package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisClient{client: client}
}

func (r *RedisClient) SetUserOnline(ctx context.Context, userID string) error {
	return r.client.Set(ctx, "user:"+userID+":online", time.Now().Unix(), 0).Err()
}

func (r *RedisClient) SetUserOffline(ctx context.Context, userID string) error {
	return r.client.Del(ctx, "user:"+userID+":online").Err()
}

func (r *RedisClient) IsUserOnline(ctx context.Context, userID string) (bool, error) {
	exists, err := r.client.Exists(ctx, "user:"+userID+":online").Result()
	return exists > 0, err
}

func (r *RedisClient) AddTyping(ctx context.Context, conversationID, userID string) error {
	key := "conversation:" + conversationID + ":typing"
	return r.client.SAdd(ctx, key, userID).Err()
}

func (r *RedisClient) RemoveTyping(ctx context.Context, conversationID, userID string) error {
	key := "conversation:" + conversationID + ":typing"
	return r.client.SRem(ctx, key, userID).Err()
}

func (r *RedisClient) GetTypingUsers(ctx context.Context, conversationID string) ([]string, error) {
	key := "conversation:" + conversationID + ":typing"
	return r.client.SMembers(ctx, key).Result()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
