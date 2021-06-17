package main

import (
	"time"

	gateway2 "github.com/isnlan/coral/pkg/blink/gateway"
	"github.com/isnlan/coral/pkg/blink/gateway/rabbitmq"

	"github.com/isnlan/coral/pkg/xgin"

	"github.com/gin-gonic/gin"
	"github.com/isnlan/coral/pkg/trace"
)

const host = "172.20.107.44:32333"
const local = "127.0.0.1:6831"

func main() {
	_, closer, err := trace.NewTracer("scpkg-gin-server", local)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	r := xgin.New()

	url := "amqp://admin:admin@localhost:5672/"
	produce := rabbitmq.NewProduce(url)
	apis := []string{"df05961e491bb6a77edeb7fc", "f15883f21409ed3f0eb34cff"}

	r.GET("/ping", func(c *gin.Context) {
		i := time.Now().Unix()
		err := produce.APICallRecord(&gateway2.APICallEntity{
			APIID:    apis[i%2],
			Latency:  int64(i % 1000),
			HttpCode: 200,
			ClientID: "473d78a37c640099",
		})
		if err != nil {
			c.JSON(404, gin.H{
				"message": "error",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		}
	})
	r.Run("0.0.0.0:8099") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
