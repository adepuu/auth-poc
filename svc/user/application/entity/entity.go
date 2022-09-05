package entity

type User struct {
	Email       string `json:"email" bson:"email"`
	FullName    string `json:"full_name" bson:"fullname"`
	PhoneNumber string `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	UserID      string `json:"user_id" bson:"-"`
	UserType    uint32 `json:"user_type" bson:"user_type"`
}
