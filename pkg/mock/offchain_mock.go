package mock

import (
	"context"
	"time"

	"github.com/isnlan/coral/pkg/blink/offchain"
	"github.com/isnlan/coral/pkg/entity"
	"github.com/isnlan/coral/pkg/protos"
)

var _ offchain.QueryService = &MockOffchain{}

type MockOffchain struct {
}

func (m *MockOffchain) QueryChannelInfo(ctx context.Context, chainID, channelName string) (*protos.ChannelInformation, error) {
	return &protos.ChannelInformation{
		ChannelId:         channelName,
		BlockNumber:       1,
		TotalTransactions: 1,
		StartTime:         uint64(time.Now().Unix()),
	}, nil
}

func (m *MockOffchain) QueryBlocks(ctx context.Context, chainID, channelName string, query interface{}) ([]*entity.Block, error) {
	return nil, nil
}

func (m *MockOffchain) QueryTxs(ctx context.Context, chainID, channelName string, query interface{}) ([]*entity.Transaction, error) {
	return nil, nil
}
