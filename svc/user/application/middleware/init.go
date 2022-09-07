package middleware

import (
	"auth-poc/svc/user/adapter/grpc"
	"auth-poc/svc/user/adapter/repository/mongo"
)

// Auth middleware checks if the request contains Token
// on the headers and if it is valid.
type Interactors struct {
	RpcClients *grpc.Clients
	DataStore  *mongo.DB
}

func NewMiddleware(i *Interactors) *Interactors {
	return i
}
