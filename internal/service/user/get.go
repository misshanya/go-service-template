package user

import (
	"context"
	"fmt"
	"go-service-template/internal/errorz"
	"go-service-template/internal/models"
	"go-service-template/pkg/crypto"

	"github.com/google/uuid"
)

func (s *service) Login(ctx context.Context, username string, password string) (*models.User, string, error) {
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}

	ok, err := crypto.ComparePasswordAndHash(password, u.HashedPassword)
	if err != nil {
		return nil, "", err
	}
	if !ok {
		return nil, "", errorz.ErrInvalidCredentials
	}

	token, _, err := s.jwtManager.GenerateToken(u.ID, u.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return u, token, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}
