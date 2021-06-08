package lease

import (
	"context"
	"sync"
)

type Register struct {
	h        *Handler
	rw       sync.RWMutex
	acls     map[string]*AclLease
	chains   map[string]*ChainLease
	channels map[string]*ChannelLease
}

func NewRegistry(h *Handler) *Register {
	return &Register{
		h:        h,
		rw:       sync.RWMutex{},
		acls:     map[string]*AclLease{},
		chains:   map[string]*ChainLease{},
		channels: map[string]*ChannelLease{},
	}
}

func (r *Register) Start(ctx context.Context) {
	chainCh := make(chan []*ChainLease)
	channelCh := make(chan []*ChannelLease)
	aclCh := make(chan []*AclLease)

	go r.h.WatchChainLeaseList(ctx, chainCh)
	go r.h.WatchChannelLeaseList(ctx, channelCh)
	go r.h.WatchACLLeaseList(ctx, aclCh)

	for {
		select {
		case chains := <-chainCh:
			r.filerChains(chains)
		case channels := <-channelCh:
			r.filerChannels(channels)
		case acls := <-aclCh:
			r.filerACLs(acls)
		case <-ctx.Done():
			return
		}
	}
}

func (r *Register) filerChains(chains []*ChainLease) {
	r.rw.Lock()
	defer r.rw.Unlock()

	for id := range r.chains {
		delete(r.chains, id)
	}

	for _, chain := range chains {
		r.chains[chain.NetworkID] = chain
	}
}

func (r *Register) filerChannels(channels []*ChannelLease) {
	r.rw.Lock()
	defer r.rw.Unlock()

	for id := range r.channels {
		delete(r.channels, id)
	}

	for _, channel := range channels {
		r.channels[channel.NetworkID+":"+channel.Name] = channel
	}
}

func (r *Register) filerACLs(acls []*AclLease) {
	r.rw.Lock()
	defer r.rw.Unlock()

	for id := range r.acls {
		delete(r.acls, id)
	}

	for _, acl := range acls {
		r.acls[acl.ClientId] = acl
	}
}

func (r *Register) GetChainByID(networkID string) *ChainLease {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.chains[networkID]
}

func (r *Register) GetChannelBy(networkID string, channelName string) *ChannelLease {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.channels[networkID+":"+channelName]
}

func (r *Register) GetAclByClientID(clientId string) *AclLease {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.acls[clientId]
}
