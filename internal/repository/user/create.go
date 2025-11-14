package user

import (
	"context"
	"go-service-template/internal/db/sqlc/storage"
	"go-service-template/internal/errorz"
	"go-service-template/internal/models"
	"go-service-template/pkg/dberrors"
)

func (r *repo) Create(ctx context.Context, u *models.User) (*models.User, error) {
	newUser, err := r.queries.CreateUser(ctx, storage.CreateUserParams{
		Username: u.Username,
		Age:      u.Age,
	})
	if err != nil {
		if dberrors.IsUniqueViolation(err) {
			return nil, errorz.UserAlreadyExists
		}
		return nil, err
	}

	return &models.User{
		ID:       newUser.ID,
		Username: newUser.Username,
		Age:      newUser.Age,
	}, nil
}
