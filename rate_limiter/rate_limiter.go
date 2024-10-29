package rate_limiter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	rateLimitKey = "rate_limit:%s:%s"
)

type RateLimiter struct {
	redisClient     *redis.Client
}

func NewRateLimiter(redisClient *redis.Client) *RateLimiter {
	rateLimiter := &RateLimiter{
		redisClient: redisClient,
	}

	return rateLimiter
}

func (r *RateLimiter) Limit(endpoint string, rateLimitTTL time.Duration, rateLimitTokens int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf(rateLimitKey, endpoint, ip)

		tokens, err := r.redisClient.Get(c, key).Int64()
		if err != nil {
			// Initialize the rate limit for the IP address
			fmt.Println(rateLimitTokens - 1)
			r.redisClient.Set(c, key, rateLimitTokens - 1, rateLimitTTL)
			c.Next()
			return
		}

		if tokens > 0 {
			// Decrement the token count
			r.redisClient.Decr(c, key)
			c.Next()
			return
		}

		c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		c.Abort()
	}
}
