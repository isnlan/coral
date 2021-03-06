package network

import (
	"context"

	"github.com/isnlan/coral/pkg/xgrpc"

	"github.com/isnlan/coral/pkg/protos"
)

//go:generate mockgen -source=./pkg/network/network.go -destination=./pkg/mock/network_mock.go -package=mock
type Network interface {
	BuildChain(ctx context.Context) error
	BuildChannel(ctx context.Context) error
	StartChain(ctx context.Context) error
	IsRunning(ctx context.Context) (bool, error)
	StopChain(ctx context.Context) error
	IsStopped(ctx context.Context) (bool, error)
	DeleteChain(ctx context.Context) error
	EnableSyncDB(ctx context.Context) (string, error)
	DisableSyncDB(ctx context.Context) error
	DownloadArtifacts(ctx context.Context) ([]byte, error)
	Register(ctx context.Context, user string, pwd string) (*protos.DigitalIdentity, error)
	InstallContract(ctx context.Context, contract *protos.Contract) (string, error)
	UpdateContract(ctx context.Context, contract *protos.Contract) (string, error)
	QueryContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) ([]byte, error)
	InvokeContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) (string, []byte, error)
	QueryChainNodes(ctx context.Context) ([]*protos.Node, error)
	QueryChannelList(ctx context.Context) ([]string, error)
	QueryChannel(ctx context.Context) (*protos.ChannelInformation, error)
	QueryContractList(ctx context.Context) ([]string, error)
	QueryLatestBlock(ctx context.Context) (*protos.InnerBlock, error)
	QueryBlockByNum(ctx context.Context, unm uint64) (*protos.InnerBlock, error)
	QueryBlockByTxId(ctx context.Context, txId string) (*protos.InnerBlock, error)
	QueryBlockByHash(ctx context.Context, hash []byte) (*protos.InnerBlock, error)
	QueryTxById(ctx context.Context, txId string) (*protos.InnerTransaction, error)
}

type networkImpl struct {
	chain   *protos.Chain
	channel string
	client  protos.NetworkClient
}

func newNetworkImpl(chain *protos.Chain, channel string, client protos.NetworkClient) *networkImpl {
	return &networkImpl{
		chain:   chain,
		channel: channel,
		client:  client,
	}
}

func (n *networkImpl) BuildChain(ctx context.Context) error {
	_, err := n.client.BuildChain(ctx, n.chain)
	return xgrpc.Unwrap(err)
}

func (n *networkImpl) BuildChannel(ctx context.Context) error {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	_, err := n.client.BuildChannel(ctx, c)
	return xgrpc.Unwrap(err)
}

func (n *networkImpl) StartChain(ctx context.Context) error {
	_, err := n.client.StartChain(ctx, n.chain)
	return xgrpc.Unwrap(err)
}

func (n *networkImpl) IsRunning(ctx context.Context) (bool, error) {
	status, err := n.client.IsRunning(ctx, n.chain)
	if err != nil {
		return false, xgrpc.Unwrap(err)
	}

	return status.Status, nil
}

func (n *networkImpl) StopChain(ctx context.Context) error {
	_, err := n.client.StopChain(ctx, n.chain)
	return xgrpc.Unwrap(err)
}

func (n *networkImpl) IsStopped(ctx context.Context) (bool, error) {
	status, err := n.client.IsStopped(ctx, n.chain)
	if err != nil {
		return false, xgrpc.Unwrap(err)
	}

	return status.Status, nil
}

func (n *networkImpl) DeleteChain(ctx context.Context) error {
	_, err := n.client.DeleteChain(ctx, n.chain)
	return xgrpc.Unwrap(err)
}

func (n *networkImpl) EnableSyncDB(ctx context.Context) (string, error) {
	uri, err := n.client.EnableSyncDB(ctx, n.chain)
	if err != nil {
		return "", xgrpc.Unwrap(err)

	}
	return uri.Uri, nil
}

func (n *networkImpl) DisableSyncDB(ctx context.Context) error {
	_, err := n.client.DisableSyncDB(ctx, n.chain)
	return xgrpc.Unwrap(err)
}

