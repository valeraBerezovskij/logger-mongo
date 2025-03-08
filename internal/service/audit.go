package service

import (
	"context"
	audit "github.com/valeraBerezovskij/logger-mongo/pkg/domain"
)

type Repository interface {
	Insert(ctx context.Context, item audit.LogItem) error
}
type Audit struct {
	repo Repository
}

func NewAudit(repo Repository) *Audit {
	return &Audit{
		repo: repo,
	}
}

func (a *Audit) Insert(ctx context.Context, req *audit.LogItem) error {
	item := audit.LogItem{
		Action:    req.Action,
		Entity:    req.Entity,
		EntityID:  req.EntityID,
		Timestamp: req.Timestamp,
	}
	return a.repo.Insert(ctx, item)
}
