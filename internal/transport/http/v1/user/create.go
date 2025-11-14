package user

import (
	"errors"
	"go-service-template/internal/errorz"
	"go-service-template/internal/models"
	"go-service-template/internal/transport/http/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) Create(c echo.Context) error {
	var body dto.UserCreateRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, &dto.HTTPStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, &dto.HTTPStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	u := &models.User{
		Username: body.Username,
		Age:      body.Age,
	}
	user, err := h.service.Create(c.Request().Context(), u)
	if err != nil {
		if errors.Is(err, errorz.UserAlreadyExists) {
			return c.JSON(http.StatusConflict, &dto.HTTPStatus{
				Code:    http.StatusConflict,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, &dto.HTTPStatus{
			Code:    http.StatusInternalServerError,
			Message: errorz.InternalServerError.Error(),
		})
	}

	resp := &dto.UserCreateResponse{
		ID:       user.ID,
		Username: user.Username,
		Age:      user.Age,
	}
	return c.JSON(http.StatusCreated, resp)
}
