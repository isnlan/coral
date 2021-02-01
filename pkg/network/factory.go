package network

import (
	"fmt"
	"sync"

	"github.com/snlansky/coral/pkg/discovery"

	"github.com/snlansky/coral/pkg/logging"

	"google.golang.org/grpc"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/snlansky/coral/pkg/errors"
	"github.com/snlansky/coral/pkg/protos"
	"github.com/snlansky/coral/pkg/xgrpc"
)

const maxCallRecvMsgSize = 1024 * 1024 * 20

var logger = logging.MustGetLogger("network")

type Factory struct {
	mu      *sync.RWMutex
	url     string
	clients map[string]*grpc.ClientConn
	opts    []grpc.DialOption
}

func New(url string) *Factory {
	return &Factory{
		mu:      new(sync.RWMutex),
		url:     url,
		clients: map[string]*grpc.ClientConn{},
		opts: []grpc.DialOption{grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize)),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`)},
	}
}

func (f *Factory) Builder(chain *protos.Chain) (*Builder, error) {
	client, err := f.getClient(chain.NetworkType)
	if err != nil {
		return nil, errors.WithMessage(err, "get network error")
	}

	return &Builder{chain: chain, client: client}, nil
}

func (f *Factory) getClient(netType string) (*grpc.ClientConn, error) {
	f.mu.RLock()
	client, find := f.clients[netType]
	f.mu.RUnlock()

	if find {
		return client, nil
	}

	url := f.makeConsulResolver(netType)
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

func (f *Factory) makeConsulResolver(netType string) string {
	var svr *protos.NetworkServer
	return fmt.Sprintf("consul://%s/%s?wait=30s&tag=%s&healthy=true&require-consistent=true",
		f.url, discovery.MakeTypeName(svr), netType)
}

func (f *Factory) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, client := range f.clients {
		_ = client.Close()
	}
}

type Builder struct {
	chain   *protos.Chain
	channel string
	client  *grpc.ClientConn
}

func (b *Builder) SetChannel(channel string) *Builder {
	b.channel = channel
	return b
}

func (b *Builder) Build() Network {
	return newNetworkImpl(b.chain, b.channel, protos.NewNetworkClient(b.client))
}
