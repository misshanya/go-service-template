package user

import (
	"context"
	"go-service-template/internal/models"
)

func (s *service) Create(ctx context.Context, u *models.User) (*models.User, error) {
	return s.repo.Create(ctx, u)
}
