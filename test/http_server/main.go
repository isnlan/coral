package main

import (
	"os"
	"time"

	"github.com/isnlan/coral/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/isnlan/coral/pkg/trace"
)

const host = "172.20.107.44:32333"
const local = "127.0.0.1:6831"

var logger = logging.MustGetLogger("mysvr")

func main() {
	_, closer, err := trace.NewTracer("scpkg-gin-server", local)
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	logging.Init(logging.NewFileConfig("http", os.Stderr))

	r := gin.Default()
	r.Use(trace.TracerWrapper)

	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 200)

		ctx := trace.GetContextFrom(c)
		spanA, ctxS := trace.StartSpanFromContext(ctx, "save object to db")
		time.Sleep(time.Second)

		spanA.Finish()

		logger.With(trace.GetFields(ctx)).Info("--------")

		spanB, _ := trace.StartSpanFromContext(ctx, "save cache")
		time.Sleep(time.Second)
		spanB.Finish()

		ctxS.Err()

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("0.0.0.0:8090") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
