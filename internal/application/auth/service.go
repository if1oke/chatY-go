package auth

import "chatY-go/pkg/logger"

type AuthService struct {
	l logger.ILogger
}

func NewAuthService(l logger.ILogger) *AuthService { return &AuthService{l: l} }

func (service *AuthService) Login(username, password string) (bool, string) {
	service.l.Infof("Accept login request: %v:%v", username, password)
	return true, "ok"
}
