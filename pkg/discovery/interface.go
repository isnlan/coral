package discovery

import (
	"context"
	"errors"
	"net"
	"reflect"

	"google.golang.org/grpc"
)

type ServiceInfo struct {
	ID      string
	Address string
	Tags    []string
}

type Deregister func()

type ServiceDiscover interface {
	RegisterHealthServer(s *grpc.Server)
	ServiceRegister(name, address string, port int, tags ...string) (Deregister, error)
	WatchService(ctx context.Context, name string, tag string, ch chan<- []*ServiceInfo)
}

func MakeTypeName(tpy interface{}) string {
	return reflect.TypeOf(tpy).Elem().Name()
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
