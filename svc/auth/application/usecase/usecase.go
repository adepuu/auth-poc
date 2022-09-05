package usecase

import (
	"auth-poc/svc/auth/application/dto"
	"auth-poc/svc/auth/application/entity"
	"errors"
)

func (u *AuthUseCase) CheckAuth(req *dto.CheckAuthRequest) (*dto.CheckAuthResponse, error) {
	if len(req.Token) == 0 {
		return nil, errors.New("[Auth][UseCase][CheckAuth] token can't be empty")
	}
	tokenClaims, valid, err := u.services.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}

	if !valid {
		return &dto.CheckAuthResponse{
			IsAuthorized: false,
		}, nil
	}

	return &dto.CheckAuthResponse{
		IsAuthorized: true,
		UserID:       tokenClaims.UserID,
		UserType:     uint(tokenClaims.UserType),
	}, nil
}

func (u *AuthUseCase) Login(req *dto.LoginRequest) (*entity.TokenPair, error) {
	if len(req.PhoneNumber) < 5 || len(req.Password) < 5 {
		return nil, errors.New("[Auth][UseCase][Login] bad format for phoneNumber or password")
	}
	match, err := u.services.CompareHash(req.PhoneNumber, req.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, nil
	}

	newTokenPair, err := u.services.GenerateToken(req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return newTokenPair, nil
}

func (u *AuthUseCase) RefreshToken(req string) (*entity.TokenPair, error) {
	newTokenPair, err := u.services.RefreshToken(req)
	if err != nil {
		return nil, err
	}

	return newTokenPair, nil
}

func (u *AuthUseCase) StorePassword(rawPassword, phoneNumber string) error {
	return u.services.StorePassword(rawPassword, phoneNumber)
}

func (u *AuthUseCase) RemovePassword(phoneNumber string) error {
	return u.services.RemovePassword(phoneNumber)
}
