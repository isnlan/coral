package network

import (
	"fmt"
	"testing"
	"time"

	"github.com/snlansky/coral/pkg/protos"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	factory := New("127.0.0.1:8500")
	chain := &protos.Chain{
		Id:                 "5ffbe60730a09d3ccb722477",
		NetworkType:        "fabric",
		Name:               "chain1",
		Account:            "admin",
		Consensus:          "etcdraft",
		SignatureAlgorithm: "",
		NodeCount:          2,
		TlsEnabled:         true,
	}

	time.Sleep(time.Second)
	builder, err := factory.Builder(chain)
	assert.NoError(t, err)
	fmt.Println(builder.Build())

	time.Sleep(time.Second * 10)
	builder, err = factory.Builder(chain)
	assert.NoError(t, err)
	fmt.Println(builder)
}
