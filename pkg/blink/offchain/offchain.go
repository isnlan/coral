package offchain

import (
	"context"
	"fmt"

	"github.com/isnlan/coral/pkg/errors"
	"github.com/isnlan/coral/pkg/response"
	"github.com/isnlan/coral/pkg/trace"

	"github.com/isnlan/coral/pkg/protos"

	"github.com/isnlan/coral/pkg/entity"
)

type QueryService interface {
	QueryChannelInfo(ctx context.Context, chainID, channelName string) (*protos.ChannelInformation, error)
	QueryBlocks(ctx context.Context, chainID, channelName string, query interface{}) ([]*entity.Block, error)
	QueryTxs(ctx context.Context, chainID, channelName string, query interface{}) ([]*entity.Transaction, error)
}

type (
	QueryChannelInfoRequest struct {
		NetworkID   string `json:"network_id" validate:"required"`
		ChannelName string `json:"channel_name" validate:"required"`
	}

	QueryChannelDataRequest struct {
		NetworkID   string      `json:"network_id" validate:"required"`
		ChannelName string      `json:"channel_name" validate:"required"`
		Query       interface{} `json:"query" validate:"required"`
	}

	QueryTxsRequest struct {
		NetworkID   string `json:"network_id" validate:"required"`
		ChannelName string `json:"channel_name" validate:"required"`
		BlockNumber uint64 `json:"block_number" validate:"required"`
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

func (c *client) QueryChannelInfo(ctx context.Context, chainID, channelName string) (*protos.ChannelInformation, error) {
	var (
		resp response.Response
		info = &protos.ChannelInformation{}
	)

	resp.Data = info
	req := &QueryChannelInfoRequest{
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

	return info, nil
}

func (c *client) QueryBlocks(ctx context.Context, chainID, channelName string, query interface{}) ([]*entity.Block, error) {
	var (
		resp response.Response
		qr   = &[]*entity.Block{}
	)

	resp.Data = qr

	req := &QueryChannelDataRequest{
		NetworkID:   chainID,
		ChannelName: channelName,
		Query:       query,
	}

	err := trace.DoPost(ctx, fmt.Sprintf("%s/api/v2/blockList", c.url), req, &resp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.ErrorCode != response.SuccessCode {
		return nil, errors.New(resp.Description)
	}

	return *qr, nil
}

func (c *client) QueryTxs(ctx context.Context, chainID, channelName string, query interface{}) ([]*entity.Transaction, error) {
	var (
		resp response.Response
		qr   = &[]*entity.Transaction{}
	)

	resp.Data = qr

	req := &QueryChannelDataRequest{
		NetworkID:   chainID,
		ChannelName: channelName,
		Query:       query,
	}

	err := trace.DoPost(ctx, fmt.Sprintf("%s/api/v2/txList", c.url), req, &resp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.ErrorCode != response.SuccessCode {
		return nil, errors.New(resp.Description)
	}

	return *qr, nil
}
