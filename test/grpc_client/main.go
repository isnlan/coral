package main

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"

	grpc_middeware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/snlansky/coral/pkg/errors"
	"github.com/snlansky/coral/pkg/trace"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {

	_, closer, err := trace.NewTracer("scpkg-grpc-client", "127.0.0.1:6831")
	errors.Check(err)
	defer closer.Close()

	ctx := context.Background()
	var conn *grpc.ClientConn

	serviceAddress := "127.0.0.1:8091"

	conn, err = grpc.DialContext(
		ctx,
		serviceAddress,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			grpc_middeware.ChainUnaryClient(
				trace.OpenTracingClientInterceptor(),
			),
		),
	)
	errors.Check(err)

	spanA := opentracing.StartSpan("A7")
	time.Sleep(time.Second)
	spanA.Finish()

	client := pb.NewGreeterClient(conn)

	_, err = client.SayHello(opentracing.ContextWithSpan(context.Background(), spanA), &pb.HelloRequest{
		Name: "lucy",
	})
	errors.Check(err)

	time.Sleep(time.Second * 2)
	_, err = client.SayHello(opentracing.ContextWithSpan(context.Background(), spanA), &pb.HelloRequest{
		Name: "lilith",
	})
	errors.Check(err)
}
