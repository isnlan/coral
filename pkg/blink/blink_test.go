package blink

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	store = New("https://kchaintest.kingdeeresearch.com/blockchain-manager-backend")
)

func TestClientStoreGetAcl(t *testing.T) {
	info, err := store.AclQuery("5dbac7be00b59c0c")
	assert.NoError(t, err)

	fmt.Println(info)
}

func TestClientStoreGetChainLease(t *testing.T) {
	info, err := store.ChainLease("6030d1eb786694b2dd8b4e00")
	assert.NoError(t, err)

	fmt.Println(info)
}
