package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/redis/go-redis/v9"
	"github.com/shienlee73/url-shortener/handler"
	"github.com/shienlee73/url-shortener/rate_limiter"
	"github.com/shienlee73/url-shortener/store"
)

func main() {
	addr, port, redisAddr, redisPassword := parseFlags()

	// redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
	ctx :=context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}
	defer redisClient.Close()

	// bolt
	storm, err := storm.Open("url-shortener.db")
	if err != nil {
		panic(fmt.Sprintf("Error init BoltDB: %v", err))
	}
	defer storm.Close()

	// store
	storageService := store.NewStorageService(
		redisClient,
		storm,
		store.WithContext(ctx),
		store.WithCacheDuration(5*time.Minute),
	)

	// rate limiter
	rateLimiter := rate_limiter.NewRateLimiter(
		redisClient,
	)

	// server
	server := handler.NewServer(storageService, rateLimiter)
	if err := server.Start(fmt.Sprintf("%s:%d", addr, port)); err != nil {
		panic(err)
	}
}

func parseFlags() (string, int, string, string) {
	addr := flag.String("address", "127.0.0.1", "address to listen on")
	port := flag.Int("port", 8080, "port to listen on")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis address")
	redisPassword := flag.String("redis-password", "", "Redis password")

	flag.StringVar(addr, "a", "127.0.0.1", "address to listen on (short)")
	flag.IntVar(port, "p", 8080, "port to listen on (short)")
	flag.StringVar(redisAddr, "r", "localhost:6379", "Redis address (short)")
	flag.StringVar(redisPassword, "rp", "", "Redis password (short)")

	flag.Parse()

	return *addr, *port, *redisAddr, *redisPassword
}