package constants

import "time"

const (
	BEARER_SCHEMA         = "Bearer"
	DB_COLL               = "users"
	DB_NAME               = "account"
	DEFAULT_TIMEOUT       = 3 * time.Second
	ENV_DEVELOPMENT       = "development"
	ENV_LOCAL             = "local"
	ENV_PRODUCTION        = "production"
	USER_ID_CTX           = "user-id"
	USER_TYPE_CTX         = "user-type"
	USER_PHONE_NUMBER_CTX = "phone-number"
)

type UserType int64

const (
	USER_TYPE_SUPER_ADMIN UserType = iota + 1
	USER_TYPE_REGULAR
)
