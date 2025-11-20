package user

import (
	"errors"
	"go-service-template/internal/errorz"
	"go-service-template/internal/models"
	"go-service-template/internal/transport/http/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Create handles the creation of a user
//
//	@Summary		Create a new user
//	@Description	Creates a new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UserCreateRequest	true	"User data"
//	@Success		201		{object}	dto.UserCreateResponse
//	@Failure		400		{object}	dto.HTTPStatus	"Invalid request body"
//	@Failure		409		{object}	dto.HTTPStatus	"User already exists"
//	@Failure		500		{object}	dto.HTTPStatus	"Internal server error"
//	@Router			/user [post]
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
