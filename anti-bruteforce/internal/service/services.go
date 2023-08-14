package service

import (
	"context"

	"github.com/evgen1067/anti-bruteforce/internal/bucket"
	"github.com/evgen1067/anti-bruteforce/internal/common"
	"github.com/evgen1067/anti-bruteforce/internal/logger"
	"github.com/evgen1067/anti-bruteforce/internal/repository/psql"
	"go.uber.org/zap"
)

type Auth interface {
	Authorize(req common.APIAuthRequest) bool
	ResetBucket()
}

type Blacklist interface {
	AddToBlacklist(address string) error
	ExistsInBlacklist(address string) (bool, error)
	DeleteFromBlacklist(address string) error
}

type Whitelist interface {
	AddToWhitelist(address string) error
	ExistsInWhitelist(address string) (bool, error)
	DeleteFromWhitelist(address string) error
}

type LeakyBucket interface {
	Add(req common.APIAuthRequest) bool
	ResetBucket()
}

type Logger interface {
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
}

type Services struct {
	Auth
	Blacklist
	Whitelist
	Logger
}

func NewServices(ctx context.Context, db *psql.Repo, lb *bucket.LeakyBucket, l *logger.Logger) *Services {
	blacklist := NewBlacklistService(ctx, db)
	whitelist := NewWhitelistService(ctx, db)
	logg := NewLogger(l)
	auth := NewAuthService(ctx, blacklist, whitelist, lb)
	return &Services{
		Auth:      auth,
		Blacklist: blacklist,
		Whitelist: whitelist,
		Logger:    logg,
	}
}
