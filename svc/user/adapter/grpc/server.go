package grpc

import (
	"net"

	log "github.com/angelbirth/logger"

	"auth-poc/svc/user/adapter/grpc/pb"
	handlers "auth-poc/svc/user/application/handlers/grpc"
	"auth-poc/svc/user/config"

	"auth-poc/svc/user/application/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UseCases struct {
	User usecase.UserUseCase
}

type Server struct {
	UseCases *UseCases
	Instance *grpc.Server
}

func NewGrpcServer(uc *UseCases) *Server {
	grpcServer := grpc.NewServer()
	return &Server{
		Instance: grpcServer,
		UseCases: uc,
	}
}

func (s *Server) Run(config config.Config) error {
	var err error
	var address = net.JoinHostPort(config.RpcDefaultHost, config.RpcPortUserApp)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("[GRPC][Server] failed to listen on %s: %v", address, err)
		return err
	}

	authHandlers := handlers.UserRpcHandler{UserUseCase: s.UseCases.User}

	pb.RegisterUserServer(s.Instance, &authHandlers)
	reflection.Register(s.Instance)
	log.Infof("[GRPC][Server] GRPC listening on %s\n", address)
	if err := s.Instance.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRPC over %s: %v", address, err)
		return err
	}
	return nil
}

func (s *Server) Close() {
	s.Instance.GracefulStop()
}
