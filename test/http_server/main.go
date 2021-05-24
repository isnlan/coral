package main

import (
	"os"
	"time"

	"github.com/isnlan/coral/pkg/xgin"

	"github.com/isnlan/coral/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/isnlan/coral/pkg/trace"
)

const host = "172.20.107.44:32333"
const local = "127.0.0.1:6831"

var logger = logging.MustGetLogger("httpserver")

func main() {
	_, closer, err := trace.NewTracer("scpkg-gin-server", local)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	r := xgin.New()
	//r.Use(trace.TracerWrapper)

	logging.Init("ping", logging.NewWriteSyncerConfig(os.Stderr))

	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 200)
		ctx := trace.GetContextFrom(c)
		logger.With(trace.GetTraceFieldFrom(ctx)...).Info("get client request")
		spanA, ctxS := trace.StartSpanFromContext(ctx, "save object to db")
		time.Sleep(time.Second)

		spanA.Finish()

		spanB, _ := trace.StartSpanFromContext(ctx, "save cache")
		time.Sleep(time.Second)
		spanB.Finish()

		ctxS.Err()

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("0.0.0.0:8099") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
