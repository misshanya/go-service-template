package user

import (
	"context"
	"database/sql"
	"errors"
	"go-service-template/internal/errorz"
	"go-service-template/internal/models"
	"go-service-template/internal/repository"

	"github.com/google/uuid"
)

func (r *repo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	q := repository.Queries(ctx, r.q)

	u, err := q.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.ErrUserNotFound
		}
		return nil, err
	}

	return toModelUser(u), nil
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	q := repository.Queries(ctx, r.q)

	u, err := q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.ErrUserNotFound
		}
		return nil, err
	}

	return toModelUser(u), nil
}
