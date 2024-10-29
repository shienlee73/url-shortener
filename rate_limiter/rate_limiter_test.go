package rate_limiter

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var rateLimiter = &RateLimiter{}

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

	rateLimiter = NewRateLimiter(redisClient)
}

func afterEach() {
	rateLimiter.redisClient.FlushAll(context.Background())
}

func TestRateLimiterInit(t *testing.T) {
	assert.NotEmpty(t, rateLimiter.redisClient)
	afterEach()
}

func TestRateLimitInit(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		RemoteAddr: "192.168.1.100:8000",
	}
	rateLimiter.Limit("test_endpoint", time.Minute, 5)(c)
	assert.Equal(t, http.StatusOK, w.Code)

	remaining := rateLimiter.redisClient.Get(c, "rate_limit:test_endpoint:192.168.1.100")
	assert.Equal(t, "4", remaining.Val())
	afterEach()
}

func TestRateLimitDecrement(t *testing.T) {
	for i := 0; i <= 5; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			RemoteAddr: "192.168.1.100:12345",
		}

		rateLimiter.Limit("test_endpoint", time.Minute, 5)(c)
		if i < 5 {
			assert.Equal(t, http.StatusOK, w.Code, "Expected request to be allowed")
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Expected request to be rate-limited")
		}

		remaining := rateLimiter.redisClient.Get(c, "rate_limit:test_endpoint:192.168.1.100")
		expectedRemaining := fmt.Sprint(max(4-i, 0))
		assert.Equal(t, expectedRemaining, remaining.Val(), "Remaining requests should decrement correctly")
	}
	afterEach()
}
