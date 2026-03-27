package server

import (
	"context"
	"go-service-template/internal/transport/http/v1"
	"go-service-template/internal/transport/http/v1/user"
)

type Server struct {
	*user.UserHandler
}

func New(u *user.UserHandler) *Server {
	return &Server{
		UserHandler: u,
	}
}

func (s *Server) GetPing(ctx context.Context, request v1.GetPingRequestObject) (v1.GetPingResponseObject, error) {
	return v1.GetPing200JSONResponse{Status: "PONG_PUBLIC"}, nil
}

func (s *Server) GetAdminPing(ctx context.Context, request v1.GetAdminPingRequestObject) (v1.GetAdminPingResponseObject, error) {
	return v1.GetAdminPing200JSONResponse{Secret: "PONG_ADMIN"}, nil
}
