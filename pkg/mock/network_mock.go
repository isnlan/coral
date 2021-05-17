// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/network/network.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	protos "github.com/isnlan/coral/pkg/protos"
	reflect "reflect"
)

// MockNetwork is a mock of Network interface
type MockNetwork struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkMockRecorder
}

// MockNetworkMockRecorder is the mock recorder for MockNetwork
type MockNetworkMockRecorder struct {
	mock *MockNetwork
}

// NewMockNetwork creates a new mock instance
func NewMockNetwork(ctrl *gomock.Controller) *MockNetwork {
	mock := &MockNetwork{ctrl: ctrl}
	mock.recorder = &MockNetworkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNetwork) EXPECT() *MockNetworkMockRecorder {
	return m.recorder
}

// BuildChain mocks base method
func (m *MockNetwork) BuildChain(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildChain", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuildChain indicates an expected call of BuildChain
func (mr *MockNetworkMockRecorder) BuildChain(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildChain", reflect.TypeOf((*MockNetwork)(nil).BuildChain), ctx)
}

// BuildChannel mocks base method
func (m *MockNetwork) BuildChannel(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildChannel", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuildChannel indicates an expected call of BuildChannel
func (mr *MockNetworkMockRecorder) BuildChannel(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildChannel", reflect.TypeOf((*MockNetwork)(nil).BuildChannel), ctx)
}

// StartChain mocks base method
func (m *MockNetwork) StartChain(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartChain", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartChain indicates an expected call of StartChain
func (mr *MockNetworkMockRecorder) StartChain(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartChain", reflect.TypeOf((*MockNetwork)(nil).StartChain), ctx)
}

// IsRunning mocks base method
func (m *MockNetwork) IsRunning(ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRunning", ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRunning indicates an expected call of IsRunning
func (mr *MockNetworkMockRecorder) IsRunning(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRunning", reflect.TypeOf((*MockNetwork)(nil).IsRunning), ctx)
}

// StopChain mocks base method
func (m *MockNetwork) StopChain(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopChain", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopChain indicates an expected call of StopChain
func (mr *MockNetworkMockRecorder) StopChain(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopChain", reflect.TypeOf((*MockNetwork)(nil).StopChain), ctx)
}

// IsStopped mocks base method
func (m *MockNetwork) IsStopped(ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsStopped", ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsStopped indicates an expected call of IsStopped
func (mr *MockNetworkMockRecorder) IsStopped(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsStopped", reflect.TypeOf((*MockNetwork)(nil).IsStopped), ctx)
}

// DeleteChain mocks base method
func (m *MockNetwork) DeleteChain(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChain", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteChain indicates an expected call of DeleteChain
func (mr *MockNetworkMockRecorder) DeleteChain(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChain", reflect.TypeOf((*MockNetwork)(nil).DeleteChain), ctx)
}

// DownloadArtifacts mocks base method
func (m *MockNetwork) DownloadArtifacts(ctx context.Context) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadArtifacts", ctx)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadArtifacts indicates an expected call of DownloadArtifacts
func (mr *MockNetworkMockRecorder) DownloadArtifacts(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadArtifacts", reflect.TypeOf((*MockNetwork)(nil).DownloadArtifacts), ctx)
}

// Register mocks base method
func (m *MockNetwork) Register(ctx context.Context, user, pwd string) (*protos.DigitalIdentity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, user, pwd)
	ret0, _ := ret[0].(*protos.DigitalIdentity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register
func (mr *MockNetworkMockRecorder) Register(ctx, user, pwd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockNetwork)(nil).Register), ctx, user, pwd)
}

// InstallContract mocks base method
func (m *MockNetwork) InstallContract(ctx context.Context, contract *protos.Contract) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstallContract", ctx, contract)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstallContract indicates an expected call of InstallContract
func (mr *MockNetworkMockRecorder) InstallContract(ctx, contract interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstallContract", reflect.TypeOf((*MockNetwork)(nil).InstallContract), ctx, contract)
}

// UpdateContract mocks base method
func (m *MockNetwork) UpdateContract(ctx context.Context, contract *protos.Contract) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContract", ctx, contract)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateContract indicates an expected call of UpdateContract
func (mr *MockNetworkMockRecorder) UpdateContract(ctx, contract interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContract", reflect.TypeOf((*MockNetwork)(nil).UpdateContract), ctx, contract)
}

// QueryContract mocks base method
func (m *MockNetwork) QueryContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryContract", ctx, identity, contract, arg)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContract indicates an expected call of QueryContract
func (mr *MockNetworkMockRecorder) QueryContract(ctx, identity, contract, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContract", reflect.TypeOf((*MockNetwork)(nil).QueryContract), ctx, identity, contract, arg)
}

// InvokeContract mocks base method
func (m *MockNetwork) InvokeContract(ctx context.Context, identity *protos.DigitalIdentity, contract string, arg []string) (string, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokeContract", ctx, identity, contract, arg)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// InvokeContract indicates an expected call of InvokeContract
func (mr *MockNetworkMockRecorder) InvokeContract(ctx, identity, contract, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokeContract", reflect.TypeOf((*MockNetwork)(nil).InvokeContract), ctx, identity, contract, arg)
}

// QueryChainNodes mocks base method
func (m *MockNetwork) QueryChainNodes(ctx context.Context) ([]*protos.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryChainNodes", ctx)
	ret0, _ := ret[0].([]*protos.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryChainNodes indicates an expected call of QueryChainNodes
func (mr *MockNetworkMockRecorder) QueryChainNodes(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryChainNodes", reflect.TypeOf((*MockNetwork)(nil).QueryChainNodes), ctx)
}

// QueryChannelList mocks base method
func (m *MockNetwork) QueryChannelList(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryChannelList", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryChannelList indicates an expected call of QueryChannelList
func (mr *MockNetworkMockRecorder) QueryChannelList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryChannelList", reflect.TypeOf((*MockNetwork)(nil).QueryChannelList), ctx)
}

// QueryChannel mocks base method
func (m *MockNetwork) QueryChannel(ctx context.Context) (*protos.ChannelInformation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryChannel", ctx)
	ret0, _ := ret[0].(*protos.ChannelInformation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryChannel indicates an expected call of QueryChannel
func (mr *MockNetworkMockRecorder) QueryChannel(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryChannel", reflect.TypeOf((*MockNetwork)(nil).QueryChannel), ctx)
}

// EnableSyncChannelDB mocks base method
func (m *MockNetwork) EnableSyncChannelDB(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableSyncChannelDB", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnableSyncChannelDB indicates an expected call of EnableSyncChannelDB
func (mr *MockNetworkMockRecorder) EnableSyncChannelDB(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableSyncChannelDB", reflect.TypeOf((*MockNetwork)(nil).EnableSyncChannelDB), ctx)
}

// DisableSyncChannelDB mocks base method
func (m *MockNetwork) DisableSyncChannelDB(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableSyncChannelDB", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DisableSyncChannelDB indicates an expected call of DisableSyncChannelDB
func (mr *MockNetworkMockRecorder) DisableSyncChannelDB(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableSyncChannelDB", reflect.TypeOf((*MockNetwork)(nil).DisableSyncChannelDB), ctx)
}

// QueryContractList mocks base method
func (m *MockNetwork) QueryContractList(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryContractList", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContractList indicates an expected call of QueryContractList
func (mr *MockNetworkMockRecorder) QueryContractList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContractList", reflect.TypeOf((*MockNetwork)(nil).QueryContractList), ctx)
}

// QueryLatestBlock mocks base method
func (m *MockNetwork) QueryLatestBlock(ctx context.Context) (*protos.InnerBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryLatestBlock", ctx)
	ret0, _ := ret[0].(*protos.InnerBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryLatestBlock indicates an expected call of QueryLatestBlock
func (mr *MockNetworkMockRecorder) QueryLatestBlock(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryLatestBlock", reflect.TypeOf((*MockNetwork)(nil).QueryLatestBlock), ctx)
}

// QueryBlockByNum mocks base method
func (m *MockNetwork) QueryBlockByNum(ctx context.Context, unm uint64) (*protos.InnerBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryBlockByNum", ctx, unm)
	ret0, _ := ret[0].(*protos.InnerBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryBlockByNum indicates an expected call of QueryBlockByNum
func (mr *MockNetworkMockRecorder) QueryBlockByNum(ctx, unm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryBlockByNum", reflect.TypeOf((*MockNetwork)(nil).QueryBlockByNum), ctx, unm)
}

// QueryBlockByTxId mocks base method
func (m *MockNetwork) QueryBlockByTxId(ctx context.Context, txId string) (*protos.InnerBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryBlockByTxId", ctx, txId)
	ret0, _ := ret[0].(*protos.InnerBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryBlockByTxId indicates an expected call of QueryBlockByTxId
func (mr *MockNetworkMockRecorder) QueryBlockByTxId(ctx, txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryBlockByTxId", reflect.TypeOf((*MockNetwork)(nil).QueryBlockByTxId), ctx, txId)
}

// QueryBlockByHash mocks base method
func (m *MockNetwork) QueryBlockByHash(ctx context.Context, hash []byte) (*protos.InnerBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryBlockByHash", ctx, hash)
	ret0, _ := ret[0].(*protos.InnerBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryBlockByHash indicates an expected call of QueryBlockByHash
func (mr *MockNetworkMockRecorder) QueryBlockByHash(ctx, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryBlockByHash", reflect.TypeOf((*MockNetwork)(nil).QueryBlockByHash), ctx, hash)
}

// QueryTxById mocks base method
func (m *MockNetwork) QueryTxById(ctx context.Context, txId string) (*protos.InnerTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryTxById", ctx, txId)
	ret0, _ := ret[0].(*protos.InnerTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryTxById indicates an expected call of QueryTxById
func (mr *MockNetworkMockRecorder) QueryTxById(ctx, txId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryTxById", reflect.TypeOf((*MockNetwork)(nil).QueryTxById), ctx, txId)
}
