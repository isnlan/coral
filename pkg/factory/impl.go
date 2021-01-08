package factory

import (
	"context"

	"github.com/snlansky/coral/pkg/protos"
)

type network struct {
	chain   *protos.Chain
	channel string
	cli     protos.NetworkClient
	closer  func()
}

func (n *network) BuildChain(ctx context.Context) error {
	defer n.closer()
	_, err := n.cli.BuildChain(ctx, n.chain)
	return err
}

func (n *network) BuildChannel(ctx context.Context) error {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	_, err := n.cli.BuildChannel(ctx, c)
	return err
}

func (n *network) StartChain(ctx context.Context) error {
	defer n.closer()
	_, err := n.cli.StartChain(ctx, n.chain)
	return err
}

func (n *network) IsRunning(ctx context.Context) (bool, error) {
	defer n.closer()
	status, err := n.cli.IsRunning(ctx, n.chain)
	if err != nil {
		return false, err
	}

	return status.Status, nil
}

func (n *network) StopChain(ctx context.Context) error {
	defer n.closer()
	_, err := n.cli.StopChain(ctx, n.chain)
	return err
}

func (n *network) IsStopped(ctx context.Context) (bool, error) {
	defer n.closer()
	status, err := n.cli.IsStopped(ctx, n.chain)
	if err != nil {
		return false, err
	}

	return status.Status, nil
}

func (n *network) DeleteChain(ctx context.Context) error {
	defer n.closer()
	_, err := n.cli.DeleteChain(ctx, n.chain)
	return err
}

func (n *network) DownloadArtifacts(ctx context.Context) ([]byte, error) {
	defer n.closer()
	art, err := n.cli.DownloadArtifacts(ctx, n.chain)
	if err != nil {
		return nil, err
	}

	return art.Data, nil
}

func (n *network) Register(ctx context.Context, user string, pwd string) (*protos.DigitalIdentity, error) {
	defer n.closer()
	req := &protos.RequestRegister{
		Chain: n.chain,
		User:  user,
		Pwd:   pwd,
	}
	return n.cli.Register(ctx, req)
}

func (n *network) InstallContract(ctx context.Context, contract *protos.Contract) (string, error) {
	defer n.closer()
	req := &protos.RequestSetupContract{
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
	}
	resp, err := n.cli.InstallContract(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.TxId, nil
}

func (n *network) UpdateContract(ctx context.Context, contract *protos.Contract) (string, error) {
	defer n.closer()
	req := &protos.RequestSetupContract{
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
	}
	resp, err := n.cli.UpdateContract(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.TxId, nil
}

func (n *network) QueryContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) ([]byte, error) {
	defer n.closer()
	req := &protos.RequestQueryOrInvokeContract{
		Identity: identity,
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
		Args:     arg,
	}
	resp, err := n.cli.QueryContract(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (n *network) InvokeContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) (string, []byte, error) {
	defer n.closer()
	req := &protos.RequestQueryOrInvokeContract{
		Identity: identity,
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
		Args:     arg,
	}
	resp, err := n.cli.InvokeContract(ctx, req)
	if err != nil {
		return "", nil, err
	}
	return resp.TxId, resp.Data, nil
}

func (n *network) QueryChainNodes(ctx context.Context) ([]*protos.Node, error) {
	defer n.closer()
	nodes, err := n.cli.QueryChainNodes(ctx, n.chain)
	if err != nil {
		return nil, err
	}

	return nodes.Nodes, nil
}

func (n *network) QueryChannelList(ctx context.Context) ([]string, error) {
	defer n.closer()
	list, err := n.cli.QueryChannelList(ctx, n.chain)
	if err != nil {
		return nil, err
	}

	return list.Channels, nil
}

func (n *network) QueryChannel(ctx context.Context) (*protos.ChannelInformation, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	return n.cli.QueryChannel(ctx, c)
}

func (n *network) QueryContractList(ctx context.Context) ([]*protos.Contract, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	list, err := n.cli.QueryContractList(ctx, c)
	if err != nil {
		return nil, err
	}

	return list.Contracts, nil
}

func (n *network) QueryLatestBlock(ctx context.Context) (*protos.Block, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	block, err := n.cli.QueryLatestBlock(ctx, c)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (n *network) QueryBlockByNum(ctx context.Context, unm uint64) (*protos.Block, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryBlockByNum{
		Channel: c,
		Num:     unm,
	}
	block, err := n.cli.QueryBlockByNum(ctx, req)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (n *network) QueryBlockByTxId(ctx context.Context, txId string) (*protos.Block, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryBlockTxId{
		Channel: c,
		TxId:    txId,
	}
	block, err := n.cli.QueryBlockByTxId(ctx, req)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (n *network) QueryBlockByHash(ctx context.Context, hash []byte) (*protos.Block, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryBlockByHash{
		Channel: c,
		Hash:    hash,
	}
	block, err := n.cli.QueryBlockByHash(ctx, req)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (n *network) QueryTxById(ctx context.Context, txId string) (*protos.Transaction, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryTxById{
		Channel: c,
		TxId:    txId,
	}
	return n.cli.QueryTxById(ctx, req)
}
