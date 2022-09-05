package ports

import "auth-poc/svc/auth/application/entity"

// Detail implementation can be seen on
// each usecase's DAO
type AuthRepository interface {
	GetStoredHash(phoneNumber string) (string, error)
	GetUserByPhoneNumber(phoneNumber string) (*entity.UserData, error)
	RemoveCredential(phoneNumber string) error
	UpsertCredential(byteArr []byte, phoneNumber string) error
}
