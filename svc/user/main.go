package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	systemlog "log"

	"github.com/nightlyone/lockfile"

	"auth-poc/svc/user/adapter/grpc"
	http "auth-poc/svc/user/adapter/http"
	"auth-poc/svc/user/adapter/repository/mongo"
	"auth-poc/svc/user/application/dao"
	httpMiddleware "auth-poc/svc/user/application/middleware/http"
	"auth-poc/svc/user/application/service"
	"auth-poc/svc/user/application/usecase"
	"auth-poc/svc/user/config"
	"auth-poc/svc/user/constants"

	log "github.com/angelbirth/logger"
	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

var (
	pidFile   string
	serverEnv string
)

func init() {
	flag.StringVarP(&pidFile, "pid-file", "p", "pid", "pid file")
	flag.StringVarP(&serverEnv, "env", "e", constants.ENV_LOCAL, "sever ENV")
	log.Init("", false, false, os.Stdout)
	log.SetFlags(systemlog.LstdFlags | systemlog.Lshortfile)
}

func main() {
	flag.Parse()
	config, err := config.LoadConfig("./config", serverEnv)
	if err != nil {
		log.Fatal("cannot load config :", err)
	}

	dataStore, err := mongo.New(config)
	if err != nil {
		log.Fatal(err)
	}

	rpcClients, err := grpc.New(config)
	if err != nil {
		log.Fatal(err)
	}

	authDao := dao.NewUserDao(&dao.User{
		DataStore: dataStore,
		Rpc:       rpcClients,
	})

	authSvc := service.NewUserService(authDao)
	authUseCase := usecase.NewUserUseCase(authSvc)

	rpcServer := grpc.NewGrpcServer(&grpc.UseCases{
		User: *authUseCase,
	})

	authMiddleware := httpMiddleware.NewAuthMiddleware(rpcClients)

	httpHandler := http.NewHttpHandler(&http.Options{
		User:           *authUseCase,
		AuthMiddleware: authMiddleware,
	})

	httpHandler.SetupRoutes()

	pidFile, _ = filepath.Abs(pidFile)
	lockFile, e := lockfile.New(pidFile)
	if e != nil {
		log.Fatal(e)
		return
	}
	e = lockFile.TryLock()
	if e != nil {
		log.Fatal(e)
		return
	}
	defer lockFile.Unlock()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-quitCh
		log.Infof("Caught signal %s, gracefully exiting...\n", sig)
		httpHandler.Shutdown()
		rpcServer.Close()
		rpcClients.Close()
		dataStore.Close()
	}()

	var eg errgroup.Group

	eg.Go(func() error {
		return httpHandler.Run(config)
	})

	eg.Go(func() error {
		return rpcServer.Run(config)
	})
	_ = eg.Wait()
}
