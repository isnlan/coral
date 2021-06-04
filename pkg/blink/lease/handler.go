package lease

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/consul/api"

	"github.com/isnlan/coral/pkg/logging"

	"github.com/isnlan/coral/pkg/discovery"
	"github.com/isnlan/coral/pkg/utils"
)

const ns = "blink"

var errNotFindSource = errors.New("not find source")
var logger = logging.MustGetLogger("source")

type Handler struct {
	ds discovery.ServiceDiscover
}

func New(ds discovery.ServiceDiscover) *Handler {
	return &Handler{ds: ds}
}

func (h *Handler) SetAclLease(acl *AclLease) error {
	return h.SetSource(acl.ClientId, acl)
}

func (h *Handler) GetAclLease(clientId string) (*AclLease, error) {
	var acl AclLease

	err := h.GetSource(clientId, &acl)
	if err == errNotFindSource {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &acl, err
}

func (h *Handler) DeleteAclLease(acl *AclLease) error {
	return h.DeleteSource(acl.ClientId, acl)
}

func (h *Handler) DeleteAclLeaseList() error {
	return h.DeleteSourceList(&AclLease{})
}

func (h *Handler) SetChainLease(chain *ChainLease) error {
	return h.SetSource(chain.NetworkID, chain)
}

func (h *Handler) GetChainLease(networkId string) (*ChainLease, error) {
	var chain ChainLease

	err := h.GetSource(networkId, &chain)
	if err == errNotFindSource {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &chain, err
}

func (h *Handler) DeleteChainLease(chain *ChainLease) error {
	return h.DeleteSource(chain.NetworkID, chain)
}

func (h *Handler) DeleteChainLeaseList() error {
	return h.DeleteSourceList(&ChainLease{})
}

func (h *Handler) WatchChainLeaseList(ctx context.Context, ch chan<- []string) {
	keys := make(chan []string)
	prefix := utils.MakeTypeName(&ChainLease{})

	h.ds.WatchKeysByPrefix(ctx, ns, prefix, keys)

	for {
		select {
		case list := <-keys:

			var tmp []string

			for _, key := range list {
				split := strings.Split(key, ":")
				if len(split) != 2 {
					logger.Errorf("find error key: %s not match [%s:NetworkID], we will skip it", key, prefix)
				} else {
					tmp = append(tmp, split[1])
				}
			}
			ch <- tmp
		case <-ctx.Done():
			logger.Warn("context done, stop watching chain list")
			return
		}
	}
}

func (h *Handler) WatchChainLease(ctx context.Context, networkId string, ch chan *ChainLease) {
	key := fmt.Sprintf("%s:%s", utils.MakeTypeName(&ChainLease{}), networkId)
	pairs := make(chan *api.KVPair)

	h.ds.WatchKey(ctx, ns, key, pairs)

	for {
		select {
		case pair := <-pairs:
			if pair == nil || len(pair.Value) == 0 {
				ch <- nil
				continue
			}

			var chain ChainLease

			err := json.Unmarshal(pair.Value, &chain)
			if err != nil {
				logger.Errorf("json unmarshal error: %w, consul key: %v", err, key)
				continue
			}

			ch <- &chain
		case <-ctx.Done():
			logger.Warn("context done, stop watching chain")
			return
		}
	}
}

func (h *Handler) WatchChannelLease(ctx context.Context, networkId, channelName string, ch chan *ChannelLease) {
	key := fmt.Sprintf("%s:%s:%s", utils.MakeTypeName(&ChannelLease{}), networkId, channelName)
	pairs := make(chan *api.KVPair)

	h.ds.WatchKey(ctx, ns, key, pairs)

	for {
		select {
		case pair := <-pairs:
			if pair == nil || len(pair.Value) == 0 {
				ch <- nil
				continue
			}

			var channel ChannelLease

			err := json.Unmarshal(pair.Value, &channel)
			if err != nil {
				logger.Errorf("json unmarshal error: %w, consul key: %v", err, key)
				continue
			}

			ch <- &channel
		case <-ctx.Done():
			logger.Warn("context done, stop watching channel")
			return
		}
	}
}

func (h *Handler) SetChannelLease(channel *ChannelLease) error {
	key := fmt.Sprintf("%s:%s", channel.NetworkID, channel.Name)

	return h.SetSource(key, channel)
}

func (h *Handler) GetChannelLease(networkId, channelName string) (*ChannelLease, error) {
	var lease ChannelLease

	key := fmt.Sprintf("%s:%s", networkId, channelName)

	err := h.GetSource(key, &lease)
	if err == errNotFindSource {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &lease, err
}

func (h *Handler) DeleteChannelLease(channel *ChannelLease) error {
	return h.ds.DeleteKey(ns, fmt.Sprintf("%s:%s:%s", utils.MakeTypeName(&ChannelLease{}), channel.NetworkID, channel.Name))
}

func (h *Handler) DeleteChannelLeaseListByNetworkId(networkId string) error {
	return h.ds.DeleteKeyByPrefix(ns, fmt.Sprintf("%s:%s", utils.MakeTypeName(&ChannelLease{}), networkId))
}

func (h *Handler) DeleteChannelLeaseList() error {
	return h.ds.DeleteKeyByPrefix(ns, utils.MakeTypeName(&ChannelLease{}))
}

func (h *Handler) SetSource(key string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return h.ds.SetKey(ns, fmt.Sprintf("%s:%s", utils.MakeTypeName(v), key), bytes)
}

func (h *Handler) GetSource(key string, v interface{}) error {
	bytes, err := h.ds.GetKey(ns, fmt.Sprintf("%s:%s", utils.MakeTypeName(v), key))
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		return errNotFindSource
	}

	return json.Unmarshal(bytes, v)
}

func (h *Handler) DeleteSource(key string, v interface{}) error {
	return h.ds.DeleteKey(ns, fmt.Sprintf("%s:%s", utils.MakeTypeName(v), key))
}

func (h *Handler) DeleteSourceList(v interface{}) error {
	return h.ds.DeleteKeyByPrefix(ns, utils.MakeTypeName(v))
}
