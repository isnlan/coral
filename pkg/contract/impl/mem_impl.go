package impl

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/snlansky/coral/pkg/contract/identity"

	"github.com/snlansky/coral/pkg/contract"
	"github.com/snlansky/coral/pkg/utils"
)

type memoryStub struct {
	address string
	factory *MemoryFactoryChain
	t       time.Time
}

func (m *memoryStub) GetArgs() [][]byte {
	panic("implement me")
}

func (m *memoryStub) GetTxID() string {
	return string(utils.Sha256Encode([]byte(m.t.String())))
}

func (m *memoryStub) GetChannelID() string {
	return "mem-channel"
}

func (m *memoryStub) GetAddress() (string, error) {
	addr, err := identity.AddressFromHexString(m.address)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func (m *memoryStub) GetState(key string) ([]byte, error) {
	v := m.factory.states[key]
	return v, nil
}

func (m *memoryStub) PutState(key string, value []byte) error {
	m.factory.states[key] = value
	return nil
}

func (m *memoryStub) DelState(key string) ([]byte, error) {
	v, ok := m.factory.states[key]
	if ok {
		delete(m.factory.states, key)
		return v, nil
	}
	return nil, nil
}

func (m *memoryStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return contract.CreateKey(objectType, attributes)
}

func (m *memoryStub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return contract.SplitKey(compositeKey)
}

func (m *memoryStub) GetTxTimestamp() (time.Time, error) {
	return m.t, nil
}

func (m *memoryStub) SetEvent(name string, payload []byte) error {
	m.factory.events[name] = payload
	return nil
}

func (m *memoryStub) InvokeContract(contractName string, args [][]byte, channel string) ([]byte, error) {
	panic("implement me")
}

func (m *memoryStub) GetOriginStub() interface{} {
	panic("implement me")
}

type MemoryFactoryChain struct {
	states map[string][]byte
	events map[string][]byte
}

func NewMemoryFactoryChain() *MemoryFactoryChain {
	return &MemoryFactoryChain{
		states: map[string][]byte{},
		events: map[string][]byte{},
	}
}

func (m *MemoryFactoryChain) NewStub(addr string) contract.IContractStub {
	return &memoryStub{
		address: addr,
		factory: m,
		t:       time.Now(),
	}
}

func (m *MemoryFactoryChain) Debug(prefix ...string) {
	fmt.Println("------------STATES-------------")
	var keys []string
	for k := range m.states {
		if len(prefix) == 0 {
			keys = append(keys, k)
		} else {
			if strings.HasPrefix(k, prefix[0]) {
				keys = append(keys, k)
			}
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s -> %s\n", k, string(m.states[k]))
	}

	fmt.Println("------------EVENTS-------------")
	for k, v := range m.events {
		fmt.Printf("%s -> %s\n", k, string(v))
	}
}
