package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type StorageService struct {
	redisClient   *redis.Client
	ctx           context.Context
	CacheDuration time.Duration
}

func WithContext(ctx context.Context) func(*StorageService) {
	return func(s *StorageService) {
		s.ctx = ctx
	}
}

func WithCacheDuration(duration time.Duration) func(*StorageService) {
	return func(s *StorageService) {
		s.CacheDuration = duration
	}
}

func NewStorageService(redisClient *redis.Client, options ...func(*StorageService)) *StorageService {
	storageService := &StorageService{}

	storageService.redisClient = redisClient

	for _, option := range options {
		option(storageService)
	}

	return storageService
}

func (s *StorageService) SaveUrlMapping(shortUrl string, originalUrl string, userId string) error {
	err := s.redisClient.Set(s.ctx, shortUrl, originalUrl, s.CacheDuration).Err()
	if err != nil {
		return fmt.Errorf("failed to save url mapping: %v", err)
	}
	return nil
}

func (s *StorageService) RetrieveOriginalUrl(shortUrl string) (string, error) {
	result, err := s.redisClient.Get(s.ctx, shortUrl).Result()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve original url: %v", err)
	}
	return result, nil
}
