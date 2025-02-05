package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thiri-lwin/gopher-tech-blog/internal/pkg/jwt"
)

type UserCreds struct {
	Email    string
	Password string
}

// ServeSignIn serves the signin page
func (h Handler) ServeSignIn(c *gin.Context) {
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/contact-bg.jpg", h.cfg.ImageURL), // update image
		"Heading":         "Welcome back!",
		"Subheading":      "Please sign in to access your account",
		"IsAuthenticated": false,
	}
	h.tmpl.ExecuteTemplate(c.Writer, "signin.html", data)
}

func (h Handler) SignIn(c *gin.Context) {
	var creds UserCreds
	if err := c.ShouldBindJSON(&creds); err != nil {
		log.Println("Failed to bind form data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// // Validate user credentials
	// storedPassword, exists := users[creds.Username]
	// if !exists {
	// 	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	// 	return
	// }

	// if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
	// 	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	// 	return
	// }

	tokenString, err := jwt.GenerateJWT(creds.Email, h.cfg.JWTKey)
	if err != nil {
		log.Println("Failed to generate JWT", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to signin"})
		return
	}

	// Set JWT in a cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute * time.Duration(h.cfg.JWTExpirationTime)),
	})
	c.JSON(http.StatusOK, gin.H{"message": "Sign-in successful"})
}

func (h Handler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(), // Expire immediately
	})
	c.Redirect(http.StatusPermanentRedirect, "/")
}
