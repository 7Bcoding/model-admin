package utils

import (
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateLongTermToken generates a JWT token that expires in 100 years
func GenerateLongTermToken(userID int, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(100, 0, 0).Unix(), // 100 years from now
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func TestGenerateLongTermToken(t *testing.T) {
	token, err := GenerateLongTermToken(1, "api-user", "api-user")
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	log.Printf("Generated token: %s", token)
}
