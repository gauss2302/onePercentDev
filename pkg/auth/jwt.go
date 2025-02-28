package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTManager struct {
	config *JWTConfig
}

func NewJWTManger(cfg *JWTConfig) *JWTManager {
	return &JWTManager{
		config: cfg,
	}
}

func (m *JWTManager) GenerateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(m.config.AccessExpire).Unix(),
	})

	return token.SignedString([]byte(m.config.AccessSecret))
}

func (m *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(m.config.RefreshExpire).Unix(),
	})

	return token.SignedString([]byte(m.config.RefreshSecret))
}

func (m *JWTManager) VerifyAccessToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.config.AccessSecret), nil
	})
}

func (m *JWTManager) VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.config.RefreshSecret), nil
	})
}
