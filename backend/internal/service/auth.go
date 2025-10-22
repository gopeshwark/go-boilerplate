package service

import (
	"github.com/gopeshwark/go-boilerplate/internal/server"

	"github.com/clerk/clerk-sdk-go/v2"
)

type AuthService struct {
	servcer *server.Server
}

func NewAuthService(s *server.Server) *AuthService {
	clerk.SetKey(s.Config.Auth.SecretKey)
	return &AuthService{
		servcer: s,
	}
}
