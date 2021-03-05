package xgrpc

import (
	"log"
	"net"

	"github.com/isnlan/coral/pkg/errors"
	"github.com/isnlan/coral/pkg/trace"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type Server struct {
	addr     string
	listener net.Listener
	svr      *grpc.Server
}

func NewServer(addr string, opts ...grpc.ServerOption) (*Server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to listen")
	}

	streamInterceptor := grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recoverHandler)),
	))

	unaryInterceptor := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractorForInitialReq(requestFieldExtractor)),
		trace.OpenTracingServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		LoggerUnaryServerInterceptor(),
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recoverHandler)),
	))

	opts = append(opts, streamInterceptor, unaryInterceptor)

	s := &Server{
		addr:     addr,
		listener: listener,
		svr:      grpc.NewServer(opts...),
	}

	return s, nil
}

func (s *Server) Server() *grpc.Server {
	return s.svr
}

func (s *Server) Start() {
	log.Printf("Listening and serving GRPC on %s\n", s.addr)
	if err := s.svr.Serve(s.listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) Close() {
	s.svr.Stop()
}
