package xgrpc

import (
	"context"
	"time"

	"github.com/isnlan/coral/pkg/errors"
	"google.golang.org/grpc/status"

	grpc_middeware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/isnlan/coral/pkg/trace"

	"google.golang.org/grpc"
)

func NewClient(url string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			grpc_middeware.ChainUnaryClient(
				trace.OpenTracingClientInterceptor(),
			),
		),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return grpc.DialContext(ctx, url, opts...)
}

func Unwrap(err error) error {
	if err == nil {
		return err
	}
	if se, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return errors.New(se.GRPCStatus().Message())
	}

	return err
}
