# AUTH-POC
A simple/minimal authentication and authorization proof of concept using JWT.

## Development
- Test Admin ID: `6285151500400`
- Test Admin Password : `sapananya2`

## How to set session
Main idea is to get `access token` for session authentication and use `refresh token` when `access token` expired.

### TTL (Time To Live)
- `access token` is *5 Minutes*
- `refresh token` is *60 Minutes*

### Step
- login using `Login` on AUTH folder with admin test account above to get `access_token` and `refresh token`
- Authentication header is automatically set everytime `Login` or `Refresh` endpoint called (see AUTH POC Collection -> Test)
- If `access token` is expired, try use `Refresh` endpoint to get new token

**postman collection login staging** : https://www.getpostman.com/collections/b93ea1f4d200eb85385d


### Run service
For Testing you use [this postman](https://github.com/adepuu/auth-poc/blob/main/files/documents/AuthPOC.postman_collection.json)
To test Grpc  [grpcox](https://github.com/gusaul/grpcox)

### Reset Data
- Run `make docker-rebuilddb` to reset data
