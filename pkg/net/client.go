package net

import (
	"context"
	"time"

	grpc_middeware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/snlansky/coral/pkg/trace"

	"github.com/snlansky/coral/pkg/errors"

	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

type Client struct {
	url  string
	pool *grpcpool.Pool
}

func New(url string, opts ...grpc.DialOption) (*Client, error) {
	opts = append(opts, grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			grpc_middeware.ChainUnaryClient(
				trace.OpenTracingClientInterceptor(),
			),
		),
	)

	factory := func() (conn *grpc.ClientConn, err error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		conn, err = grpc.DialContext(ctx, url, opts...)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return conn, nil
	}

	pool, err := grpcpool.New(factory, 8, 16, time.Hour, time.Hour*24)
	if err != nil {
		return nil, err
	}
	m := &Client{
		url:  url,
		pool: pool,
	}
	return m, nil
}

func (c *Client) Get() (*grpcpool.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conn, err := c.pool.Get(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get grpc connection error")
	}
	return conn, err
}
