package entity

import jwt "github.com/golang-jwt/jwt/v4"

type UserData struct {
	PhoneNumber string
	UserID      string
	UserType    uint
}

type TokenClaims struct {
	UserType    uint   `json:"user_type"`
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
