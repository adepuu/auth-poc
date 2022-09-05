package handlers

import (
	"context"

	log "github.com/angelbirth/logger"

	"auth-poc/svc/auth/adapter/grpc/pb"
	"auth-poc/svc/auth/application/dto"
	"auth-poc/svc/auth/application/usecase"
)

type AuthRpcHandler struct {
	pb.UnimplementedAuthServer
	AuthUseCase usecase.AuthUseCase
}

func (s *AuthRpcHandler) CheckAuth(ctx context.Context, req *pb.CheckAuthArgs) (*pb.CheckAuthReply, error) {
	input := dto.CheckAuthRequest{
		Token: req.GetToken(),
	}
	authResult, err := s.AuthUseCase.CheckAuth(&input)
	if err != nil {
		log.Error("[gRPC][CheckAuthHandler] Error: ", err)
		return nil, err
	}

	return &pb.CheckAuthReply{
		IsAuthorized: authResult.IsAuthorized,
		UserID:       authResult.UserID,
		UserType:     uint32(authResult.UserType),
	}, nil
}

func (s *AuthRpcHandler) StorePassword(ctx context.Context, req *pb.StorePasswordRequest) (*pb.StateReply, error) {
	err := s.AuthUseCase.StorePassword(req.GetRawPassword(), req.GetPhoneNumber())
	if err != nil {
		log.Error("[gRPC][StorePassword] Error: ", err)
		return nil, err
	}
	return &pb.StateReply{
		Success: true,
	}, nil
}

func (s *AuthRpcHandler) RemovePassword(ctx context.Context, req *pb.RemovePasswordRequest) (*pb.StateReply, error) {
	err := s.AuthUseCase.RemovePassword(req.GetPhoneNumber())
	if err != nil {
		log.Error("[gRPC][RemovePassword] Error: ", err)
		return nil, err
	}
	return &pb.StateReply{
		Success: true,
	}, nil
}
