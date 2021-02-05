package consul

import (
	"context"
	"fmt"
	"time"

	"github.com/snlansky/coral/pkg/discovery"

	"github.com/snlansky/coral/pkg/logging"

	"google.golang.org/grpc/health"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var logger = logging.MustGetLogger("consul")

type consulImpl struct {
	client *api.Client
}

func New(url string) (*consulImpl, error) {
	config := api.DefaultConfig()
	config.Address = url
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &consulImpl{client: client}, nil
}

func (c *consulImpl) ServiceRegister(name, address string, port int, tags ...string) (discovery.Deregister, error) {
	svr := &api.AgentServiceRegistration{
		Name:    name,
		ID:      fmt.Sprintf("%s-%s:%d", name, address, port),
		Address: address,
		Port:    port,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", address, port),
			Interval:                       "3s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	register := func() error {
		return c.client.Agent().ServiceRegister(svr)
	}

	if err := register(); err != nil {
		return nil, err
	}
	logger.Infof("service: %s register success!", svr.ID)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		tick := time.NewTicker(time.Minute)
		for {
			select {
			case <-tick.C:
				_, _, err := c.client.Agent().Service(svr.ID, &api.QueryOptions{RequireConsistent: true})
				if err != nil {
					if err := register(); err != nil {
						logger.Errorf("service: %s register error", svr.ID)
					} else {
						logger.Infof("service: %s register success!", svr.ID)
					}
				}
			case <-ctx.Done():
				err := c.client.Agent().ServiceDeregister(svr.ID)
				if err != nil {
					logger.Errorf("service: %s deregister error: %v", svr.ID, err)
				} else {
					logger.Warnf("service: %s deregister", svr.ID)
				}
				return
			}
		}
	}()

	return discovery.Deregister(cancel), nil
}

func (c *consulImpl) RegisterHealthServer(s *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
}

func (c *consulImpl) WatchService(ctx context.Context, name string, tag string, ch chan<- []*discovery.ServiceInfo) {
	go func() {
		var waitIndex uint64
		for {
			entries, lastIndex, err := c.serviceEntriesWatch(name, tag, waitIndex)
			if err != nil {
				logger.Errorf("service entries watch error: %v", err)
				time.Sleep(30 * time.Second)
				continue
			}
			if waitIndex != lastIndex {
				waitIndex = lastIndex
				ch <- entries
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
}

func (c *consulImpl) serviceEntriesWatch(name, tag string, waitIndex uint64) ([]*discovery.ServiceInfo, uint64, error) {
	opt := &api.QueryOptions{
		RequireConsistent: true,
		WaitIndex:         waitIndex,
		WaitTime:          time.Minute * 3,
	}

	entries, meta, err := c.client.Health().Service(name, tag, true, opt)
	if err != nil {
		return nil, 0, err
	}

	var list []*discovery.ServiceInfo
	for _, entry := range entries {
		if entry.Service == nil {
			continue
		}
		info := &discovery.ServiceInfo{
			ID:      entry.Service.ID,
			Address: fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port),
			Tags:    entry.Service.Tags,
		}

		list = append(list, info)
	}

	return list, meta.LastIndex, nil
}
