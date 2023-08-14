package tests

import "github.com/evgen1067/anti-bruteforce/internal/common"

var (
	authReq = common.APIAuthRequest{
		Login:    "l",
		Password: "p",
		IP:       "127.0.0.1/25",
	}
	authSuccess = &common.APIAuthResponse{
		Ok: true,
	}
	authFail = &common.APIAuthResponse{
		Ok: false,
	}
	address = "127.0.0.1/25"
)
