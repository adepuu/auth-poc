package constants

import "time"

const (
	DEFAULT_TIMEOUT = 3 * time.Second
	DB_NAME         = "account"
	DB_COLL         = "creds"

	ENV_DEVELOPMENT = "development"
	ENV_LOCAL       = "local"
	ENV_PRODUCTION  = "production"

	ACCESS_TOKEN_SIGN_KEY  = "37cd7f3d-a7c6-4b96-b94b-0b4fd982b7c5"
	REFRESH_TOKEN_SIGN_KEY = "cb6b1936-b349-4fe5-925d-64139bd69af4"

	ACCESS_TOKEN_EXPIRATION  = time.Minute * 5
	REFRESH_TOKEN_EXPIRATION = time.Minute * 60 * 2
)
