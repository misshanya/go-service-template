package middleware

import (
	"go-service-template/internal/errorz"
	v1 "go-service-template/internal/transport/http/v1"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func StrictErrorMiddleware(f v1.StrictHandlerFunc, operationID string) v1.StrictHandlerFunc {
	return func(ctx echo.Context, request interface{}) (response interface{}, err error) {
		res, err := f(ctx, request)
		if err != nil {
			slog.Error("operation failed", "operation_id", operationID, "error", err)
			return nil, echo.NewHTTPError(http.StatusInternalServerError, errorz.ErrInternalServerError.Error())
		}
		return res, nil
	}
}
