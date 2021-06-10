package mock

import (
	"context"

	"github.com/hashicorp/consul/api"
	"github.com/isnlan/coral/pkg/discovery"
	"google.golang.org/grpc"
)

type ServiceDiscoverMock struct {
}

func (s2 *ServiceDiscoverMock) RegisterHealthServer(s *grpc.Server) {
}

func (s2 *ServiceDiscoverMock) ServiceRegister(name, address string, port int, tags ...string) (discovery.Deregister, error) {
	return func() {}, nil
}

func (s2 *ServiceDiscoverMock) WatchService(ctx context.Context, name string, tag string, ch chan<- []*discovery.ServiceInfo) {
}

func (s2 *ServiceDiscoverMock) SetKey(ns, key string, value []byte) error {
	return nil
}

func (s2 *ServiceDiscoverMock) GetKey(ns, key string) ([]byte, error) {
	return []byte{}, nil
}

func (s2 *ServiceDiscoverMock) DeleteKey(ns, key string) error {
	return nil
}

func (s2 *ServiceDiscoverMock) DeleteKeyByPrefix(ns, prefix string) error {
	return nil
}

func (s2 *ServiceDiscoverMock) WatchKey(ctx context.Context, ns, key string, ch chan<- *api.KVPair) {

}

func (s2 *ServiceDiscoverMock) WatchKeysByPrefix(ctx context.Context, ns, prefix string, ch chan<- []string) {
}

func (s2 *ServiceDiscoverMock) WatchValuesByKeyPrefix(ctx context.Context, ns, prefix string, ch chan<- []*api.KVPair) {
}
