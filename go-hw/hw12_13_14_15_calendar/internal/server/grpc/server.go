package grpc

import (
	"github.com/evgen1067/hw12_13_14_15_calendar/api"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/services"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	address  string
	services *services.Services
	api.UnimplementedEventServiceServer
	*grpc.Server
}

func NewGRPC(services *services.Services, cfg *config.Config) *Server {
	srv := grpc.NewServer()
	server := &Server{
		address:  net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port),
		services: services,
		Server:   srv,
	}
	api.RegisterEventServiceServer(srv, server)
	return server
}

func (s *Server) ListenAndServe() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() {
	s.GracefulStop()
}
