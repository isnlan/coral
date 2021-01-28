package consul

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/snlansky/coral/pkg/discovery"

	"github.com/snlansky/coral/pkg/net"
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
			server, err := net.NewServer(fmt.Sprintf("0.0.0.0:%d", port))
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
