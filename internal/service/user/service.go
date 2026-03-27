package user

import (
	"context"
	"go-service-template/internal/models"

	"github.com/google/uuid"
)

type repo interface {
	Create(ctx context.Context, input *models.CreateUser) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type jwtManager interface {
	GenerateToken(userID uuid.UUID, role string) (string, int64, error)
}

type service struct {
	repo       repo
	jwtManager jwtManager
}

func New(repo repo, jwtManager jwtManager) *service {
	return &service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}
