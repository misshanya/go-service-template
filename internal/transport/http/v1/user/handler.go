package user

import (
	"context"
	"go-service-template/internal/models"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type service interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
}

type handler struct {
	service   service
	validator *validator.Validate
}

func New(s service) *handler {
	return &handler{
		service:   s,
		validator: validator.New(),
	}
}

func (h *handler) Setup(group *echo.Group) {
	group.POST("", h.Create)
}
