package grpc

import (
	"auth-poc/svc/auth/adapter/grpc/pb"
	"net"

	log "github.com/angelbirth/logger"

	"auth-poc/svc/auth/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	User     pb.UserClient
	ListConn []*grpc.ClientConn
}

func New(config config.Config) (*Clients, error) {
	log.Infof("[GRPC][Client] Initializing clients conn with rpc servers")
	UserConn, e := grpc.Dial(
		net.JoinHostPort(config.RpcDefaultHost, config.RpcUserService),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithBlock(),
	)
	if e != nil {
		log.Error(e)
		return nil, e
	}

	log.Infof("[GRPC][Client] Client conn with rpc servers completed!")
	return &Clients{
		ListConn: []*grpc.ClientConn{UserConn},
		User:     pb.NewUserClient(UserConn),
	}, nil
}

func (c *Clients) Close() error {
	for _, conn := range c.ListConn {
		err := conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
