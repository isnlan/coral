package main

import (
	"context"

	"github.com/isnlan/coral/pkg/errors"
	"github.com/isnlan/coral/pkg/trace"
	"github.com/isnlan/coral/pkg/xgrpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var serviceAddress = "0.0.0.0:8091"

func main() {
	_, closer, err := trace.NewTracer(serviceName, "127.0.0.1:6831")
	errors.Check(err)
	defer closer.Close()

	//lis, err := net.Listen("tcp", serviceAddress)
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//s := grpc.NewServer(grpc.UnaryInterceptor(trace.OpenTracingServerInterceptor()))
	//pb.RegisterGreeterServer(s, &server{})
	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}

	s, err := xgrpc.NewServer(serviceAddress)
	if err != nil {
		panic(err)
	}

	pb.RegisterGreeterServer(s.Server(), &server{})
	s.Start()
}

const (
	serviceName = "HelloServer"
)

type server struct {
}

func (s server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	if request.Name == "lucy" {
		//n := 0
		//_ = 1 / n
		return nil, errors.New("name error")
	}
	return &pb.HelloReply{
		Message: "hi: " + request.Name,
	}, nil
}

var _ pb.GreeterServer = &server{}
