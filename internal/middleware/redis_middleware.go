package middleware

import (
	"drivex/internal/redis"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit ( redisClient *redis.Client , limit int ,window time.Duration ) gin.HandlerFunc{

	return func (c *gin.Context){
		ip := c.ClientIP()

		key := fmt.Sprintf(
			"rate_limit:%s",
			ip,
		)
		ctx := c.Request.Context()
		count, err := redisClient.RDB.Incr(ctx, key).Result()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter unavailable",
			})
			c.Abort()
			return
		}

		if count == 1 {
			redisClient.RDB.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}
		c.Next()

	}

}