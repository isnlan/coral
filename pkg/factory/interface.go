package factory

import "github.com/snlansky/coral/pkg/protos"

type INetwork interface {
	BuildChain() error
	BuildChannel() error
	StartChain() error
	IsRunning() (bool, error)
	StopChain() error
	IsStopped() (bool, error)
	DeleteChain() error
	DownloadArtifacts() ([]byte, error)
	Register(user string, pwd string) (*protos.DigitalIdentity, error)
	InstallContract(contract *protos.Contract) (string, error)
	UpdateContract(contract *protos.Contract) (string, error)
	QueryContract(identity *protos.DigitalIdentity, contract string, arg []string) ([]byte, error)
	InvokeContract(identity *protos.DigitalIdentity, contract string, arg []string) (string, []byte, error)
	QueryChainNodes() ([]*protos.Node, error)
	QueryChannelList() ([]string, error)
	QueryChannel() (*protos.ChannelInformation, error)
	QueryContractList() ([]*protos.Contract, error)
	QueryLatestBlock() (*protos.Block, []*protos.Transaction, error)
	QueryBlockByNum(unm uint64) (*protos.Block, []*protos.Transaction, error)
	QueryBlockByTxId(txId string) (*protos.Block, []*protos.Transaction, error)
	QueryBlockByHash(hash []byte) (*protos.Block, []*protos.Transaction, error)
	QueryTxById(txId string) (*protos.Transaction, error)
}
