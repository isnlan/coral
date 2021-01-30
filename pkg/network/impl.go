package network

import (
	"context"

	"github.com/snlansky/coral/pkg/protos"
)

type network struct {
	chain   *protos.Chain
	channel string
	client  protos.NetworkClient
}

func (n *network) BuildChain(ctx context.Context) error {
	_, err := n.client.BuildChain(ctx, n.chain)
	return err
}

func (n *network) BuildChannel(ctx context.Context) error {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	_, err := n.client.BuildChannel(ctx, c)
	return err
}

func (n *network) StartChain(ctx context.Context) error {
	_, err := n.client.StartChain(ctx, n.chain)
	return err
}

func (n *network) IsRunning(ctx context.Context) (bool, error) {
	status, err := n.client.IsRunning(ctx, n.chain)
	if err != nil {
		return false, err
	}

	return status.Status, nil
}

func (n *network) StopChain(ctx context.Context) error {
	_, err := n.client.StopChain(ctx, n.chain)
	return err
}

func (n *network) IsStopped(ctx context.Context) (bool, error) {
	status, err := n.client.IsStopped(ctx, n.chain)
	if err != nil {
		return false, err
	}

	return status.Status, nil
}

func (n *network) DeleteChain(ctx context.Context) error {
	_, err := n.client.DeleteChain(ctx, n.chain)
	return err
}

func (n *network) DownloadArtifacts(ctx context.Context) ([]byte, error) {
	art, err := n.client.DownloadArtifacts(ctx, n.chain)
	if err != nil {
		return nil, err
	}

	return art.Data, nil
}

func (n *network) Register(ctx context.Context, user string, pwd string) (*protos.DigitalIdentity, error) {
	req := &protos.RequestRegister{
		Chain: n.chain,
		User:  user,
		Pwd:   pwd,
	}
	return n.client.Register(ctx, req)
}

func (n *network) InstallContract(ctx context.Context, contract *protos.Contract) (string, error) {
	req := &protos.RequestSetupContract{
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
	}
	resp, err := n.client.InstallContract(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.TxId, nil
}

func (n *network) UpdateContract(ctx context.Context, contract *protos.Contract) (string, error) {
	req := &protos.RequestSetupContract{
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
	}
	resp, err := n.client.UpdateContract(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.TxId, nil
}

func (n *network) QueryContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) ([]byte, error) {
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
		return nil, err
	}
	return resp.Data, nil
}

func (n *network) InvokeContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) (string, []byte, error) {
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
		return "", nil, err
	}
	return resp.TxId, resp.Data, nil
}

func (n *network) QueryChainNodes(ctx context.Context) ([]*protos.Node, error) {
	nodes, err := n.client.QueryChainNodes(ctx, n.chain)
	if err != nil {
		return nil, err
	}

	return nodes.Nodes, nil
}

func (n *network) QueryChannelList(ctx context.Context) ([]string, error) {
	list, err := n.client.QueryChannelList(ctx, n.chain)
	if err != nil {
		return nil, err
	}

	return list.Channels, nil
}

func (n *network) QueryChannel(ctx context.Context) (*protos.ChannelInformation, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	return n.client.QueryChannel(ctx, c)
}

func (n *network) QueryContractList(ctx context.Context) ([]string, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	list, err := n.client.QueryContractList(ctx, c)
	if err != nil {
		return nil, err
	}

	return list.Contracts, nil
}

func (n *network) QueryLatestBlock(ctx context.Context) (*protos.InnerBlock, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	block, err := n.client.QueryLatestBlock(ctx, c)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (n *network) QueryBlockByNum(ctx context.Context, unm uint64) (*protos.InnerBlock, error) {
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
		return nil, err
	}
	return block, nil
}

func (n *network) QueryBlockByTxId(ctx context.Context, txId string) (*protos.InnerBlock, error) {
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
		return nil, err
	}
	return block, nil
}

func (n *network) QueryBlockByHash(ctx context.Context, hash []byte) (*protos.InnerBlock, error) {
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
		return nil, err
	}
	return block, nil
}

func (n *network) QueryTxById(ctx context.Context, txId string) (*protos.InnerTransaction, error) {
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryTxById{
		Channel: c,
		TxId:    txId,
	}
	return n.client.QueryTxById(ctx, req)
}
