package service

import (
	"auth-poc/svc/auth/ports"
)

type AuthService struct {
	repository ports.AuthRepository
}

func NewAuthService(AuthRepo ports.AuthRepository) *AuthService {
	return &AuthService{
		repository: AuthRepo,
	}
}
