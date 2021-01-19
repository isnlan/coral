package factory

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"

	grpcpool "github.com/processout/grpc-go-pool"
	"github.com/snlansky/coral/pkg/errors"
	"github.com/snlansky/coral/pkg/protos"

	"github.com/snlansky/coral/pkg/net"
)

type Factory struct {
	lock sync.RWMutex
	nets map[string]*net.Client
}

func New() *Factory {
	return &Factory{lock: sync.RWMutex{}, nets: map[string]*net.Client{}}
}

func (mgr *Factory) Register(networkType string, addr string) error {
	cli, err := net.New(addr, grpc.WithMaxMsgSize(20*1024*1024))
	if err != nil {
		return err
	}

	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	mgr.nets[networkType] = cli
	return nil
}

func (mgr *Factory) Builder(chain *protos.Chain) (*Builder, error) {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()

	client, find := mgr.nets[chain.NetworkType]
	if !find {
		return nil, errors.New(fmt.Sprintf("network %s not register", chain.NetworkType))
	}

	conn, err := client.Get()
	if err != nil {
		return nil, err
	}

	return &Builder{chain: chain, conn: conn}, nil
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
