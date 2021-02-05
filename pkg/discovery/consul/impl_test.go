package consul

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/snlansky/coral/pkg/discovery"

	"github.com/snlansky/coral/pkg/xgrpc"
	"github.com/stretchr/testify/assert"

	"github.com/snlansky/coral/pkg/protos"
)

func TestNew(t *testing.T) {
	var vms *protos.VMServer
	sname := discovery.MakeTypeName(vms)
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
