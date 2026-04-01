package jwt

import (
	"room_service/internal/domain"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret           []byte
	accessExpiry     time.Duration
}


func NewJWTService(secret string, accessExpiry time.Duration) *JWTService {
	return &JWTService{
		secret:           []byte(secret),
		accessExpiry:     accessExpiry,
	}
}

func (j *JWTService) GenerateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":         user.ID,
		"role":       user.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	return signed, err
}

func (j *JWTService) ValidateToken(tokenStr string) (*domain.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}
	user := &domain.User{}
	if id, ok := claims["sub"].(string); ok {
		user.ID = id
	} else if idf, ok := claims["sub"].(float64); ok {
		user.ID = fmt.Sprintf("%d", int64(idf))
	}
	if role, ok := claims["sub"].(string); ok {
		user.Role = role
	}
	return user, nil
}
