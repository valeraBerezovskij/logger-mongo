package server

import (
	"context"
	audit "github.com/valeraBerezovskij/logger-mongo/pkg/domain"
)

type AuditService interface {
	Insert(ctx context.Context, req *audit.LogRequest) error
}

type AuditServer struct {
	audit.UnimplementedAuditServiceServer
	service AuditService
}

func NewAuditServer(service AuditService) *AuditServer {
	return &AuditServer{service: service}
}

func (s *AuditServer) Log(ctx context.Context, req *audit.LogRequest) (*audit.Empty, error) {
	err := s.service.Insert(ctx, req)

	return &audit.Empty{}, err
}
