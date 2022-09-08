package constants

import "time"

const (
	AS_ADMIN_VALUE        = "iknowitsdangerous"
	BEARER_SCHEMA         = "Bearer"
	DB_COLL               = "users"
	DB_NAME               = "account"
	DEFAULT_TIMEOUT       = 3 * time.Second
	ENV_DEVELOPMENT       = "development"
	ENV_LOCAL             = "local"
	ENV_PRODUCTION        = "production"
	USER_ID_CTX           = "user-id"
	USER_PHONE_NUMBER_CTX = "phone-number"
	USER_TYPE_CTX         = "user-type"
)

type UserType int64

const (
	USER_TYPE_SUPER_ADMIN UserType = iota + 1
	USER_TYPE_REGULAR
)
