package common

import "time"

type APIAuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

type APIAuthResponse struct {
	Ok bool `json:"ok"`
}

type APIListRequest struct {
	Address string `json:"address"`
}

type APIException struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AddressItem struct {
	ID      int       `json:"id"`
	Address string    `json:"address"`
	Created time.Time `json:"created"`
}
