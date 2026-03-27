package user

import (
	"context"
	"errors"
	"go-service-template/internal/errorz"
	"go-service-template/internal/transport/http/middleware"
	v1 "go-service-template/internal/transport/http/v1"

	"github.com/google/uuid"
)

func (h *UserHandler) Login(ctx context.Context, req v1.LoginRequestObject) (v1.LoginResponseObject, error) {
	u, token, err := h.svc.Login(ctx, req.Body.Username, req.Body.Password)
	if err != nil {
		if errors.Is(err, errorz.ErrUserNotFound) || errors.Is(err, errorz.ErrInvalidCredentials) {
			return v1.Login401JSONResponse{
				UnauthorizedJSONResponse: v1.UnauthorizedJSONResponse{
					Message: "invalid username or password",
					Code:    "INVALID_CREDENTIALS",
				},
			}, nil
		}
		return nil, err
	}

	return v1.Login200JSONResponse{
		User:  toV1User(u),
		Token: token,
	}, nil
}

func (h *UserHandler) GetCurrentUser(ctx context.Context, req v1.GetCurrentUserRequestObject) (v1.GetCurrentUserResponseObject, error) {
	userID := ctx.Value(middleware.UserIDContextKey).(uuid.UUID)

	u, err := h.svc.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errorz.ErrUserNotFound) {
			return v1.GetCurrentUser404JSONResponse{
				NotFoundJSONResponse: v1.NotFoundJSONResponse{
					Message: "user not found",
					Code:    "USER_NOT_FOUND",
				},
			}, nil
		}
		return nil, err
	}

	return v1.GetCurrentUser200JSONResponse{
		Id:       u.ID,
		Role:     v1.UserRole(u.Role),
		Username: u.Username,
	}, nil
}
