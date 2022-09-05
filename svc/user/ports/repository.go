package ports

import "auth-poc/svc/user/application/entity"

// Detail implementation can be seen on
// each usecase's DAO
type UserRepository interface {
	DeleteUser(id, phoneNumber string) error
	GetAllUser(limit, page int64) ([]*entity.User, error)
	GetUserByKey(userID, phoneNumber string) (*entity.User, error)
	InsertUniqueUser(usr *entity.User) (*entity.User, error)
	StorePassword(rawPassword string, phoneNumber string) error
	UpsertUser(usr *entity.User) (*entity.User, error)
}
