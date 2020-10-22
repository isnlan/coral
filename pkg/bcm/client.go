package bcm

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
	"github.com/snlansky/coral/pkg/entity"
	"github.com/snlansky/coral/pkg/errors"
	"github.com/snlansky/coral/pkg/response"
)

type Client struct {
	baseUrl string
}

func New(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
	}
}

func (c *Client) AclQuery(clientId string) (*entity.AclClient, error) {
	var resp response.JsonResponse
	var acl entity.AclClient
	resp.Data = &acl

	data := map[string]string{
		"client_id": clientId,
	}

	_, _, errs := gorequest.New().Post(fmt.Sprintf("%s/api/private/acl/query", c.baseUrl)).
		Send(data).
		EndStruct(&resp)
	if errs != nil && len(errs) != 0 {
		return nil, errs[0]
	}

	if resp.ErrorCode != response.SuccessCode {
		return nil, errors.Errorf("request acl error: %s", resp.Description)
	}

	return &acl, nil
}

func (c *Client) ChainLease(chainId string) (*entity.Lease, error) {
	var resp response.JsonResponse
	var lease entity.Lease
	resp.Data = &lease

	_, _, errs := gorequest.New().Get(fmt.Sprintf("%s/api/private/chains/lease?chain_id=%s", c.baseUrl, chainId)).
		EndStruct(&resp)
	if errs != nil && len(errs) != 0 {
		return nil, errs[0]
	}

	if resp.ErrorCode != response.SuccessCode {
		return nil, errors.Errorf("request acl error: %s", resp.Description)
	}

	return &lease, nil
}
