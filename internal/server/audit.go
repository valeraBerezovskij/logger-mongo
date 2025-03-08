package server

import (
	"context"
	audit "github.com/valeraBerezovskij/logger-mongo/pkg/domain"
)

type AuditService interface {
	Insert(ctx context.Context, req *audit.LogItem) error
}

type AuditServer struct {
	service AuditService
}

func NewAuditServer(service AuditService) *AuditServer {
	return &AuditServer{service: service}
}

func (s *AuditServer) Insert(ctx context.Context, req *audit.LogItem) error {
	err := s.service.Insert(ctx, req)

	return err
}
