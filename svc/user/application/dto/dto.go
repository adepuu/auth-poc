package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type RegisterRequest struct {
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Password    string `json:"password,omitempty"`
}

type UpdateRequest struct {
	UserID   string `json:"-"`
	UserType uint32 `json:"user_type"`
	RegisterRequest
}

type UserMongoDoc struct {
	DocumentID  primitive.ObjectID `bson:"_id,omitempty"`
	Email       string             `bson:"email"`
	FullName    string             `bson:"fullname"`
	PhoneNumber string             `bson:"phone_number"`
	UserType    uint32             `bson:"user_type"`
}
