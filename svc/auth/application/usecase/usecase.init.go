package usecase

import (
	"auth-poc/svc/auth/application/service"
)

type AuthUseCase struct {
	services *service.AuthService
}

func NewAuthUseCase(AuthSvc *service.AuthService) *AuthUseCase {
	return &AuthUseCase{
		services: AuthSvc,
	}
}
