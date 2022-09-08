package usecase

import (
	"auth-poc/svc/user/application/dto"
	"auth-poc/svc/user/application/entity"
)

func (u *UserUseCase) GetUserByKey(userID, phoneNumber string) (*entity.User, error) {
	return u.services.GetUserByKey(userID, phoneNumber)
}

func (u *UserUseCase) DeleteUser(userID string) error {
	return u.services.DeleteUser(userID)
}

func (u *UserUseCase) GetAllUser(limit, page int64) ([]*entity.User, error) {
	return u.services.GetAllUser(limit, page)
}

func (u *UserUseCase) Register(req *dto.RegisterRequest) (*entity.User, error) {
	newUser, err := u.services.AddUser(&entity.User{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		FullName:    req.FullName,
		UserType:    req.UserType,
	})

	if err != nil {
		return nil, err
	}

	if newUser != nil {
		err = u.services.StorePassword(req.Password, req.PhoneNumber, nil)
		if err != nil {
			// Try rollback inserted data when calling Auth rpc store password failed
			rollbackErr := u.services.DeleteUser(newUser.UserID)
			if rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}
	}

	return newUser, nil
}

func (u *UserUseCase) Update(req *dto.UpdateRequest) (*entity.User, error) {
	updatedUser, err := u.services.UpdateUser(&entity.User{
		UserID:   req.UserID,
		Email:    req.Email,
		FullName: req.FullName,
		UserType: req.UserType,
	})

	if err != nil {
		return nil, err
	}

	if len(req.Password) > 0 {
		err = u.services.StorePassword(req.Password, req.PhoneNumber, &req.UserID)
		if err != nil {
			return nil, err
		}
	}

	return updatedUser, nil
}
