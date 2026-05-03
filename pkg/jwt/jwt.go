package jwt

import (
	"errors"
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	Role   int8  `json:"role"`
	TokenID string `json:"token_id"`
	jwt.RegisteredClaims
}

type Manager struct {
	accessSecret     []byte
	refreshSecret    []byte
	accessExpireHours time.Duration
	refreshExpireDays time.Duration
}

func NewManager(accessSecret, refreshSecret string, accessExpireHours, refreshExpireDays int) *Manager {
	return &Manager{
		accessSecret:      []byte(accessSecret),
		refreshSecret:     []byte(refreshSecret),
		accessExpireHours: time.Duration(accessExpireHours) * time.Hour,
		refreshExpireDays: time.Duration(refreshExpireDays) * 24 * time.Hour,
	}
}

func (m *Manager) GenerateTokenPair(userID int64, role int8) (accessToken, refreshToken, tokenID string, err error) {
	tokenID = uuid.New().String()
	now := time.Now()

	accessClaims := CustomClaims{
		UserID:  userID,
		Role:    role,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpireHours)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString(m.accessSecret)
	if err != nil {
		return "", "", "", fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := CustomClaims{
		UserID:  userID,
		Role:    role,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshExpireDays)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = rt.SignedString(m.refreshSecret)
	if err != nil {
		return "", "", "", fmt.Errorf("sign refresh token: %w", err)
	}

	return accessToken, refreshToken, tokenID, nil
}

func (m *Manager) ParseAccessToken(tokenString string) (*CustomClaims, error) {
	return m.parseToken(tokenString, m.accessSecret)
}

func (m *Manager) ParseRefreshToken(tokenString string) (*CustomClaims, error) {
	return m.parseToken(tokenString, m.refreshSecret)
}

func (m *Manager) parseToken(tokenString string, secret []byte) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
