package user

import (
	"context"
	"go-service-template/internal/models"
	v1 "go-service-template/internal/transport/http/v1"

	"github.com/google/uuid"
)

type service interface {
	Create(ctx context.Context, input *models.CreateUser) (*models.User, string, error)
	Login(ctx context.Context, username string, password string) (*models.User, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type UserHandler struct {
	svc service
}

func New(svc service) *UserHandler {
	return &UserHandler{svc: svc}
}

func toV1User(u *models.User) v1.User {
	return v1.User{
		Id:       u.ID,
		Username: u.Username,
		Role:     v1.UserRole(u.Role),
	}
}
