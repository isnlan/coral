package net

import (
	"log"
	"net"

	"github.com/snlansky/coral/pkg/errors"
	"github.com/snlansky/coral/pkg/trace"

	"google.golang.org/grpc"
)

type Server struct {
	addr     string
	listener net.Listener
	svr      *grpc.Server
}

func NewServer(addr string) (*Server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to listen")
	}

	var opts []grpc.ServerOption
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
