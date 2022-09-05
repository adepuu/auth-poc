package service

import (
	"auth-poc/svc/user/application/entity"
	"auth-poc/svc/user/constants"
	"auth-poc/svc/user/helper/sanitize"

	"net/mail"
)

func (s *UserService) GetUserByKey(userID, phoneNumber string) (*entity.User, error) {
	return s.repository.GetUserByKey(userID, phoneNumber)
}

func (s *UserService) GetAllUser(limit, page int64) ([]*entity.User, error) {
	return s.repository.GetAllUser(limit, page)
}

func (s *UserService) AddUser(usr *entity.User) (*entity.User, error) {
	_, err := mail.ParseAddress(usr.Email)
	if err != nil {
		return nil, err
	}
	usr.FullName = sanitize.GeneralSanitize(usr.FullName)
	usr.PhoneNumber = sanitize.GeneralSanitize(usr.PhoneNumber)
	usr.Email = sanitize.GeneralSanitize(usr.Email)

	return s.repository.InsertUniqueUser(usr)
}

func (s *UserService) UpdateUser(usr *entity.User) (*entity.User, error) {
	_, err := mail.ParseAddress(usr.Email)
	if err != nil {
		return nil, err
	}
	usr.FullName = sanitize.GeneralSanitize(usr.FullName)
	usr.Email = sanitize.GeneralSanitize(usr.Email)

	if usr.UserType == 0 {
		usr.UserType = uint32(constants.USER_TYPE_REGULAR)
	}

	return s.repository.UpsertUser(usr)
}

func (s *UserService) StorePassword(rawPassword, phoneNumber string, userID *string) error {
	if len(phoneNumber) == 0 && userID != nil {
		userData, err := s.repository.GetUserByKey(*userID, "")
		if err != nil {
			return err
		}
		phoneNumber = userData.PhoneNumber
	}
	return s.repository.StorePassword(rawPassword, phoneNumber)
}

func (s *UserService) DeleteUser(id string) error {
	userData, err := s.repository.GetUserByKey(id, "")
	if err != nil {
		return err
	}
	return s.repository.DeleteUser(id, userData.PhoneNumber)
}
