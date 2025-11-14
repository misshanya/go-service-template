package dto

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Age      int32     `json:"age"`
}

type UserCreateRequest struct {
	Username string `json:"username" validate:"required"`
	Age      int32  `json:"age" validate:"required,gt=0,lt=200"`
}

type UserCreateResponse User
