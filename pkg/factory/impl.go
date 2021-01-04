package factory

import (
	"context"
	"time"

	"github.com/snlansky/coral/pkg/protos"
)

type network struct {
	chain   *protos.Chain
	channel string
	cli     protos.NetworkClient
	closer  func()
}

func (n *network) getContext() context.Context {
	ctx := context.Background()
	timeout, _ := context.WithTimeout(ctx, time.Second*5)
	return timeout
}

func (n *network) BuildChain() error {
	defer n.closer()
	_, err := n.cli.BuildChain(n.getContext(), n.chain)
	return err
}

func (n *network) BuildChannel() error {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	_, err := n.cli.BuildChannel(n.getContext(), c)
	return err
}

func (n *network) StartChain() error {
	defer n.closer()
	_, err := n.cli.StartChain(n.getContext(), n.chain)
	return err
}

func (n *network) IsRunning() (bool, error) {
	defer n.closer()
	status, err := n.cli.IsRunning(n.getContext(), n.chain)
	if err != nil {
		return false, err
	}

	return status.Status, nil
}

func (n *network) StopChain() error {
	defer n.closer()
	_, err := n.cli.StopChain(n.getContext(), n.chain)
	return err
}

func (n *network) IsStopped() (bool, error) {
	defer n.closer()
	status, err := n.cli.IsStopped(n.getContext(), n.chain)
	if err != nil {
		return false, err
	}

	return status.Status, nil
}

func (n *network) DeleteChain() error {
	defer n.closer()
	_, err := n.cli.DeleteChain(n.getContext(), n.chain)
	return err
}

func (n *network) DownloadArtifacts() ([]byte, error) {
	defer n.closer()
	art, err := n.cli.DownloadArtifacts(n.getContext(), n.chain)
	if err != nil {
		return nil, err
	}

	return art.Data, nil
}

func (n *network) Register(user string, pwd string) (*protos.DigitalIdentity, error) {
	defer n.closer()
	req := &protos.RequestRegister{
		Chain: n.chain,
		User:  user,
		Pwd:   pwd,
	}
	return n.cli.Register(n.getContext(), req)
}

func (n *network) InstallContract(contract *protos.Contract) (string, error) {
	defer n.closer()
	req := &protos.RequestSetupContract{
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
	}
	resp, err := n.cli.InstallContract(n.getContext(), req)
	if err != nil {
		return "", err
	}
	return resp.TxId, nil
}

func (n *network) UpdateContract(contract *protos.Contract) (string, error) {
	defer n.closer()
	req := &protos.RequestSetupContract{
		Channel: &protos.Channel{
			Chain: n.chain,
			Name:  n.channel,
		},
		Contract: contract,
	}
	resp, err := n.cli.UpdateContract(n.getContext(), req)
	if err != nil {
		return "", err
	}
	return resp.TxId, nil
}

func (n *network) QueryContract(identity *protos.DigitalIdentity, contract string, arg []string) ([]byte, error) {
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
	resp, err := n.cli.QueryContract(n.getContext(), req)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (n *network) InvokeContract(identity *protos.DigitalIdentity, contract string, arg []string) (string, []byte, error) {
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
	resp, err := n.cli.InvokeContract(n.getContext(), req)
	if err != nil {
		return "", nil, err
	}
	return resp.TxId, resp.Data, nil
}

func (n *network) QueryChainNodes() ([]*protos.Node, error) {
	defer n.closer()
	nodes, err := n.cli.QueryChainNodes(n.getContext(), n.chain)
	if err != nil {
		return nil, err
	}

	return nodes.Nodes, nil
}

func (n *network) QueryChannelList() ([]string, error) {
	defer n.closer()
	list, err := n.cli.QueryChannelList(n.getContext(), n.chain)
	if err != nil {
		return nil, err
	}

	return list.Channels, nil
}

func (n *network) QueryChannel() (*protos.ChannelInfo, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	return n.cli.QueryChannel(n.getContext(), c)
}

func (n *network) QueryContractList() ([]*protos.Contract, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	list, err := n.cli.QueryContractList(n.getContext(), c)
	if err != nil {
		return nil, err
	}

	return list.Contracts, nil
}

func (n *network) QueryLatestBlock() (*protos.Block, []*protos.Transaction, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	fullBlock, err := n.cli.QueryLatestBlock(n.getContext(), c)
	if err != nil {
		return nil, nil, err
	}
	return fullBlock.Block, fullBlock.Txs, nil
}

func (n *network) QueryBlockByNum(unm uint64) (*protos.Block, []*protos.Transaction, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryBlockByNum{
		Channel: c,
		Num:     unm,
	}
	fullBlock, err := n.cli.QueryBlockByNum(n.getContext(), req)
	if err != nil {
		return nil, nil, err
	}
	return fullBlock.Block, fullBlock.Txs, nil
}

func (n *network) QueryBlockByTxId(txId string) (*protos.Block, []*protos.Transaction, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryBlockTxId{
		Channel: c,
		TxId:    txId,
	}
	fullBlock, err := n.cli.QueryBlockByTxId(n.getContext(), req)
	if err != nil {
		return nil, nil, err
	}
	return fullBlock.Block, fullBlock.Txs, nil
}

func (n *network) QueryBlockByHash(hash []byte) (*protos.Block, []*protos.Transaction, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryBlockByHash{
		Channel: c,
		Hash:    hash,
	}
	fullBlock, err := n.cli.QueryBlockByHash(n.getContext(), req)
	if err != nil {
		return nil, nil, err
	}
	return fullBlock.Block, fullBlock.Txs, nil
}

func (n *network) QueryTxById(txId string) (*protos.Transaction, error) {
	defer n.closer()
	c := &protos.Channel{
		Chain: n.chain,
		Name:  n.channel,
	}
	req := &protos.RequestQueryTxById{
		Channel: c,
		TxId:    txId,
	}
	return n.cli.QueryTxById(n.getContext(), req)
}