func (n *networkImpl) DownloadArtifacts(ctx context.Context) ([]byte, error) {
	art, err := n.client.DownloadArtifacts(ctx, n.chain)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}

	return art.Data, nil
}

func (n *networkImpl) Register(ctx context.Context, user string, pwd string) (*protos.DigitalIdentity, error) {
	req := &protos.RequestRegister{
		Chain: n.chain,
		User:  user,
		Pwd:   pwd,
	}
	dig, err := n.client.Register(ctx, req)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}
	return dig, nil
}

func (n *networkImpl) InstallContract(ctx context.Context, contract *protos.Contract) (string, error) {
	req := &protos.RequestSetupContract{
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
	}
	resp, err := n.client.InstallContract(ctx, req)
	if err != nil {
		return "", xgrpc.Unwrap(err)
	}
	return resp.TxId, nil
}

func (n *networkImpl) UpdateContract(ctx context.Context, contract *protos.Contract) (string, error) {
	req := &protos.RequestSetupContract{
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
	}
	resp, err := n.client.UpdateContract(ctx, req)
	if err != nil {
		return "", xgrpc.Unwrap(err)
	}
	return resp.TxId, nil
}

func (n *networkImpl) QueryContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) ([]byte, error) {
	req := &protos.RequestQueryOrInvokeContract{
		Identity: identity,
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
		Args:     arg,
	}
	resp, err := n.client.QueryContract(ctx, req)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}
	return resp.Data, nil
}

func (n *networkImpl) InvokeContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) (string, []byte, error) {
	req := &protos.RequestQueryOrInvokeContract{
		Identity: identity,
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
		Args:     arg,
	}
	resp, err := n.client.InvokeContract(ctx, req)
	if err != nil {
		return "", nil, xgrpc.Unwrap(err)
	}
	return resp.TxId, resp.Data, nil
}

func (n *networkImpl) QueryChainNodes(ctx context.Context) ([]*protos.Node, error) {
	nodes, err := n.client.QueryChainNodes(ctx, n.chain)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}

	return nodes.Nodes, nil
}

func (n *networkImpl) QueryChannelList(ctx context.Context) ([]string, error) {
	list, err := n.client.QueryChannelList(ctx, n.chain)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}

	return list.Channels, nil
}

func (n *networkImpl) QueryChannel(ctx context.Context) (*protos.ChannelInformation, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	info, err := n.client.QueryChannel(ctx, c)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}

	return info, nil
}

func (n *networkImpl) QueryContractList(ctx context.Context) ([]string, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	list, err := n.client.QueryContractList(ctx, c)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}

	return list.Contracts, nil
}

func (n *networkImpl) QueryLatestBlock(ctx context.Context) (*protos.InnerBlock, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	block, err := n.client.QueryLatestBlock(ctx, c)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}
	return block, nil
}

func (n *networkImpl) QueryBlockByNum(ctx context.Context, unm uint64) (*protos.InnerBlock, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryBlockByNum{
		Channel: c,
		Num:     unm,
	}
	block, err := n.client.QueryBlockByNum(ctx, req)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}
	return block, nil
}

func (n *networkImpl) QueryBlockByTxId(ctx context.Context, txId string) (*protos.InnerBlock, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryBlockTxId{
		Channel: c,
		TxId:    txId,
	}
	block, err := n.client.QueryBlockByTxId(ctx, req)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}
	return block, nil
}

func (n *networkImpl) QueryBlockByHash(ctx context.Context, hash []byte) (*protos.InnerBlock, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryBlockByHash{
		Channel: c,
		Hash:    hash,
	}
	block, err := n.client.QueryBlockByHash(ctx, req)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}
	return block, nil
}

func (n *networkImpl) QueryTxById(ctx context.Context, txId string) (*protos.InnerTransaction, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryTxById{
		Channel: c,
		TxId:    txId,
	}
	tx, err := n.client.QueryTxById(ctx, req)
	if err != nil {
		return nil, xgrpc.Unwrap(err)
	}
	return tx, nil
}
