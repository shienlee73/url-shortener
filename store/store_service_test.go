package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var testStoreService = &StorageService{}

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	storm, err := storm.Open("url-shortener.db")
	if err != nil {
		panic(fmt.Sprintf("Error init BoltDB: %v", err))
	}

	testStoreService = NewStorageService(
		redisClient,
		storm,
		WithContext(ctx),
		WithCacheDuration(5*time.Minute),
	)
}

func TestStoreInit(t *testing.T) {
	assert.NotEmpty(t, testStoreService.redisClient)
}

func TestInsertionAndRetrieval(t *testing.T) {
	originalUrl := "https://www.eddywm.com/lets-build-a-url-shortener-in-go-with-redis-part-2-storage-layer/"
	userId := "e0dba740-fc4b-4977-872c-d360239e6b1a"
	shortURL := "Jsz4k57oAX"

	testStoreService.SaveUrlMapping(shortURL, originalUrl, userId)
	retrievedUrl, err := testStoreService.RetrieveOriginalUrl(shortURL)

	assert.NoError(t, err)
	assert.Equal(t, originalUrl, retrievedUrl)
}
