package server

import (
	"fmt"
	audit "github.com/valeraBerezovskij/logger-mongo/pkg/domain"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	grpcServer  *grpc.Server
	auditServer audit.AuditServiceServer
}

func New(auditServer audit.AuditServiceServer) *Server {
	return &Server{
		grpcServer:  grpc.NewServer(),
		auditServer: auditServer,
	}
}

func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	//Server init
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	//Server register
	audit.RegisterAuditServiceServer(s.grpcServer, s.auditServer)

	//Server start
	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
