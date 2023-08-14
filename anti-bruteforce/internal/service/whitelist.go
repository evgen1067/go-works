package service

import (
	"context"

	"github.com/evgen1067/anti-bruteforce/internal/repository"
)

type WhitelistService struct {
	ctx context.Context
	db  repository.ListRepo
}

func NewWhitelistService(ctx context.Context, db repository.ListRepo) *WhitelistService {
	return &WhitelistService{
		ctx: ctx,
		db:  db,
	}
}

func (w *WhitelistService) AddToWhitelist(address string) error {
	return w.db.AddToWhitelist(w.ctx, address)
}

func (w *WhitelistService) ExistsInWhitelist(address string) (bool, error) {
	return w.db.ExistsInWhitelist(w.ctx, address)
}

func (w *WhitelistService) DeleteFromWhitelist(address string) error {
	return w.db.DeleteFromWhitelist(w.ctx, address)
}
