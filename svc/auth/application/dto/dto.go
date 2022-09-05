package dto

type CheckAuthRequest struct {
	Token string
}

type LoginRequest struct {
	PhoneNumber string
	Password    string
}

type CheckAuthResponse struct {
	IsAuthorized bool   `json:"is_authorized"`
	UserID       string `json:"user_id"`
	UserType     uint   `json:"user_type"`
}

type TokenCollection struct {
	PhoneNumber string `bson:"phone_number"`
	Password    string `bosn:"password"`
}
