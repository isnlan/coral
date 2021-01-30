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

const maxCallRecvMsgSize = 20 * 1024 * 1024

var logger = logging.MustGetLogger("network")

type Factory struct {
	lock *sync.RWMutex
	url  string
	nets map[string]*grpc.ClientConn
	opts []grpc.DialOption
}

func New(url string) *Factory {
	return &Factory{
		lock: new(sync.RWMutex),
		url:  url,
		nets: map[string]*grpc.ClientConn{},
		opts: []grpc.DialOption{grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize)),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`)},
	}
}

func (f *Factory) Builder(chain *protos.Chain) (*Builder, error) {
	client, err := f.getNetwork(chain.NetworkType)
	if err != nil {
		return nil, errors.WithMessage(err, "get network error")
	}

	return &Builder{chain: chain, client: client}, nil
}

func (f *Factory) getNetwork(netType string) (*grpc.ClientConn, error) {
	f.lock.RLock()
	client, find := f.nets[netType]
	f.lock.RUnlock()

	if find {
		return client, nil
	}

	var svr *protos.NetworkServer

	client, err := xgrpc.NewClient(fmt.Sprintf("consul://%s/%s?wait=30s&tag=%s&healthy=true&require-consistent=true",
		f.url, discovery.MakeTypeName(svr), netType), f.opts...)
	if err != nil {
		return nil, err
	}

	f.lock.Lock()
	defer f.lock.Unlock()
	f.nets[netType] = client
	logger.Infof("find service: %v", netType)
	return client, nil
}

func (f *Factory) Close() {
	f.lock.Lock()
	defer f.lock.Unlock()
	for _, client := range f.nets {
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
	return &networkImpl{
		chain:   b.chain,
		channel: b.channel,
		client:  protos.NewNetworkClient(b.client),
	}
}
