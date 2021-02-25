package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isnlan/coral/pkg/gateway"
	"github.com/isnlan/coral/pkg/gateway/rabbitmq"
	"github.com/isnlan/coral/pkg/trace"
	"github.com/isnlan/coral/pkg/xgin"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	go func() {
		url := "amqp://admin:admin@localhost:5672/"
		consume := rabbitmq.NewConsume(url, &mockConsume{})
		consume.Start()
	}()
	time.Sleep(1 * time.Second)

	rc := gateway.New("UserCenter", rabbitmq.NewProduce("amqp://admin:admin@localhost:5672/"))
	router := xgin.New(rc.Handler)

	f := func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
	f2 := func(c *gin.Context) {
		gateway.SetClientId(trace.GetContextFrom(c), "aaaaaa")
		c.JSON(200, gin.H{
			"message": "注册用户",
		})
	}

	router.GET("/ping", rc.RegisterHandler("PING", "PING_TYPE", f))
	router.GET("/register", rc.RegisterHandler("注册用户", "USER_TYPE", f2))
	router.GET("/token", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "token",
		})
	})

	err := rc.RecordeApi(router.Routes())
	check(err)
	router.Run(":8085")
}

type mockConsume struct {
}

func (m mockConsume) ApiHandler(api *gateway.Api) error {
	fmt.Printf("api: %+v\n", api)
	return nil
}

func (m mockConsume) ApiCallHandler(entity *gateway.ApiCallEntity) error {
	fmt.Printf("entity: %+v\n", entity)
	return nil
}

func (m mockConsume) ContractCallHandler(entity *gateway.ContractCallEntity) error {
	fmt.Printf("entity: %+v\n", entity)
	return nil
}
