package blink

import (
	"fmt"

	"github.com/isnlan/coral/pkg/entity"
	"github.com/isnlan/coral/pkg/errors"
	"github.com/isnlan/coral/pkg/response"
	"github.com/parnurzeal/gorequest"
)

type Blink interface {
	AclQuery(clientId string) (*entity.AclClient, error)
	ChainLease(chainId string) (*entity.Lease, error)
}

type blinkImpl struct {
	baseUrl string
}

func New(baseUrl string) Blink {
	return &blinkImpl{
		baseUrl: baseUrl,
	}
}

func (c *blinkImpl) AclQuery(clientId string) (*entity.AclClient, error) {
	var resp response.Response
	var acl entity.AclClient
	resp.Data = &acl

	data := map[string]string{
		"client_id": clientId,
	}

	_, _, errs := gorequest.New().Post(fmt.Sprintf("%s/api/private/acl/query", c.baseUrl)).
		Send(data).
		EndStruct(&resp)
	if len(errs) != 0 {
		return nil, errs[0]
	}

	if resp.ErrorCode != response.SuccessCode {
		return nil, errors.Errorf("request acl error: %s", resp.Description)
	}

	return &acl, nil
}

func (c *blinkImpl) ChainLease(chainId string) (*entity.Lease, error) {
	var resp response.Response
	var lease entity.Lease
	resp.Data = &lease

	_, _, errs := gorequest.New().Get(fmt.Sprintf("%s/api/private/chains/lease?chain_id=%s", c.baseUrl, chainId)).
		EndStruct(&resp)
	if len(errs) != 0 {
		return nil, errs[0]
	}

	if resp.ErrorCode != response.SuccessCode {
		return nil, errors.Errorf("request acl error: %s", resp.Description)
	}

	return &lease, nil
}
