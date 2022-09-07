package middleware

import (
	"auth-poc/svc/auth/adapter/repository/mongo"
)

// Auth middleware checks if the request contains Token
// on the headers and if it is valid.
type Interactors struct {
	DataStore *mongo.DB
}

func NewMiddleware(i *Interactors) *Interactors {
	return i
}
