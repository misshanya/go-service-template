package user

import (
	"context"
	"errors"
	"go-service-template/internal/errorz"
	"go-service-template/internal/models"
	v1 "go-service-template/internal/transport/http/v1"
)

func (h *UserHandler) RegisterUser(ctx context.Context, req v1.RegisterUserRequestObject) (v1.RegisterUserResponseObject, error) {
	input := &models.CreateUser{
		Username:       req.Body.Username,
		Password:       req.Body.Password,
		HashedPassword: req.Body.Password,
		Role:           "user",
	}
	u, token, err := h.svc.Create(ctx, input)
	if err != nil {
		if errors.Is(err, errorz.ErrUserAlreadyExists) {
			return v1.RegisterUser409JSONResponse{
				ConflictJSONResponse: v1.ConflictJSONResponse{
					Message: "user with the given username already exists",
					Code:    "USER_ALREADY_EXISTS",
				},
			}, nil
		}
		return nil, err
	}

	return v1.RegisterUser201JSONResponse{
		User:  toV1User(u),
		Token: token,
	}, nil
}
