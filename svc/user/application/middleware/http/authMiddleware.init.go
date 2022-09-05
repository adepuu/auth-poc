package middleware

import (
	"auth-poc/svc/user/adapter/grpc"
)

// Auth middleware checks if the request contains Token
// on the headers and if it is valid.
type Auth struct {
	RpcClients *grpc.Clients
}

func NewAuthMiddleware(rpc *grpc.Clients) *Auth {
	return &Auth{
		RpcClients: rpc,
	}
}
