package dao

import (
	"auth-poc/svc/auth/adapter/grpc"
	"auth-poc/svc/auth/adapter/repository/mongo"
)

type Auth struct {
	DataStore *mongo.DB
	Rpc       *grpc.Clients
}

func NewAuthDao(a *Auth) *Auth {
	return a
}
