package offchain

import (
	"context"
	"fmt"

	"github.com/isnlan/coral/pkg/errors"
	"github.com/isnlan/coral/pkg/response"
	"github.com/isnlan/coral/pkg/trace"

	"github.com/isnlan/coral/pkg/entity"
)

type QueryService interface {
	QueryChannelInfo(ctx context.Context, chainID, channelName string) (*entity.CheckPoint, error)
	QueryBlocks(ctx context.Context, chainID, channelName string, query map[string]interface{}, page, limit int) ([]*entity.Block, int64, error)
	QueryTxs(ctx context.Context, chainID, channelName string, query map[string]interface{}, page, limit int) ([]*entity.Transaction, int64, error)
}

type (
	RequestChannelInfo struct {
		NetworkID   string `json:"network_id" validate:"required"`
		ChannelName string `json:"channel_name" validate:"required"`
	}

	RequestChannelData struct {
		NetworkID   string      `json:"network_id" validate:"required"`
		ChannelName string      `json:"channel_name" validate:"required"`
		Query       interface{} `json:"query" validate:"required"`
		Page        int         `json:"page" validate:"omitempty"`
		Limit       int         `json:"limit" validate:"omitempty"`
	}

	ResponseBlockList struct {
		Blocks []*entity.Block `json:"blocks"`
		Count  int64           `json:"count"`
	}

	ResponseTransactionList struct {
		Txs   []*entity.Transaction `json:"txs"`
		Count int64                 `json:"count"`
	}
)

type client struct {
	url string
}

func New(url string) QueryService {
	return &client{
		url: url,
	}
}

func (c *client) QueryChannelInfo(ctx context.Context, chainID, channelName string) (*entity.CheckPoint, error) {
	var (
		resp response.Response
		cp   = &entity.CheckPoint{}
	)

	resp.Data = cp
	req := &RequestChannelInfo{
		NetworkID:   chainID,
		ChannelName: channelName,
	}

	err := trace.DoPost(ctx, fmt.Sprintf("%s/api/v2/channelInfo", c.url), req, &resp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.ErrorCode != response.SuccessCode {
		return nil, errors.New(resp.Description)
	}

	return cp, nil
}

func (c *client) QueryBlocks(ctx context.Context, chainID, channelName string, query map[string]interface{}, page, limit int) ([]*entity.Block, int64, error) {
	var (
		resp response.Response
		qr   = new(ResponseBlockList)
	)

	resp.Data = qr

	req := &RequestChannelData{
		NetworkID:   chainID,
		ChannelName: channelName,
		Query:       query,
		Page:        page,
		Limit:       limit,
	}

	err := trace.DoPost(ctx, fmt.Sprintf("%s/api/v2/blockList", c.url), req, &resp)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	if resp.ErrorCode != response.SuccessCode {
		return nil, 0, errors.New(resp.Description)
	}

	return qr.Blocks, qr.Count, nil
}

func (c *client) QueryTxs(ctx context.Context, chainID, channelName string, query map[string]interface{}, page, limit int) ([]*entity.Transaction, int64, error) {
	var (
		resp response.Response
		qr   = new(ResponseTransactionList)
	)

	resp.Data = qr

	req := &RequestChannelData{
		NetworkID:   chainID,
		ChannelName: channelName,
		Query:       query,
		Page:        page,
		Limit:       limit,
	}

	err := trace.DoPost(ctx, fmt.Sprintf("%s/api/v2/txList", c.url), req, &resp)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	if resp.ErrorCode != response.SuccessCode {
		return nil, 0, errors.New(resp.Description)
	}

	return qr.Txs, qr.Count, nil
}
