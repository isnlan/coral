package factory

import (
	"context"

	"github.com/snlansky/coral/pkg/protos"
)

type INetwork interface {
	BuildChain(ctx context.Context) error
	BuildChannel(ctx context.Context) error
	StartChain(ctx context.Context) error
	IsRunning(ctx context.Context) (bool, error)
	StopChain(ctx context.Context) error
	IsStopped(ctx context.Context) (bool, error)
	DeleteChain(ctx context.Context) error
	DownloadArtifacts(ctx context.Context) ([]byte, error)
	Register(ctx context.Context, user string, pwd string) (*protos.DigitalIdentity, error)
	InstallContract(ctx context.Context, contract *protos.Contract) (string, error)
	UpdateContract(ctx context.Context, contract *protos.Contract) (string, error)
	QueryContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) ([]byte, error)
	InvokeContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) (string, []byte, error)
	QueryChainNodes(ctx context.Context) ([]*protos.Node, error)
	QueryChannelList(ctx context.Context) ([]string, error)
	QueryChannel(ctx context.Context) (*protos.ChannelInformation, error)
	QueryContractList(ctx context.Context) ([]*protos.Contract, error)
	QueryLatestBlock(ctx context.Context) (*protos.Block, []*protos.Transaction, error)
	QueryBlockByNum(ctx context.Context, unm uint64) (*protos.Block, []*protos.Transaction, error)
	QueryBlockByTxId(ctx context.Context, txId string) (*protos.Block, []*protos.Transaction, error)
	QueryBlockByHash(ctx context.Context, hash []byte) (*protos.Block, []*protos.Transaction, error)
	QueryTxById(ctx context.Context, txId string) (*protos.Transaction, error)
}
