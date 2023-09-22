package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimitMiddleware() gin.HandlerFunc {
	// Define rate limits
	rate := limiter.Rate{
		Limit:  10,        // Maximum number of requests allowed
		Period: time.Hour, // Time period for the rate limit
	}

	// Create a memory store for rate limiting
	store := memory.NewStore()

	// Create a rate limiter instance
	limiter := limiter.New(store, rate)

	// Return the Gin middleware handler
	return func(c *gin.Context) {
		// Use the rate limiter to check if the request is allowed
		contextKey := "contextKey"
		ctx, err := limiter.Get(c, contextKey)
		if err != nil {
			// Handle the error (e.g., log it)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		if ctx.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
