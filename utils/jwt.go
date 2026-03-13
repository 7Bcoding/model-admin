package utils

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("beta-llm-secret-1234") // 在生产环境中应该使用环境变量

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(userID int, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken 验证 JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

type contextKey string

const UserContextKey contextKey = "user"

// SetUserToContext 将用户信息添加到上下文
func SetUserToContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, UserContextKey, claims)
}

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value("user").(*Claims)
	if !ok {
		return nil, false
	}
	return claims, true
}
