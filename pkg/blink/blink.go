package blink

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
	"github.com/snlansky/coral/pkg/entity"
	"github.com/snlansky/coral/pkg/errors"
	"github.com/snlansky/coral/pkg/response"
)

type Blink interface {
	AclQuery(clientId string) (*entity.AclClient, error)
	ChainLease(chainId string) (*entity.Lease, error)
	CallRecord(data interface{}) error
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

func (c *blinkImpl) CallRecord(data interface{}) error {
	var resp response.Response

	_, _, errs := gorequest.New().Post(fmt.Sprintf("%s/api/private/calls/record", c.baseUrl)).
		Send(data).
		EndStruct(&resp)
	if len(errs) != 0 {
		return errs[0]
	}

	if resp.ErrorCode != response.SuccessCode {
		return errors.Errorf("request call record error: %s", resp.Description)
	}

	return nil
}
