package main

import (
	"github.com/gin-gonic/gin"
	"github.com/snlansky/coral/pkg/gateway"
	"github.com/snlansky/coral/pkg/gateway/rabbitmq"
	"github.com/snlansky/coral/pkg/trace"
	"github.com/snlansky/coral/pkg/xgin"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
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

	router.GET("/ping", rc.RegisterHandler("PING", f))
	router.GET("/register", rc.RegisterHandler("注册用户", f2))
	router.GET("/token", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "token",
		})
	})

	err := rc.RecordeApi(router.Routes())
	check(err)
	router.Run(":8085")
}
