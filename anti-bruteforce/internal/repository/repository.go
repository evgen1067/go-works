package repository

import (
	"context"
)

type DatabaseRepo interface {
	Connect(ctx context.Context) error
	Close() error
	ListRepo
}

type ListRepo interface {
	AddToBlacklist(ctx context.Context, address string) error
	AddToWhitelist(ctx context.Context, address string) error
	ExistsInBlacklist(ctx context.Context, address string) (bool, error)
	ExistsInWhitelist(ctx context.Context, address string) (bool, error)
	DeleteFromBlacklist(ctx context.Context, address string) error
	DeleteFromWhitelist(ctx context.Context, address string) error
}
