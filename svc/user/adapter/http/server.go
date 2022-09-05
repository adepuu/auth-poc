package ginAdapter

import (
	"context"

	log "github.com/angelbirth/logger"

	"net"
	"net/http"

	handlers "auth-poc/svc/user/application/handlers/http"
	middleware "auth-poc/svc/user/application/middleware/http"
	"auth-poc/svc/user/application/usecase"
	"auth-poc/svc/user/config"

	"github.com/gin-gonic/gin"
)

// Register another usecases below
type Options struct {
	User           usecase.UserUseCase
	AuthMiddleware *middleware.Auth
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
	h := handlers.UserHttpHandler{UserUseCase: s.Options.User}
	regularRoutes := s.Router.Group("/regular").Use(s.Options.AuthMiddleware.ValidateToken)
	internalRoutes := s.Router.Group("/internal")
	internalRoutes.Use(s.Options.AuthMiddleware.ValidateToken)
	internalRoutes.Use(s.Options.AuthMiddleware.AdminOnly)

	// Public Routes
	s.Router.POST("/register", h.Register)

	//Regular -> registered regular user Routes
	regularRoutes.GET("/profile", h.Profile)

	// Internal -> registered admin user Routes
	internalRoutes.GET("/detail/:id", h.Detail)
	internalRoutes.GET("/all-user", h.AllUser)
	internalRoutes.DELETE("/delete/:id", h.Delete)
	internalRoutes.PUT("/update/:id", h.Update)
}

func (s *Server) Shutdown() error {
	return s.Server.Shutdown(context.Background())
}
