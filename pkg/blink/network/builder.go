package network

import (
	"github.com/isnlan/coral/pkg/protos"
	"google.golang.org/grpc"
)

type Builder interface {
	SetChannel(channel string) Builder
	Build() Network
}

type builderImpl struct {
	chain   *protos.Chain
	channel string
	client  *grpc.ClientConn
}

func (b *builderImpl) SetChannel(channel string) Builder {
	b.channel = channel
	return b
}

func (b *builderImpl) Build() Network {
	return newNetworkImpl(b.chain, b.channel, protos.NewNetworkClient(b.client))
}
