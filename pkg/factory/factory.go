package factory

import (
	"fmt"
	"sync"

	"github.com/snlansky/coral/pkg/discovery"

	"github.com/snlansky/coral/pkg/logging"

	"google.golang.org/grpc"

	_ "github.com/mbobakov/grpc-consul-resolver"
	grpcpool "github.com/processout/grpc-go-pool"
	"github.com/snlansky/coral/pkg/errors"
	"github.com/snlansky/coral/pkg/net"
	"github.com/snlansky/coral/pkg/protos"
)

const maxCallRecvMsgSize = 20 * 1024 * 1024

var logger = logging.MustGetLogger("factory")

type Factory struct {
	lock sync.RWMutex
	url  string
	nets map[string]*net.Client
	opts []grpc.DialOption
}

func New(url string) *Factory {
	return &Factory{
		lock: sync.RWMutex{},
		url:  url,
		nets: map[string]*net.Client{},
		opts: []grpc.DialOption{grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize)),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`)},
	}
}

func (mgr *Factory) Builder(chain *protos.Chain) (*Builder, error) {
	client, err := mgr.getNetwork(chain.NetworkType)
	if err != nil {
		return nil, errors.WithMessage(err, "get network error")
	}

	conn, err := client.Get()
	if err != nil {
		return nil, err
	}

	return &Builder{chain: chain, conn: conn}, nil
}

func (mgr *Factory) getNetwork(netType string) (*net.Client, error) {
	mgr.lock.RLock()
	client, find := mgr.nets[netType]
	mgr.lock.RUnlock()

	if find {
		return client, nil
	}

	var svr *protos.NetworkServer

	client, err := net.New(fmt.Sprintf("%s/%s?wait=3s&tag=%s", mgr.url, discovery.MakeTypeName(svr), netType), mgr.opts...)
	if err != nil {
		return nil, err
	}

	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	mgr.nets[netType] = client
	logger.Infof("find service: %v", netType)
	return client, nil
}

func (mgr *Factory) Close() {
}

type Builder struct {
	chain   *protos.Chain
	conn    *grpcpool.ClientConn
	channel string
}

func (b *Builder) SetChannel(channel string) *Builder {
	b.channel = channel
	return b
}

func (b *Builder) Build() INetwork {
	return &network{
		chain:   b.chain,
		channel: b.channel,
		cli:     protos.NewNetworkClient(b.conn.ClientConn),
		closer: func() {
			_ = b.conn.Close()
		},
	}
}
