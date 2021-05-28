package lease

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/isnlan/coral/pkg/discovery/consul"
	"github.com/isnlan/coral/pkg/entity"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNew(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)
	acl := &entity.AclClient{
		ID:           primitive.NewObjectID(),
		Name:         "苍穹开发应用（请勿删除！！！）",
		ClientId:     "5dbac7be00b59c0c",
		ClientSecret: "447729b7724e28a9795d24a8c1500f773398b910d5c66c5058849d0136ea1220",
		Account:      "kdcloud55538816",
		Team:         "kdcloud-t5",
		ChainId:      "5fa26e237a065ea42dc993dd",
		Nodes:        nil,
		Enable:       true,
		CreateTime:   1604480693,
		Description:  "",
	}
	err = impl.SetSource(acl.ClientId, acl)
	assert.NoError(t, err)

	var acl2 entity.AclClient
	err = impl.GetSource("5dbac7be00b59c0c", &acl2)
	assert.NoError(t, err)
	fmt.Println(acl2)

	//err = impl.DeleteSource("5dbac7be00b59c0c", &entity.AclClient{})
	//assert.NoError(t, err)

	err = impl.DeleteSourceList(&entity.AclClient{})
	assert.NoError(t, err)
}

func TestBlinkSourceImpl_GetSource(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)
	acl := &AclLease{
		ID:           "acl1",
		Name:         "苍穹开发应用（请勿删除！！！）",
		ClientId:     "5dbac7be00b59c0c",
		ClientSecret: "447729b7724e28a9795d24a8c1500f773398b910d5c66c5058849d0136ea1220",
		Account:      "kdcloud55538816",
		Team:         "kdcloud-t5",
		NetworkID:    "5fa26e237a065ea42dc993dd",
		Enable:       true,
	}
	err = impl.SetAclLease(acl)
	assert.NoError(t, err)

	acl2, err := impl.GetAclLease(acl.ClientId)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(acl, acl2))

	//err = impl.DeleteAclLease(acl)
	//assert.NoError(t, err)

	err = impl.DeleteAclLeaseList()
	assert.NoError(t, err)
}

func TestBlinkSourceImpl_SetChainLease(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)

	lease := &ChainLease{
		NetworkID:   "5fa26e237a065ea42dc993dd",
		NetworkType: "fabric",
		NetworkName: "importantnetwork",
		Account:     "kdcloud55538816",
		Team:        "kdcloud-t5",
		IsRunning:   false,
		TlsEnabled:  false,
	}
	err = impl.SetChainLease(lease)
	assert.NoError(t, err)

	lease2, err := impl.GetChainLease(lease.NetworkID)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(lease, lease2))

	//err = impl.DeleteChainLease(lease)
	//assert.NoError(t, err)
	//
	//chainLease, err := impl.GetChainLease(lease.NetworkID)
	//assert.NoError(t, err)
	//assert.Nil(t, chainLease)
	//
	//err = impl.DeleteChainLeaseList()
	//assert.NoError(t, err)

}

func TestBlinkSourceImpl_WatchChainLeaseList(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)

	ch := make(chan []string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	go impl.WatchChainLeaseList(ctx, ch)

	go func() {
		for keys := range ch {
			fmt.Println(keys)
		}
	}()
	<-ctx.Done()
}