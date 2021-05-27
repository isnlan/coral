package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/isnlan/coral/pkg/utils"

	"github.com/hashicorp/consul/api"

	"github.com/isnlan/coral/pkg/discovery"

	"github.com/isnlan/coral/pkg/xgrpc"
	"github.com/stretchr/testify/assert"

	"github.com/isnlan/coral/pkg/protos"
)

func TestNew(t *testing.T) {
	var vms *protos.VMServer
	sname := utils.MakeTypeName(vms)
	fmt.Println(sname)

	tags := []string{"fabric", "silk", "fisco"}

	go func() {
		for i := 0; i < 5; i++ {
			client, err := New("127.0.0.1:8500")
			assert.NoError(t, err)

			port := 7000 + i
			server, err := xgrpc.NewServer(fmt.Sprintf("0.0.0.0:%d", port))
			assert.NoError(t, err)

			client.RegisterHealthServer(server.Server())

			ip, err := discovery.GetLocalIP()
			assert.NoError(t, err)

			cancel, err := client.ServiceRegister(sname, ip, port, tags[rand.Intn(len(tags))])
			assert.NoError(t, err)

			go func() {
				time.Sleep(time.Duration(rand.Intn(40)) * time.Second)
				cancel()
				//t.Logf("%s close", id)
			}()

			go server.Start()
		}
	}()

	handler := make(chan []*discovery.ServiceInfo)

	go func() {
		client, err := New("127.0.0.1:8500")
		assert.NoError(t, err)

		time.Sleep(time.Second)
		client.WatchService(context.Background(), sname, "silk", handler)
	}()

	go func() {
		for infos := range handler {
			for _, e := range infos {
				fmt.Printf("%v ", e)
			}
			fmt.Println()
		}
	}()

	time.Sleep(time.Minute)
}

func TestNew2(t *testing.T) {
	client, err := New("127.0.0.1:8500")
	assert.NoError(t, err)

	port := 7000
	server, err := xgrpc.NewServer(fmt.Sprintf("0.0.0.0:%d", port))
	assert.NoError(t, err)

	client.RegisterHealthServer(server.Server())
	go server.Start()
	opt := &api.QueryOptions{
		RequireConsistent: true,
	}

	ip, err := discovery.GetLocalIP()
	assert.NoError(t, err)

	cancel, err := client.ServiceRegister("aac", ip, port, "Fsf")
	assert.NoError(t, err)

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
		//t.Logf("%s close", id)
	}()

	for {
		time.Sleep(time.Second)
		_, _, err := client.client.Agent().Service("aac-172.20.158.73:7000", opt)
		fmt.Println(err)

	}
}

func TestConsulImpl_GetKey(t *testing.T) {
	value1 := map[string]interface{}{"name": "lucy", "age": 2}
	bytes1, _ := json.Marshal(value1)
	value2 := map[string]interface{}{"name": "lili", "age": 20}
	bytes2, _ := json.Marshal(value2)

	client, err := New("127.0.0.1:8500")
	assert.NoError(t, err)
	err = client.SetKey("blink:chain", "chain1", bytes1)
	assert.NoError(t, err)

	err = client.SetKey("blink:acl", "chain1", bytes2)
	assert.NoError(t, err)

	v, err := client.GetKey("blink:acl", "chain1")
	assert.NoError(t, err)
	assert.Equal(t, v, bytes2)
}

func TestConsulImpl_GetKeys(t *testing.T) {
	value1 := map[string]interface{}{"name": "lucy", "age": 2}
	bytes1, _ := json.Marshal(value1)

	client, err := New("127.0.0.1:8500")
	assert.NoError(t, err)

	err = client.SetKey("ns", "blink:chain2", bytes1)
	assert.NoError(t, err)

	keys, err := client.GetList("blink")
	assert.NoError(t, err)
	for _, k := range keys {
		fmt.Printf("-> %+#v\n", k)
	}

	key, err := client.GetKeys("blink")
	assert.NoError(t, err)
	fmt.Println(key)
}

func TestConsulImpl_WatchKey(t *testing.T) {
	client, err := New("127.0.0.1:8500")
	assert.NoError(t, err)

	c := make(chan *api.KVPair)
	ctx, cancel := context.WithCancel(context.Background())
	client.WatchKey(ctx, "", "blink:chain1", c)

	for k := range c {
		fmt.Printf("--> %+#v\n ", k)
		if k == nil {
			break
		}
	}
	cancel()

	time.Sleep(time.Second * 10)
}

func TestConsulImpl_WatchKeysByPrefix(t *testing.T) {
	client, err := New("127.0.0.1:8500")
	assert.NoError(t, err)

	c := make(chan []string)
	ctx, cancel := context.WithCancel(context.Background())
	client.WatchKeysByPrefix(ctx, "", "blink", c)

	for k := range c {
		fmt.Printf("--> %+#v\n ", k)
		if k == nil {
			break
		}
	}
	cancel()

	time.Sleep(time.Second * 10)
}

func TestConsulImpl_DelKey(t *testing.T) {
	client, err := New("127.0.0.1:8500")
	assert.NoError(t, err)

	err = client.DeleteKey("", "blink")
	assert.NoError(t, err)

}
