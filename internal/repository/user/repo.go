package user

import "go-service-template/internal/db/sqlc/storage"

type repo struct {
	queries *storage.Queries
}

func New(queries *storage.Queries) *repo {
	return &repo{
		queries: queries,
	}
}
