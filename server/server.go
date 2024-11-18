package server

import (
	"context"
	"log"
	"net"
	handler "tumbnail/internal/handlers"
	"tumbnail/pkg/api"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	handlers   *handler.Handler
}

func NewServer(handlers *handler.Handler) *Server {
	return &Server{
		handlers: handlers,
	}
}

func (s *Server) Run(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	s.listener = lis

	s.grpcServer = grpc.NewServer()
	api.RegisterThumbnailServer(s.grpcServer, s.handlers)

	log.Printf("gRPC server is running on port %s", port)
	return s.grpcServer.Serve(lis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return s.listener.Close()
}
