package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/thiri-lwin/gopher-tech-blog/internal/pkg/jwt"
	repo "github.com/thiri-lwin/gopher-tech-blog/internal/repo/postgres"
	"github.com/thiri-lwin/gopher-tech-blog/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

type UserCreds struct {
	Email    string
	Password string
}

// ServeSignIn serves the signin page
func (h Handler) ServeSignIn(c *gin.Context) {
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/home-bg.jpg", h.cfg.ImageURL), // update image
		"Heading":         "Welcome back!",
		"Subheading":      "Please sign in to share your thoughts and connect with others",
		"IsAuthenticated": false,
		"GoogleClientID":  h.cfg.GoogleClientID,
	}
	h.tmpl.ExecuteTemplate(c.Writer, "signin.html", data)
}

func (h Handler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()
	var creds UserCreds
	if err := c.ShouldBindJSON(&creds); err != nil {
		log.Println("Failed to bind form data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Validate user credentials
	user, err := h.repo.GetUser(ctx, creds.Email)
	if err != nil {
		log.Println("Failed to get user:", err)
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not fount"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Println("password mismatched")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	h.handleSuccessfulSignin(c, user)
}

func (h Handler) ServeSignUp(c *gin.Context) {
	data := gin.H{
		"BackgroundImage": fmt.Sprintf("%s/home-bg.jpg", h.cfg.ImageURL), // update image
		"Heading":         "Join Gopher Blog today!",
		"Subheading":      "Sign up now to share your thoughts and connect with others",
		"IsAuthenticated": false,
		"GoogleClientID":  h.cfg.GoogleClientID,
	}
	h.tmpl.ExecuteTemplate(c.Writer, "signup.html", data)
}

func (h Handler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()
	var user repo.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Failed to bind form data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Check if the user already exists
	_, err := h.repo.GetUser(ctx, user.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	} else if err != pgx.ErrNoRows {
		log.Println("Failed to check if user exists:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if user exists"})
		return
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("Failed to hash password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Create the user in the database
	if _, err := h.repo.CreateUser(ctx, user); err != nil {
		log.Println("Failed to create user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (h Handler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	c.Writer.Header().Set("Clear-Site-Data", `"cookies"`)

	c.Redirect(http.StatusSeeOther, "/")
}

func (h Handler) GoogleSignIn(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind form data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Verify Google ID token
	payload, err := idtoken.Validate(ctx, req.Token, h.cfg.GoogleClientID)
	if err != nil {
		log.Println("Failed to validate Google ID token:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Google ID token"})
		return
	}

	// Extract user info
	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	user := repo.User{
		Email:     email,
		FirstName: name,
	}

	// Check if the user already exists
	if dbUser, err := h.repo.GetUser(c.Request.Context(), email); err == nil {
		log.Println("User already exists")
		h.handleSuccessfulSignin(c, dbUser)
		return
	}

	// Create the user in the database
	userID, err := h.repo.CreateUser(ctx, user)
	if err != nil {
		log.Println("Failed to create user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	user.ID = userID
	h.handleSuccessfulSignin(c, user)
}

func (h Handler) handleSuccessfulSignin(c *gin.Context, user repo.User) {
	tokenString, err := jwt.GenerateJWT(user, h.cfg.JWTKey)
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
