package service

import (
	"auth-poc/svc/auth/application/entity"
	"auth-poc/svc/auth/constants"
	"encoding/json"
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) ValidateToken(token string) (*entity.TokenClaims, bool, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.ACCESS_TOKEN_SIGN_KEY), nil
	})
	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return nil, false, nil
		}
		return nil, false, err
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false, errors.New("[Auth][Service][ValidateToken] failed getting token claims")
	}

	marshalled, err := json.Marshal(claims)
	if err != nil {
		return nil, false, err
	}

	var tokenClaims entity.TokenClaims
	if err := json.Unmarshal(marshalled, &tokenClaims); err != nil {
		return nil, false, err
	}

	return &tokenClaims, true, nil
}

func (s *AuthService) GenerateToken(phoneNumber string) (*entity.TokenPair, error) {
	userData, err := s.repository.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}

	signedAccessToken, err := signToken(userData, constants.ACCESS_TOKEN_EXPIRATION, constants.ACCESS_TOKEN_SIGN_KEY)
	if err != nil {
		return nil, err
	}

	signedRefreshToken, err := signToken(userData, constants.REFRESH_TOKEN_EXPIRATION, constants.REFRESH_TOKEN_SIGN_KEY)
	if err != nil {
		return nil, err
	}

	return &entity.TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*entity.TokenPair, error) {
	parsed, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.REFRESH_TOKEN_SIGN_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("[Auth][Service][ValidateToken] failed getting token claims")
	}

	marshalled, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	var tokenClaims entity.TokenClaims

	if err := json.Unmarshal(marshalled, &tokenClaims); err != nil {
		return nil, err
	}

	userData := &entity.UserData{
		UserID:      tokenClaims.UserID,
		UserType:    tokenClaims.UserType,
		PhoneNumber: tokenClaims.PhoneNumber,
	}

	signedAccessToken, err := signToken(userData, constants.ACCESS_TOKEN_EXPIRATION, constants.ACCESS_TOKEN_SIGN_KEY)
	if err != nil {
		return nil, err
	}

	signedRefreshToken, err := signToken(userData, constants.REFRESH_TOKEN_EXPIRATION, constants.REFRESH_TOKEN_SIGN_KEY)
	if err != nil {
		return nil, err
	}

	return &entity.TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (s *AuthService) CompareHash(phoneNumber, password string) (bool, error) {
	storedHash, err := s.repository.GetStoredHash(phoneNumber)
	if err != nil {
		return false, err
	}

	// user not exist
	if len(storedHash) == 0 {
		return false, nil
	}

	errCompare := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if errCompare != nil {
		if errCompare == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, errCompare
	}
	return true, nil
}

func (s *AuthService) StorePassword(raw string, phoneNumber string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repository.UpsertCredential(hashedPassword, phoneNumber)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RemovePassword(phoneNumber string) error {
	err := s.repository.RemoveCredential(phoneNumber)
	if err != nil {
		return err
	}
	return nil
}

func signToken(data *entity.UserData, expiration time.Duration, signKey string) (string, error) {
	newToken := jwt.New(jwt.SigningMethodHS256)
	newTokenClaims := newToken.Claims.(jwt.MapClaims)
	newTokenClaims["user_id"] = data.UserID
	newTokenClaims["phone_number"] = data.PhoneNumber
	newTokenClaims["user_type"] = data.UserType
	newTokenClaims["exp"] = jwt.NewNumericDate(time.Now().Add(expiration))

	signedToken, err := newToken.SignedString([]byte(signKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
