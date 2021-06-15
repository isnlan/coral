package network

import (
	"fmt"
	"sync"

	"github.com/isnlan/coral/pkg/utils"

	"github.com/isnlan/coral/pkg/logging"

	"google.golang.org/grpc"

	"github.com/isnlan/coral/pkg/errors"
	_ "github.com/isnlan/coral/pkg/grpc-consul-resolver"
	"github.com/isnlan/coral/pkg/protos"
	"github.com/isnlan/coral/pkg/xgrpc"
)

const maxCallRecvMsgSize = 1024 * 1024 * 20

var logger = logging.MustGetLogger("network")

type Factory interface {
	Builder(chain *protos.Chain) (Builder, error)
	Close()
}

type consulBaseFactoryImpl struct {
	mu      *sync.RWMutex
	url     string
	clients map[string]*grpc.ClientConn
	opts    []grpc.DialOption
}

func New(url string) Factory {
	return &consulBaseFactoryImpl{
		mu:      new(sync.RWMutex),
		url:     url,
		clients: map[string]*grpc.ClientConn{},
		opts: []grpc.DialOption{grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize)),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`)},
	}
}

func (f *consulBaseFactoryImpl) Builder(chain *protos.Chain) (Builder, error) {
	client, err := f.getClient(chain.NetworkType)
	if err != nil {
		return nil, errors.WithMessage(err, "get network error")
	}

	return &builderImpl{chain: chain, client: client}, nil
}

func (f *consulBaseFactoryImpl) getClient(netType string) (*grpc.ClientConn, error) {
	f.mu.RLock()
	client, find := f.clients[netType]
	f.mu.RUnlock()

	if find {
		return client, nil
	}

	url := f.makeConsulUrl(netType)
	client, err := xgrpc.NewClient(url, f.opts...)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("create gprc client %s error", url))
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.clients[netType] = client
	logger.Infof("find service: %v", netType)
	return client, nil
}

func (f *consulBaseFactoryImpl) makeConsulUrl(netType string) string {
	var svr *protos.NetworkServer
	return fmt.Sprintf("consul://%s/%s?wait=30m&tag=%s&healthy=true&require-consistent=true",
		f.url, utils.MakeTypeName(svr), netType)
}

func (f *consulBaseFactoryImpl) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, client := range f.clients {
		_ = client.Close()
	}
}
