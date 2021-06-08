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

	//err = impl.DeleteAclLeaseList()
	//assert.NoError(t, err)
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
		IsRunning:   true,
		TlsEnabled:  false,
	}
	err = impl.SetChainLease(lease)
	assert.NoError(t, err)

	lease2, err := impl.GetChainLease(lease.NetworkID)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(lease, lease2))

	channel := &ChannelLease{
		ID:         primitive.NewObjectID().Hex(),
		NetworkID:  lease.NetworkID,
		Name:       "mychannel",
		Endpoint:   "f/b",
		IsRunning:  true,
		SyncEnable: true,
		SyncDB:     "",
	}

	err = impl.SetChannelLease(channel)
	assert.NoError(t, err)

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

func TestBlinkSourceImpl_WatchChainIDList(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)

	ch := make(chan []string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	go impl.WatchChainIDList(ctx, ch)

	go func() {
		for keys := range ch {
			fmt.Println(keys)
		}
	}()
	<-ctx.Done()
}

func TestHandler_WatchChainLeaseByID(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)

	ch := make(chan *ChainLease)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go impl.WatchChainLeaseByID(ctx, "60b995c42b778ea9b3d0c2fc", ch)

	go func() {
		for lease := range ch {
			if lease == nil {
				fmt.Println("chain lease is nil")
				time.Sleep(time.Second * 2)
				cancel()
			} else {
				fmt.Printf("chain lease changed: %v\n", lease)
			}
		}
	}()

	<-ctx.Done()
}

func TestHandler_WatchChannelLeaseByName(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)

	ch := make(chan *ChannelLease)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go impl.WatchChannelLeaseByName(ctx, "60b995c42b778ea9b3d0c2fc", "testchannel", ch)

	go func() {
		for lease := range ch {
			if lease == nil {
				fmt.Println("chain lease is nil")
				time.Sleep(time.Second * 2)
				cancel()
			} else {
				fmt.Printf("chain lease changed: %v\n", lease)
			}
		}
	}()

	<-ctx.Done()
}

func TestHandler_WatchChannelLeaseList(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)

	ch := make(chan []*ChannelLease)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go impl.WatchChannelLeaseList(ctx, ch)

	go func() {
		for list := range ch {
			if len(list) == 0 {
				fmt.Println("channel list is nil")
				time.Sleep(time.Second * 2)
				cancel()
			} else {
				for _, channel := range list {
					fmt.Printf("channel: %+#v\n", channel)
				}
			}
			fmt.Println("")
		}
	}()

	<-ctx.Done()
}

func TestHandler_WatchChainLeaseList(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)

	ch := make(chan []*ChainLease)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go impl.WatchChainLeaseList(ctx, ch)

	go func() {
		for list := range ch {
			if len(list) == 0 {
				fmt.Println("chain list is nil")
				time.Sleep(time.Second * 2)
				cancel()
			} else {
				for _, chain := range list {
					fmt.Printf("chain: %+#v\n", chain)
				}
			}
			fmt.Println("")
		}
	}()

	<-ctx.Done()
}

func TestHandler_WatchChannelLeaseListByNetworkId(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)

	ch := make(chan []*ChannelLease)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go impl.WatchChannelLeaseListByNetworkId(ctx, "5fa26e237a065ea42dc99300", ch)

	go func() {
		for list := range ch {
			if len(list) == 0 {
				fmt.Println("channel list is nil")
				time.Sleep(time.Second * 2)
				cancel()
			} else {
				for _, channel := range list {
					fmt.Printf("channel: %+#v\n", channel)
				}
			}
			fmt.Println("")
		}
	}()

	<-ctx.Done()
}
func TestHandler_WatchACLLeaseList(t *testing.T) {
	c, err := consul.New("127.0.0.1:8500")
	assert.NoError(t, err)

	impl := New(c)

	ch := make(chan []*AclLease)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go impl.WatchACLLeaseList(ctx, ch)

	go func() {
		for list := range ch {
			if len(list) == 0 {
				fmt.Println("acl list is nil")
				time.Sleep(time.Second * 2)
				cancel()
			} else {
				for _, acl := range list {
					fmt.Printf("acl: %+#v\n", acl)
				}
			}
			fmt.Println("")
		}
	}()

	<-ctx.Done()
}
