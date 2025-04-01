package auth

import (
	"chatY-go/internal/api/grpc/auth/proto"
	auth2 "chatY-go/internal/domain/auth"
	"context"
)

type AuthHandler struct {
	proto.UnimplementedAuthServiceServer
	service auth2.IAuthService
}

func NewAuthHandler(service auth2.IAuthService) proto.AuthServiceServer {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	ok, msg := h.service.Login(req.Username, req.Password)
	return &proto.LoginResponse{
		Success: ok,
		Message: msg,
	}, nil
}
