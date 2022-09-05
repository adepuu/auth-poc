package service

import (
	"auth-poc/svc/user/ports"
)

type UserService struct {
	repository ports.UserRepository
}

func NewUserService(UserRepo ports.UserRepository) *UserService {
	return &UserService{
		repository: UserRepo,
	}
}
