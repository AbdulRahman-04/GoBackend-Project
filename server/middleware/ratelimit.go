package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitMiddleware creates independent rate limiter for each route+method
func RateLimitMiddleware(limit int) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  int64(limit),
	}

	store := memory.NewStore()

	return func(c *gin.Context) {
		// Unique key per route + method => so that each route is rate-limited independently
		key := c.FullPath() + "-" + c.Request.Method

		// Create new limiter for this route
		instance := limiter.New(store, rate)

		// Get context for this key
		context, err := instance.Get(c, key)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"msg": "Rate limiter error"})
			return
		}

		// If over limit
		if context.Reached {
			c.AbortWithStatusJSON(429, gin.H{
				"msg": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
