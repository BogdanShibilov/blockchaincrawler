package v1

import (
	"fmt"
	"net"

	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/blockinfo/gw"
	"google.golang.org/grpc"
)

type Server struct {
	port       string
	service    *Service
	grpcServer *grpc.Server
}

func NewServer(port string, grpcService *Service) *Server {
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	return &Server{
		port:       port,
		service:    grpcService,
		grpcServer: grpcServer,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("failed to listen grpc port: %s", s.port)
	}

	pb.RegisterBlockInfoServiceServer(s.grpcServer, s.service)

	go s.grpcServer.Serve(listener)

	return nil
}

func (s *Server) Close() {
	s.grpcServer.GracefulStop()
}
