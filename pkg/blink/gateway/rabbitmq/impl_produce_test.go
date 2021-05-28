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
	err := produce.ApiUpload(&gateway2.Api{
		Id:      "id1",
		Scheme:  "http",
		Method:  "get",
		Path:    "/ping",
		AppName: "myapp",
		ApiName: "PING",
		ApiType: "拼",
		DocUrl:  "",
	})
	assert.NoError(t, err)

	err = produce.ApiUpload(&gateway2.Api{
		Id:      "id1",
		Scheme:  "http",
		Method:  "get",
		Path:    "/ping",
		AppName: "myapp",
		ApiName: "PING",
		DocUrl:  "",
	})
	assert.NoError(t, err)

	err = produce.ApiCallRecord(&gateway2.ApiCallEntity{
		ApiId:    "id1",
		Latency:  10,
		HttpCode: 200,
		ClientId: "c1",
	})
	assert.NoError(t, err)

	err = produce.ContractCallRecord(&gateway2.ContractCallEntity{
		ClientId:  "ssss",
		Address:   "adress",
		ChainId:   "c1",
		ChannelId: "chann1",
		Contract:  "tv",
	})
	assert.NoError(t, err)
	time.Sleep(time.Minute)
}

type mockConsume struct {
}

func (m mockConsume) ApiHandler(api *gateway2.Api) error {
	fmt.Println("api", api)
	return nil
}

func (m mockConsume) ApiCallHandler(entity *gateway2.ApiCallEntity) error {
	fmt.Println("entity", entity)
	return nil
}

func (m mockConsume) ContractCallHandler(entity *gateway2.ContractCallEntity) error {
	fmt.Println("entity", entity)
	return nil
}
