package dao

import (
	"auth-poc/svc/user/adapter/grpc"
	"auth-poc/svc/user/adapter/repository/mongo"
)

type User struct {
	DataStore *mongo.DB
	Rpc       *grpc.Clients
}

func NewUserDao(a *User) *User {
	return a
}
