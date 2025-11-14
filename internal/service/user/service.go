package user

import (
	"context"
	"go-service-template/internal/models"
)

type repo interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
}

type service struct {
	repo repo
}

func New(repo repo) *service {
	return &service{repo}
}
