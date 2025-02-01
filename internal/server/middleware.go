package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thiri-lwin/gopher-tech-blog/internal/repo/redis"
)

// Middleware for rate limiting
func rateLimitMiddleware(c *gin.Context) {
	ctx := c.Request.Context()
	ip := c.ClientIP()
	key := "rate_limit:" + ip

	// Increment request count and set expiration
	count, err := redis.RDB.Incr(ctx, key).Result()
	if err != nil {
		data := gin.H{
			"BackgroundImage": "/static/img/error-bg.jpg",
			"Heading":         "Error",
			"Subheading":      "Something went wrong. Please try again later.",
		}
		tmpl.ExecuteTemplate(c.Writer, "error.html", data)
		c.Abort()
		return
	}

	if count == 1 {
		// Set expiration time for the key on the first request
		redis.RDB.Expire(ctx, key, time.Minute)
	}

	// Limit to 5 requests per minute
	if count > 5 {
		data := gin.H{
			"BackgroundImage": "/static/img/error-bg.jpg",
			"Heading":         "Error",
			"Subheading":      "Too many requests, please slow down.",
			"Status":          "You've hit the limit. Contact us for unlimited access!",
		}
		tmpl.ExecuteTemplate(c.Writer, "error.html", data)
		c.Abort()
		return
	}

	c.Next()
}
