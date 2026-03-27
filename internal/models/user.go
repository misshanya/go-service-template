package models

import "github.com/google/uuid"

type CreateUser struct {
	Username       string
	Password       string
	HashedPassword string
	Role           string
}

type User struct {
	ID             uuid.UUID
	Username       string
	Role           string
	HashedPassword string
}
