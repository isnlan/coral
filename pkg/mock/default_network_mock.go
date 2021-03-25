package mock

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/isnlan/coral/pkg/network"
	"github.com/isnlan/coral/pkg/protos"
)

var _ network.Network = &DefaultMockNetwork{}

type DefaultMockNetwork struct {
}

func (d *DefaultMockNetwork) BuildChain(ctx context.Context) error {
	return nil
}

func (d *DefaultMockNetwork) BuildChannel(ctx context.Context) error {
	return nil
}

func (d *DefaultMockNetwork) StartChain(ctx context.Context) error {
	return nil
}

func (d *DefaultMockNetwork) IsRunning(ctx context.Context) (bool, error) {
	return true, nil
}

func (d *DefaultMockNetwork) StopChain(ctx context.Context) error {
	return nil
}

func (d *DefaultMockNetwork) IsStopped(ctx context.Context) (bool, error) {
	return true, nil
}

func (d *DefaultMockNetwork) DeleteChain(ctx context.Context) error {
	return nil
}

func (d *DefaultMockNetwork) DownloadArtifacts(ctx context.Context) ([]byte, error) {
	return []byte("data"), nil
}

func (d *DefaultMockNetwork) Register(ctx context.Context, user string, pwd string) (*protos.DigitalIdentity, error) {
	return nil, nil
}

func (d *DefaultMockNetwork) InstallContract(ctx context.Context, contract *protos.Contract) (string, error) {
	return "tx1", nil
}

func (d *DefaultMockNetwork) UpdateContract(ctx context.Context, contract *protos.Contract) (string, error) {
	return "tx2", nil
}

func (d *DefaultMockNetwork) QueryContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) ([]byte, error) {
	return []byte("data"), nil
}

func (d *DefaultMockNetwork) InvokeContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) (string, []byte, error) {
	return "tx3", []byte("data"), nil
}

func (d *DefaultMockNetwork) QueryChainNodes(ctx context.Context) ([]*protos.Node, error) {
	return []*protos.Node{}, nil
}

func (d *DefaultMockNetwork) QueryChannelList(ctx context.Context) ([]string, error) {
	return []string{"mychannel"}, nil
}

func (d *DefaultMockNetwork) QueryChannel(ctx context.Context) (*protos.ChannelInformation, error) {
	return &protos.ChannelInformation{BlockNumber: 2}, nil
}

func (d *DefaultMockNetwork) QueryContractList(ctx context.Context) ([]string, error) {
	return []string{"kvdb"}, nil
}

func (d *DefaultMockNetwork) QueryLatestBlock(ctx context.Context) (*protos.InnerBlock, error) {
	return block, nil
}

func (d *DefaultMockNetwork) QueryBlockByNum(ctx context.Context, unm uint64) (*protos.InnerBlock, error) {
	return block, nil
}

func (d *DefaultMockNetwork) QueryBlockByTxId(ctx context.Context, txId string) (*protos.InnerBlock, error) {
	return block, nil
}

func (d *DefaultMockNetwork) QueryBlockByHash(ctx context.Context, hash []byte) (*protos.InnerBlock, error) {
	return block, nil
}

func (d *DefaultMockNetwork) QueryTxById(ctx context.Context, txId string) (*protos.InnerTransaction, error) {
	return tx, nil
}

var block = &protos.InnerBlock{
	Number:       2,
	PreviousHash: []byte("phash"),
	Hash:         []byte("hash"),
	DataHash:     []byte("dhash"),
	Transactions: []*protos.InnerTransaction{tx},
	ChannelId:    "mychannel",
	Size:         1024,
	Timestamp:    nil,
}

var tx = &protos.InnerTransaction{
	TxId:           "tx1",
	ChannelId:      "mychannel",
	BlockNumber:    2,
	Contract:       "kvdb",
	Creator:        "mycreator",
	Sign:           []byte("sign"),
	TxType:         "TX",
	Timestamp:      timestamppb.Now(),
	ValidationCode: "VALID",
	Event:          nil,
}
