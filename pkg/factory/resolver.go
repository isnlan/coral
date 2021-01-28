package factory

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/snlansky/coral/pkg/discovery"
)

type NetworkResolver struct {
	mu    *sync.Mutex
	ch    <-chan []*discovery.ServiceInfo
	addrs []string
}

func NewResolver(ch <-chan []*discovery.ServiceInfo) *NetworkResolver {
	var addrs []string
	for _, info := range <-ch {
		addrs = append(addrs, info.Address)
	}
	r := &NetworkResolver{mu: new(sync.Mutex), addrs: addrs}
	go r.start()
	return r
}

func (r *NetworkResolver) Resolve() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.addrs) == 0 {
		return "", errors.New("service not find")
	}
	return r.addrs[rand.Intn(len(r.addrs))], nil
}

func (r *NetworkResolver) start() {
	for list := range r.ch {
		var addrs []string
		for _, info := range list {
			addrs = append(addrs, info.Address)
		}

		r.mu.Lock()
		r.addrs = addrs
		r.mu.Unlock()
	}
}
