package user

import (
	"context"
	"fmt"
	"go-service-template/internal/models"
	"go-service-template/pkg/crypto"
)

func (s *service) Create(ctx context.Context, input *models.CreateUser) (*models.User, string, error) {
	hash, err := crypto.GenerateHash(input.Password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}
	input.HashedPassword = hash

	u, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, "", err
	}

	token, _, err := s.jwtManager.GenerateToken(u.ID, u.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return u, token, nil
}
