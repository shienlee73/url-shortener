package store

import (
	"context"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/redis/go-redis/v9"
)

type StorageService struct {
	redisClient   *redis.Client
	storm         *storm.DB
	ctx           context.Context
	CacheDuration time.Duration
}

func NewStorageService(redisClient *redis.Client, storm *storm.DB, options ...func(*StorageService)) *StorageService {
	storageService := &StorageService{
		redisClient:   redisClient,
		storm:         storm,
		ctx:           context.Background(),
		CacheDuration: 5 * time.Minute,
	}

	for _, option := range options {
		option(storageService)
	}

	return storageService
}
