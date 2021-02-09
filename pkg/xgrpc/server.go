package xgrpc

import (
	"log"
	"net"

	"github.com/isnlan/coral/pkg/errors"
	"github.com/isnlan/coral/pkg/trace"

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

	opts = append(opts, grpc.UnaryInterceptor(trace.OpenTracingServerInterceptor()))

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
