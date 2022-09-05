package ginAdapter

import (
	"context"

	log "github.com/angelbirth/logger"

	"net"
	"net/http"

	handlers "auth-poc/svc/auth/application/handlers/http"
	"auth-poc/svc/auth/application/usecase"
	"auth-poc/svc/auth/config"

	"github.com/gin-gonic/gin"
)

// Register another usecases below
type Options struct {
	Auth usecase.AuthUseCase
}
type Server struct {
	Options *Options
	Router  *gin.Engine
	Server  *http.Server
}

func NewHttpHandler(opts *Options) *Server {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.Default()
	return &Server{
		Options: opts,
		Router:  ginRouter,
	}
}

func (s *Server) Run(config config.Config) error {
	listener, err := net.Listen("tcp", net.JoinHostPort(config.HttpHost, config.HttpPort))
	if err != nil {
		return err
	}
	s.Server = &http.Server{Handler: s.Router}
	log.Infof("[HTTP Server] listening on %s\n", listener.Addr().String())
	s.Server.RegisterOnShutdown(func() {
		listener.Close()
	})
	return s.Server.Serve(listener)
}

func (s *Server) SetupRoutes() {
	h := handlers.AuthHttpHandler{AuthUseCase: s.Options.Auth}
	s.Router.POST("/login", h.Login)
	s.Router.GET("/refresh", h.RefreshToken)
}

func (s *Server) Shutdown() error {
	return s.Server.Shutdown(context.Background())
}
