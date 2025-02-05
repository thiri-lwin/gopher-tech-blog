package server

// import (
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt"
// 	"github.com/thiri-lwin/gopher-tech-blog/internal/config"
// 	"github.com/thiri-lwin/gopher-tech-blog/internal/repo/redis"
// )

// const (
// 	rateLimit                = 15
// 	rateLimitSendMessage     = 2
// 	rateLimitDuration        = time.Minute
// 	rateLimitMessageDuration = time.Hour
// )

// func rateLimitMW(c *gin.Context) {
// 	handleRateLimit(c, "rate_limit:", rateLimit, rateLimitDuration, true)
// }

// func rateLimitSendMessageMW(c *gin.Context) {
// 	handleRateLimit(c, "rate_limit_send_message:", rateLimitSendMessage, rateLimitMessageDuration, false)
// }

// func handleRateLimit(c *gin.Context, keyPrefix string, limit int, duration time.Duration, errTmpl bool) {
// 	ctx := c.Request.Context()
// 	ip := c.ClientIP()
// 	key := keyPrefix + ip

// 	// Increment request count and set expiration
// 	count, err := redis.RDB.Incr(ctx, key).Result()
// 	if err != nil {
// 		if errTmpl {
// 			renderError(c, "Something went wrong. Please try again later.")
// 			c.Abort()
// 			return
// 		}
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong. Please try again later."})
// 		c.Abort()
// 		return
// 	}

// 	if count == 1 {
// 		// Set expiration time for the key on the first request
// 		redis.RDB.Expire(ctx, key, duration)
// 	}
// 	// Limit requests
// 	if count > int64(limit) {
// 		if errTmpl {
// 			renderError(c, "Too many requests, please slow down.")
// 			c.Abort()
// 			return
// 		}
// 		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, please slow down."})
// 		c.Abort()
// 		return
// 	}

// 	c.Next()
// }

// func renderError(c *gin.Context, subheading string) {
// 	data := gin.H{
// 		"BackgroundImage": fmt.Sprintf("%s/error-bg.jpg", config.ImageURL),
// 		"Heading":         "Error",
// 		"Subheading":      subheading,
// 		"Status":          "You've hit the limit. Contact us for unlimited access!",
// 	}
// 	tmpl.ExecuteTemplate(c.Writer, "error.html", data)
// 	c.Abort()
// }

// // CheckAuth is middleware to check if the JWT token is valid
// func CheckAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Get the token from the cookie
// 		tokenString, err := c.Cookie("jwt_token")
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
// 			c.Abort()
// 			return
// 		}

// 		jwtSecret := ""

// 		// Parse and validate the JWT token
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			// Ensure the token's signature method is HMAC
// 			if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
// 				return nil, fmt.Errorf("unexpected signing method")
// 			}
// 			return jwtSecret, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
// 			c.Abort()
// 			return
// 		}

// 		// Optionally, you can extract user information from the JWT claims
// 		claims, _ := token.Claims.(jwt.MapClaims)
// 		c.Set("user_id", claims["user_id"])

// 		// Proceed with the request
// 		c.Next()
// 	}
// }
