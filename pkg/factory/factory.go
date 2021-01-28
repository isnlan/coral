package factory

import (
	"context"
	"fmt"
	"sync"

	"github.com/snlansky/coral/pkg/logging"

	"github.com/snlansky/coral/pkg/service_discovery"

	"google.golang.org/grpc"

	grpcpool "github.com/processout/grpc-go-pool"
	"github.com/snlansky/coral/pkg/errors"
	"github.com/snlansky/coral/pkg/protos"

	"github.com/snlansky/coral/pkg/net"
)

const maxCallRecvMsgSize = 20 * 1024 * 1024

var logger = logging.MustGetLogger("factory")

type Factory struct {
	lock      sync.RWMutex
	discovery service_discovery.ServiceDiscover
	nets      map[string]*net.Client
	opts      []grpc.DialOption
	cancels   []context.CancelFunc
}

func New(discovery service_discovery.ServiceDiscover) *Factory {
	return &Factory{
		lock:      sync.RWMutex{},
		discovery: discovery,
		nets:      map[string]*net.Client{},
		opts:      []grpc.DialOption{grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize))},
	}
}

func (mgr *Factory) Register(networkType string, addr string) error {
	cli, err := net.New(addr, mgr.opts...)
	if err != nil {
		return err
	}

	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	mgr.nets[networkType] = cli
	return nil
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

	if mgr.discovery == nil {
		return nil, errors.New(fmt.Sprintf("network %s not register", netType))
	}

	var svr *protos.NetworkServer
	resolver := NewResolver()
	ctx, cancel := context.WithCancel(context.Background())
	mgr.discovery.WatchService(ctx, service_discovery.MakeTypeName(svr), netType, resolver)
	client, err := net.NewWithResolver(resolver, mgr.opts...)
	if err != nil {
		cancel()
		return nil, err
	}

	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	mgr.nets[netType] = client
	mgr.cancels = append(mgr.cancels, cancel)
	logger.Infof("find service: %v", netType)
	return client, nil
}

func (mgr *Factory) Close() {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	for _, c := range mgr.cancels {
		c()
	}
	mgr.cancels = nil
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
