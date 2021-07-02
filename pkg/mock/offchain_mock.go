package mock

import (
	"context"
	"time"

	"github.com/isnlan/coral/pkg/blink/offchain"
	"github.com/isnlan/coral/pkg/entity"
)

var _ offchain.QueryService = &MockOffchain{}

type MockOffchain struct {
}

func (m *MockOffchain) QueryBlocks(ctx context.Context, chainID, channelName string, query map[string]interface{}, page, limit int) ([]*entity.Block, int64, error) {
	return nil, 0, nil
}

func (m *MockOffchain) QueryTxs(ctx context.Context, chainID, channelName string, query map[string]interface{}, page, limit int) ([]*entity.Transaction, int64, error) {
	return nil, 0, nil
}

func (m *MockOffchain) QueryChannelInfo(ctx context.Context, chainID, channelName string) (*entity.CheckPoint, error) {
	return &entity.CheckPoint{
		NetworkId:     chainID,
		ChainId:       channelName,
		SyncedBlock:   1,
		SyncedTotalTx: 1,
		SyncedTime:    time.Now(),
	}, nil
}
