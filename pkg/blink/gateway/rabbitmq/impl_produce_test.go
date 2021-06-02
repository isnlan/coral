package rabbitmq

import (
	"fmt"
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
