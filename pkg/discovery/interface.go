package discovery

import (
	"context"
	"errors"
	"net"

	"github.com/hashicorp/consul/api"

	"google.golang.org/grpc"
)

const HTTPHealthCheckRouter = "/health"

type ServiceInfo struct {
	ID      string
	Address string
	Tags    []string
}

type Deregister func()

type ServiceDiscover interface {
	RegisterHealthServer(s *grpc.Server)
	ServiceRegister(name, address string, port int, tags ...string) (Deregister, error)
	HTTPServiceRegister(name, address string, port int, tags ...string) (Deregister, error)
	WatchService(ctx context.Context, name string, tag string, ch chan<- []*ServiceInfo)
	SetKey(ns, key string, value []byte) error
	GetKey(ns, key string) ([]byte, error)
	DeleteKey(ns, key string) error
	DeleteKeyByPrefix(ns, prefix string) error
	WatchKey(ctx context.Context, ns, key string, ch chan<- *api.KVPair)
	WatchKeysByPrefix(ctx context.Context, ns, prefix string, ch chan<- []string)
	WatchValuesByKeyPrefix(ctx context.Context, ns, prefix string, ch chan<- []*api.KVPair)
}

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("can't find a ip addr")
}
