package main

import (
	"fmt"
	"time"

	gateway2 "github.com/isnlan/coral/pkg/blink/gateway"
	rabbitmq2 "github.com/isnlan/coral/pkg/blink/gateway/rabbitmq"

	"github.com/gin-gonic/gin"
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
		consume := rabbitmq2.NewConsume(url, &mockConsume{})
		consume.Start()
	}()
	time.Sleep(1 * time.Second)

	rc := gateway2.New("UserCenter", rabbitmq2.NewProduce("amqp://admin:admin@localhost:5672/"))
	router := xgin.New(rc.Handler)

	f := func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
	f2 := func(c *gin.Context) {
		gateway2.SetClientId(trace.GetContextFrom(c), "aaaaaa")
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

func (m mockConsume) ApiHandler(api *gateway2.Api) error {
	fmt.Printf("api: %+v\n", api)
	return nil
}

func (m mockConsume) ApiCallHandler(entity *gateway2.ApiCallEntity) error {
	fmt.Printf("entity: %+v\n", entity)
	return nil
}

func (m mockConsume) ContractCallHandler(entity *gateway2.ContractCallEntity) error {
	fmt.Printf("entity: %+v\n", entity)
	return nil
}
