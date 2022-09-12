package grpc

import (
	"auth-poc/svc/user/adapter/grpc/pb"
	"net"

	log "github.com/angelbirth/logger"

	"auth-poc/svc/user/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Auth     pb.AuthClient
	ListConn []*grpc.ClientConn
}

func New(config config.Config) (*Clients, error) {
	log.Infof("[GRPC][Client] Initializing clients conn with rpc servers")
	AuthConn, e := grpc.Dial(
		net.JoinHostPort(config.RpcAuthHost, config.RpcAuthService),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithBlock(),
	)
	if e != nil {
		log.Error(e)
		return nil, e
	}

	log.Infof("[GRPC][Client] Client conn with rpc servers completed!")
	return &Clients{
		ListConn: []*grpc.ClientConn{AuthConn},
		Auth:     pb.NewAuthClient(AuthConn),
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
