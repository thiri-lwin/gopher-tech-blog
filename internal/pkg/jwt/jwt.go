package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	repo "github.com/thiri-lwin/gopher-tech-blog/internal/repo/postgres"
)

type Claims struct {
	Email     string `json:"email"`
	UserID    int    `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}

func GenerateJWT(user repo.User, jwtKey string) (string, error) {
	// Generate JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email:     user.Email,
		UserID:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
