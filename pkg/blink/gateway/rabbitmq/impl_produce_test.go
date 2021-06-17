package rabbitmq

import (
	"fmt"
	"sync"
	"testing"
	"time"

	gateway2 "github.com/isnlan/coral/pkg/blink/gateway"

	"github.com/stretchr/testify/assert"
)

func TestNewProduce(t *testing.T) {
	url := "amqp://admin:admin@localhost:5672/"
	go func() {
		consume := NewConsume(url, &mockConsume{})
		consume.Start()
	}()
	produce := NewProduce(url)

	time.Sleep(time.Second)
	err := produce.APIUpload(&gateway2.API{
		ID:      "id1",
		Scheme:  "http",
		Method:  "get",
		Path:    "/ping",
		AppName: "myapp",
		APIName: "PING",
		APIType: "拼",
		DocURL:  "",
	})
	assert.NoError(t, err)

	err = produce.APIUpload(&gateway2.API{
		ID:      "id1",
		Scheme:  "http",
		Method:  "get",
		Path:    "/ping",
		AppName: "myapp",
		APIName: "PING",
		DocURL:  "",
	})
	assert.NoError(t, err)

	err = produce.APICallRecord(&gateway2.APICallEntity{
		APIID:    "id1",
		Latency:  10,
		HttpCode: 200,
		ClientID: "c1",
	})
	assert.NoError(t, err)

	err = produce.ContractCallRecord(&gateway2.ContractCallEntity{
		ClientID:  "ssss",
		Address:   "adress",
		ChainID:   "c1",
		ChannelID: "chann1",
		Contract:  "tv",
	})
	assert.NoError(t, err)
	time.Sleep(time.Minute)
}

func BenchmarkProduceImpl_APICallRecord(b *testing.B) {
	url := "amqp://admin:admin@localhost:5672/"
	produce := NewProduce(url)

	apis := []string{"df05961e491bb6a77edeb7fc", "f15883f21409ed3f0eb34cff"}

	for i := 0; i < b.N; i++ {
		err := produce.APICallRecord(&gateway2.APICallEntity{
			APIID:    apis[i%2],
			Latency:  int64(i % 1000),
			HttpCode: 200,
			ClientID: "473d78a37c640099",
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestProduceImpl_APICallRecord(t *testing.T) {
	url := "amqp://admin:admin@localhost:5672/"
	produce := NewProduce(url)

	apis := []string{"df05961e491bb6a77edeb7fc", "f15883f21409ed3f0eb34cff"}
	var wg sync.WaitGroup
	for j := 0; j < 1000; j++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				err := produce.APICallRecord(&gateway2.APICallEntity{
					APIID:    apis[i%2],
					Latency:  int64(i % 1000),
					HttpCode: 200,
					ClientID: "473d78a37c640099",
				})
				fmt.Println(j, i)
				if err != nil {
					t.Fatal(err)
				}
			}
		}(j)
	}

	wg.Wait()
}

type mockConsume struct {
}

func (m mockConsume) APIHandler(api *gateway2.API) error {
	fmt.Println("api", api)
	return nil
}

func (m mockConsume) APICallHandler(entity *gateway2.APICallEntity) error {
	fmt.Println("entity", entity)
	return nil
}

func (m mockConsume) ContractCallHandler(entity *gateway2.ContractCallEntity) error {
	fmt.Println("entity", entity)
	return nil
}
