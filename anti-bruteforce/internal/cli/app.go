package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/evgen1067/anti-bruteforce/internal/common"
)

type AppCLI struct {
	ctx    context.Context
	client *http.Client
}

func NewAppCLI(ctx context.Context) *AppCLI {
	return &AppCLI{
		ctx:    ctx,
		client: &http.Client{},
	}
}

func (a *AppCLI) AddToList(list, addr string) error {
	body, err := a.body(addr)
	if err != nil {
		return err
	}
	url, err := a.url(list)
	if err != nil {
		return err
	}
	err = a.makeRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppCLI) DeleteFromList(list, addr string) error {
	body, err := a.body(addr)
	if err != nil {
		return err
	}
	url, err := a.url(list)
	if err != nil {
		return err
	}
	err = a.makeRequest(http.MethodDelete, url, body)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppCLI) url(list string) (string, error) {
	switch strings.ToLower(list) {
	case "blacklist":
		return common.BlacklistURL, nil
	case "whitelist":
		return common.WhitelistURL, nil
	default:
		return "", fmt.Errorf("unsupported list type")
	}
}

func (a *AppCLI) body(addr string) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	request := common.APIListRequest{Address: addr}
	err := json.NewEncoder(b).Encode(&request)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (a *AppCLI) ResetBucket() error {
	err := a.makeRequest(http.MethodPost, common.ResetURL, new(bytes.Buffer))
	if err != nil {
		return err
	}
	return nil
}

func (a *AppCLI) makeRequest(httpMethod, url string, body *bytes.Buffer) error {
	baseURL := "http://localhost:8888"
	req, err := http.NewRequestWithContext(a.ctx, httpMethod, baseURL+url, body)
	if err != nil {
		return err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
