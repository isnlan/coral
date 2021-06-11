package lease

import (
	"context"
	"reflect"
	"sync"

	"github.com/isnlan/coral/pkg/discovery"
)

type EventHub interface {
	Start(ctx context.Context)
	GetChainByID(networkID string) *ChainLease
	GetChannelBy(networkID string, channelName string) *ChannelLease
	GetAclByClientID(clientId string) *AclLease
}

type EventHandler interface {
	Create(v interface{})
	Delete(v interface{})
	Update(v interface{})
}

type eventHubImpl struct {
	rw       sync.RWMutex
	h        *Handler
	event    EventHandler
	acls     map[string]*AclLease
	chains   map[string]*ChainLease
	channels map[string]*ChannelLease
}

func NewEventHubImpl(ds discovery.ServiceDiscover, event EventHandler) EventHub {
	return &eventHubImpl{
		rw:       sync.RWMutex{},
		h:        New(ds),
		event:    event,
		acls:     map[string]*AclLease{},
		chains:   map[string]*ChainLease{},
		channels: map[string]*ChannelLease{},
	}
}

func (r *eventHubImpl) Start(ctx context.Context) {
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

func (r *eventHubImpl) filerChains(chains []*ChainLease) {
	r.rw.Lock()
	defer r.rw.Unlock()

	for id, old := range r.chains {
		find := false
		for _, chain := range chains {
			if id == chain.UniqueID() {
				find = true
				if !reflect.DeepEqual(old, chain) {
					r.chains[id] = chain
					r.update(chain)
				}
			}
		}

		if !find {
			delete(r.chains, id)
			r.delete(old)
		}
	}

	for _, chain := range chains {
		if _, find := r.chains[chain.UniqueID()]; !find {
			r.chains[chain.UniqueID()] = chain
			r.create(chain)
		}
	}
}

func (r *eventHubImpl) filerChannels(channels []*ChannelLease) {
	r.rw.Lock()
	defer r.rw.Unlock()

	for id, old := range r.channels {
		find := false
		for _, channel := range channels {
			if id == channel.UniqueID() {
				find = true
				if !reflect.DeepEqual(old, channel) {
					r.channels[id] = channel
					r.update(channel)
				}
			}
		}

		if !find {
			delete(r.channels, id)
			r.delete(old)
		}
	}

	for _, channel := range channels {
		if _, find := r.channels[channel.UniqueID()]; !find {
			r.channels[channel.UniqueID()] = channel
			r.create(channel)
		}
	}
}

func (r *eventHubImpl) filerACLs(acls []*AclLease) {
	r.rw.Lock()
	defer r.rw.Unlock()

	for id, old := range r.acls {
		find := false
		for _, acl := range acls {
			if id == acl.UniqueID() {
				find = true
				if !reflect.DeepEqual(old, acl) {
					r.acls[id] = acl
					r.update(acl)
				}
			}
		}

		if !find {
			delete(r.acls, id)
			r.delete(old)
		}
	}

	for _, acl := range acls {
		if _, find := r.acls[acl.UniqueID()]; !find {
			r.acls[acl.UniqueID()] = acl
			r.create(acl)
		}
	}
}

func (r *eventHubImpl) create(v interface{}) {
	if r.event != nil {
		r.event.Create(v)
	}
}

func (r *eventHubImpl) delete(v interface{}) {
	if r.event != nil {
		r.event.Delete(v)
	}
}

func (r *eventHubImpl) update(v interface{}) {
	if r.event != nil {
		r.event.Update(v)
	}
}

func (r *eventHubImpl) GetChainByID(networkID string) *ChainLease {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.chains[networkID]
}

func (r *eventHubImpl) GetChannelBy(networkID string, channelName string) *ChannelLease {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.channels[networkID+":"+channelName]
}

func (r *eventHubImpl) GetAclByClientID(clientId string) *AclLease {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.acls[clientId]
}
