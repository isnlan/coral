package consul

import (
	"context"
	"fmt"
	"time"

	"github.com/isnlan/coral/pkg/discovery"

	"github.com/isnlan/coral/pkg/logging"

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

func (c *consulImpl) HTTPServiceRegister(name, address string, port int, tags ...string) (discovery.Deregister, error) {
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d%s", address, port, discovery.HTTPHealthCheckRouter),
		Interval:                       "3s",
		Timeout:                        "5s",
		DeregisterCriticalServiceAfter: "30s",
	}

	return c.register(name, address, port, check, tags...)
}

func (c *consulImpl) ServiceRegister(name, address string, port int, tags ...string) (discovery.Deregister, error) {
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", address, port),
		Interval:                       "3s",
		Timeout:                        "5s",
		DeregisterCriticalServiceAfter: "30s",
	}

	return c.register(name, address, port, check, tags...)
}

func (c *consulImpl) register(name, address string, port int, check *api.AgentServiceCheck, tags ...string) (discovery.Deregister, error) {
	svr := &api.AgentServiceRegistration{
		Name:    name,
		ID:      fmt.Sprintf("%s-%s:%d", name, address, port),
		Address: address,
		Port:    port,
		Tags:    tags,
		Check:   check,
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

func (c *consulImpl) SetKey(ns, key string, value []byte) error {
	pair := &api.KVPair{
		Key:       key,
		Value:     value,
		Namespace: ns,
	}
	opt := &api.WriteOptions{
		// Namespaces are a Consul Enterprise feature
		// Namespace: ns,
	}
	_, err := c.client.KV().Put(pair, opt)
	return err
}

func (c *consulImpl) GetKey(ns, key string) ([]byte, error) {
	pair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, nil
	}

	return pair.Value, nil
}

func (c *consulImpl) GetKeys(prefix string) ([]string, error) {
	keys, _, err := c.client.KV().Keys(prefix, "/", nil)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (c *consulImpl) GetList(prefix string) ([]*api.KVPair, error) {
	list, _, err := c.client.KV().List(prefix, nil)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *consulImpl) DeleteKey(ns, key string) error {
	_, err := c.client.KV().Delete(key, nil)
	return err
}

func (c *consulImpl) DeleteKeyByPrefix(ns, prefix string) error {
	_, err := c.client.KV().DeleteTree(prefix, nil)
	return err
}

func (c *consulImpl) WatchKey(ctx context.Context, ns, key string, ch chan<- *api.KVPair) {
	go func() {
		var waitIndex uint64
		for {
			pair, lastIndex, err := c.getKeyByIndex(ns, key, waitIndex)
			if err != nil {
				logger.Errorf("key watch error: %v", err)
				time.Sleep(30 * time.Second)
				return
			}

			if waitIndex != lastIndex {
				waitIndex = lastIndex
				// pair maybe is null, but we not care
				ch <- pair
			}

			select {
			case <-ctx.Done():
				logger.Infof("stop watching key: %s", key)
				return
			default:
			}
		}
	}()
}

func (c *consulImpl) getKeyByIndex(ns, key string, waitIndex uint64) (*api.KVPair, uint64, error) {
	opt := &api.QueryOptions{
		RequireConsistent: true,
		WaitIndex:         waitIndex,
		WaitTime:          time.Minute,
	}

	pair, meta, err := c.client.KV().Get(key, opt)
	if err != nil {
		return nil, 0, err
	}

	return pair, meta.LastIndex, nil
}

func (c *consulImpl) WatchKeysByPrefix(ctx context.Context, ns, prefix string, ch chan<- []string) {
	go func() {
		var waitIndex uint64
		for {
			keys, lastIndex, err := c.getKeysByPrefixAndIndex(ns, prefix, waitIndex)
			if err != nil {
				logger.Errorf("key watch error: %v", err)
				time.Sleep(30 * time.Second)
				return
			}
			if waitIndex != lastIndex {
				waitIndex = lastIndex
				ch <- keys
			}

			select {
			case <-ctx.Done():
				logger.Infof("stop watching keys by prefix: %s", prefix)
				return
			default:
			}
		}
	}()
}

func (c *consulImpl) getKeysByPrefixAndIndex(ns, prefix string, waitIndex uint64) ([]string, uint64, error) {
	opt := &api.QueryOptions{
		RequireConsistent: true,
		WaitIndex:         waitIndex,
		WaitTime:          time.Minute,
	}

	keys, meta, err := c.client.KV().Keys(prefix, "/", opt)
	if err != nil {
		return nil, 0, err
	}

	return keys, meta.LastIndex, nil
}

func (c *consulImpl) WatchValuesByKeyPrefix(ctx context.Context, ns, prefix string, ch chan<- []*api.KVPair) {
	go func() {
		var waitIndex uint64
		for {
			kvs, lastIndex, err := c.getValuesByKeyPrefixAndIndex(ns, prefix, waitIndex)
			if err != nil {
				logger.Errorf("key watch error: %v", err)
				time.Sleep(30 * time.Second)
				return
			}
			if waitIndex != lastIndex {
				waitIndex = lastIndex
				ch <- kvs
			}

			select {
			case <-ctx.Done():
				logger.Infof("stop watching values by prefix: %s", prefix)
				return
			default:
			}
		}
	}()
}

func (c *consulImpl) getValuesByKeyPrefixAndIndex(ns, prefix string, waitIndex uint64) ([]*api.KVPair, uint64, error) {
	opt := &api.QueryOptions{
		RequireConsistent: true,
		WaitIndex:         waitIndex,
		WaitTime:          time.Minute,
	}

	kvs, meta, err := c.client.KV().List(prefix, opt)
	if err != nil {
		return nil, 0, err
	}

	return kvs, meta.LastIndex, nil
}

func (c *consulImpl) LockKey(key string) (*api.Lock, error) {
	return c.client.LockKey(key)
}
