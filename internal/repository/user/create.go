package user

import (
	"context"
	"go-service-template/internal/db/sqlc/storage"
	"go-service-template/internal/errorz"
	"go-service-template/internal/models"
	"go-service-template/internal/repository"
	"go-service-template/pkg/dberrors"
)

func (r *repo) Create(ctx context.Context, input *models.CreateUser) (*models.User, error) {
	q := repository.Queries(ctx, r.q)

	u, err := q.CreateUser(ctx, storage.CreateUserParams{
		Username: input.Username,
		Password: input.HashedPassword,
		Role:     input.Role,
	})
	if err != nil {
		if dberrors.IsUniqueViolation(err) {
			return nil, errorz.ErrUserAlreadyExists
		}
		return nil, err
	}

	return toModelUser(u), nil
}
