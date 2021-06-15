package mock

import (
	"github.com/isnlan/coral/pkg/blink/network"
	"github.com/isnlan/coral/pkg/protos"
)

var _ network.Factory = &MockFacotry{}

type MockFacotry struct {
	b network.Builder
}

func (m *MockFacotry) SetBuilder(b network.Builder) {
	m.b = b
}

func (m *MockFacotry) Builder(chain *protos.Chain) (network.Builder, error) {
	if m.b != nil {
		return m.b, nil
	}

	return &MockBuilder{}, nil
}

func (m *MockFacotry) Close() {
	// nothing to do
}

type MockBuilder struct {
	channel string
}

func (m *MockBuilder) SetChannel(channel string) network.Builder {
	m.channel = channel
	return m
}

func (m *MockBuilder) Build() network.Network {
	return &DefaultMockNetwork{}
}
