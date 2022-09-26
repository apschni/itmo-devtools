package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	apperrors "homeworkdeliverysystem/errors"
	"log"
	"time"
)

type TokenRepository struct {
	RedisClient *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) *TokenRepository {
	return &TokenRepository{RedisClient: redisClient}
}

func (t *TokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := t.RedisClient.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}
	return nil
}

func (t *TokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	result := t.RedisClient.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}

	if result.Val() < 1 {
		log.Printf("Refresh token to redis for userID/tokenID: %s/%s does not exist\n", userID, tokenID)
		return apperrors.NewAuthorization("Invalid refresh token")
	}

	return nil
}

func (t *TokenRepository) DeleteUserRefreshToken(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s*", userID)

	iter := t.RedisClient.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := t.RedisClient.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("Failed to delete refresh token: %s\n", iter.Val())
			failCount++
		}
	}

	if err := iter.Err(); err != nil {
		log.Printf("Failed to delete refresh token: %s\n", iter.Val())
	}

	if failCount > 0 {
		return apperrors.NewInternal()
	}

	return nil
}
