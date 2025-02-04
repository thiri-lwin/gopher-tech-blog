package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thiri-lwin/gopher-tech-blog/internal/config"
	"github.com/thiri-lwin/gopher-tech-blog/internal/repo/redis"
)

const (
	rateLimit                = 15
	rateLimitSendMessage     = 2
	rateLimitDuration        = time.Minute
	rateLimitMessageDuration = time.Hour
)

func rateLimitMW(c *gin.Context) {
	handleRateLimit(c, "rate_limit:", rateLimit, rateLimitDuration, true)
}

func rateLimitSendMessageMW(c *gin.Context) {
	handleRateLimit(c, "rate_limit_send_message:", rateLimitSendMessage, rateLimitMessageDuration, false)
}

func handleRateLimit(c *gin.Context, keyPrefix string, limit int, duration time.Duration, errTmpl bool) {
	ctx := c.Request.Context()
	ip := c.ClientIP()
	key := keyPrefix + ip

	// Increment request count and set expiration
	count, err := redis.RDB.Incr(ctx, key).Result()
	if err != nil {
		if errTmpl {
			renderError(c, "Something went wrong. Please try again later.")
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong. Please try again later."})
		c.Abort()
		return
	}

	if count == 1 {
		// Set expiration time for the key on the first request
		redis.RDB.Expire(ctx, key, duration)
	}
	// Limit requests
	if count > int64(limit) {
		if errTmpl {
			renderError(c, "Too many requests, please slow down.")
			c.Abort()
			return
		}
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, please slow down."})
		c.Abort()
		return
	}

	c.Next()
}

func renderError(c *gin.Context, subheading string) {
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/error-bg.jpg", config.ImageURL),
		"Heading":         "Error",
		"Subheading":      subheading,
		"Status":          "You've hit the limit. Contact us for unlimited access!",
	}
	tmpl.ExecuteTemplate(c.Writer, "error.html", data)
	c.Abort()
}
