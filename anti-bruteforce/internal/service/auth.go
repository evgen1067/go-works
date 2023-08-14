package service

import (
	"context"
	"github.com/evgen1067/anti-bruteforce/internal/common"
)

type AuthService struct {
	ctx       context.Context
	blacklist Blacklist
	whitelist Whitelist
	bucket    LeakyBucket
}

func NewAuthService(ctx context.Context,
	blacklist Blacklist,
	whitelist Whitelist,
	bucket LeakyBucket,
) *AuthService {
	return &AuthService{
		ctx:       ctx,
		blacklist: blacklist,
		whitelist: whitelist,
		bucket:    bucket,
	}
}

func (a *AuthService) Authorize(req common.APIAuthRequest) bool {
	bl, err := a.blacklist.ExistsInBlacklist(req.IP)
	if err != nil {
		return false
	}
	if bl == true {
		return false
	}

	wh, err := a.whitelist.ExistsInWhitelist(req.IP)
	if err != nil {
		return false
	}
	if wh == true {
		return true
	}
	return a.bucket.Add(req)
}

func (a *AuthService) ResetBucket() {
	a.bucket.ResetBucket()
}
