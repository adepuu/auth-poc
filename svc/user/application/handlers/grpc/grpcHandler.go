package handlers

import (
	"context"
	"errors"

	log "github.com/angelbirth/logger"

	"auth-poc/svc/user/adapter/grpc/pb"
	"auth-poc/svc/user/application/usecase"
)

type UserRpcHandler struct {
	pb.UnimplementedUserServer
	UserUseCase usecase.UserUseCase
}

func (s *UserRpcHandler) GetUserByPhoneNumber(ctx context.Context, req *pb.GetUserByPhoneNumberArgs) (*pb.UserData, error) {
	phoneNumber := req.GetPhoneNumber()
	if len(phoneNumber) == 0 {
		return nil, errors.New("[User][Handler][GetUserByPhoneNumber] phone number can't be empty")
	}

	getResult, err := s.UserUseCase.GetUserByKey("", phoneNumber)
	if err != nil {
		log.Error("[gRPC][GetUserByPhoneNumber] Error: ", err)
		return nil, err
	}
	return &pb.UserData{
		Email:       getResult.Email,
		FullName:    getResult.FullName,
		PhoneNumber: getResult.PhoneNumber,
		UserID:      getResult.UserID,
		UserType:    getResult.UserType,
	}, nil
}
