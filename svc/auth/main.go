package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	systemlog "log"

	"github.com/nightlyone/lockfile"

	"auth-poc/svc/auth/adapter/grpc"
	http "auth-poc/svc/auth/adapter/http"
	"auth-poc/svc/auth/adapter/repository/mongo"
	"auth-poc/svc/auth/application/dao"
	"auth-poc/svc/auth/application/service"
	"auth-poc/svc/auth/application/usecase"
	"auth-poc/svc/auth/config"
	"auth-poc/svc/auth/constants"

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

	authDao := dao.NewAuthDao(&dao.Auth{
		DataStore: dataStore,
		Rpc:       rpcClients,
	})

	authSvc := service.NewAuthService(authDao)
	authUseCase := usecase.NewAuthUseCase(authSvc)

	rpcServer := grpc.NewGrpcServer(&grpc.UseCases{
		Auth: *authUseCase,
	})

	httpHandler := http.NewHttpHandler(&http.Options{
		Auth: *authUseCase,
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
