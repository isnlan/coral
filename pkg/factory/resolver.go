package factory

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/snlansky/coral/pkg/service_discovery"
)

type NetworkResolver struct {
	mu    *sync.Mutex
	addrs []string
}

func NewResolver() *NetworkResolver {
	return &NetworkResolver{mu: new(sync.Mutex), addrs: []string{}}
}

func (r *NetworkResolver) Resolve() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.addrs) == 0 {
		return "", errors.New("service not find")
	}
	return r.addrs[rand.Intn(len(r.addrs))], nil
}

func (r *NetworkResolver) Handle(list []*service_discovery.ServiceInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var addrs []string
	for _, info := range list {
		addrs = append(addrs, info.Address)
	}
	r.addrs = addrs
}
