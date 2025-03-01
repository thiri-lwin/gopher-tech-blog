package middleware

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thiri-lwin/gopher-tech-blog/internal/config"
	pkg_jwt "github.com/thiri-lwin/gopher-tech-blog/internal/pkg/jwt"
	"github.com/thiri-lwin/gopher-tech-blog/internal/repo/redis"
)

type RequestMeta struct {
	UserID    int
	Email     string
	FirstName string
	LastName  string
}

const (
	rateLimit                = 15
	rateLimitSendMessage     = 2
	rateLimitDuration        = time.Minute
	rateLimitMessageDuration = time.Hour
)

func RateLimitMW(c *gin.Context, tmpl *template.Template) {
	handleRateLimit(c, "rate_limit:", rateLimit, rateLimitDuration, tmpl)
}

func RateLimitSendMessageMW(c *gin.Context) {
	handleRateLimit(c, "rate_limit_send_message:", rateLimitSendMessage, rateLimitMessageDuration, nil)
}

func handleRateLimit(c *gin.Context, keyPrefix string, limit int, duration time.Duration, tmpl *template.Template) {
	ctx := c.Request.Context()
	ip := c.ClientIP()
	key := keyPrefix + ip

	// Increment request count and set expiration
	count, err := redis.RDB.Incr(ctx, key).Result()
	if err != nil {
		if tmpl != nil {
			renderError(c, tmpl, "Something went wrong. Please try again later.")
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
		if tmpl != nil {
			renderError(c, tmpl, "Too many requests, please slow down.")
			c.Abort()
			return
		}
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, please slow down."})
		c.Abort()
		return
	}

	c.Next()
}

func renderError(c *gin.Context, tmpl *template.Template, subheading string) {
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/error-bg.jpg", config.ImageURL),
		"Heading":         "Error",
		"Subheading":      subheading,
		"Status":          "You've hit the limit. Contact us for unlimited access!",
	}
	tmpl.ExecuteTemplate(c.Writer, "error.html", data)
	c.Abort()
}

func AuthMiddleware(jwtKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the cookie
		cookie, err := c.Cookie("jwt_token")
		if err != nil {
			c.Set("request_meta", RequestMeta{})
			c.Next()
			return
		}

		// Parse and validate the JWT token
		claims := &pkg_jwt.Claims{}
		token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if err != nil || !token.Valid {
			c.Set("request_meta", RequestMeta{})
			c.Next()
			return
		}

		// Set the username in the context
		c.Set("request_meta", RequestMeta{
			UserID:    claims.UserID,
			Email:     claims.Email,
			FirstName: claims.FirstName,
			LastName:  claims.LastName,
		})
		c.Next()
	}
}

func GetRequestMeta(c *gin.Context) RequestMeta {
	val, ok := c.Get("request_meta")
	if !ok {
		return RequestMeta{}
	}

	meta, ok := val.(RequestMeta)
	if !ok {
		return RequestMeta{}
	}

	return meta
}
