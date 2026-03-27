package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Manager struct {
	secret        []byte
	tokenDuration time.Duration
}

func NewManager(secret string, duration time.Duration) *Manager {
	return &Manager{
		secret:        []byte(secret),
		tokenDuration: duration,
	}
}

func (m *Manager) GenerateToken(userID uuid.UUID, role string) (string, int64, error) {
	now := time.Now()
	exp := now.Add(m.tokenDuration)

	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"role": role,
		"iat":  now.Unix(),
		"exp":  exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, int64(m.tokenDuration.Seconds()), nil
}

func (m *Manager) ParseToken(tokenString string) (uuid.UUID, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			return uuid.Nil, "", fmt.Errorf("invalid token claims: missing sub")
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, "", fmt.Errorf("invalid user ID in token: %w", err)
		}

		role, ok := claims["role"].(string)
		if !ok {
			return uuid.Nil, "", fmt.Errorf("invalid token claims: missing role")
		}

		return userID, role, nil
	}

	return uuid.Nil, "", fmt.Errorf("invalid token")
}
