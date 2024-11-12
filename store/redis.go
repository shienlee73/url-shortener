package store

import (
	"context"
	"fmt"
	"time"
)

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

func (s *StorageService) SaveUrlMappingToRedis(shortUrl string, originalUrl string, urlMappingId string) error {
	err := s.redisClient.HSet(s.ctx, shortUrl, "urlMappingId", urlMappingId, "originalUrl", originalUrl).Err()
	if err != nil {
		return fmt.Errorf("failed to save url mapping: %v", err)
	}
	err = s.redisClient.Expire(s.ctx, shortUrl, s.CacheDuration).Err()
	if err != nil {
		return fmt.Errorf("failed to save url mapping: %v", err)
	}
	return nil
}

func (s *StorageService) RetrieveUrlMappingFromRedis(shortUrl string) (map[string]string, error) {
	result, err := s.redisClient.HGetAll(s.ctx, shortUrl).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve url mapping: %v", err)
	}
	if len(result) == 0 {
        return nil, fmt.Errorf("no mapping found for shortUrl: %s", shortUrl)
    }
	return result, nil
}
