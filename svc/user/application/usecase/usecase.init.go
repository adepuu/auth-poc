package usecase

import (
	"auth-poc/svc/user/application/service"
)

type UserUseCase struct {
	services *service.UserService
}

func NewUserUseCase(UserSvc *service.UserService) *UserUseCase {
	return &UserUseCase{
		services: UserSvc,
	}
}
