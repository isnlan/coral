package main

import (
	"github.com/gin-gonic/gin"
	"github.com/isnlan/coral/pkg/discovery"
	"github.com/isnlan/coral/pkg/discovery/consul"
	"github.com/isnlan/coral/pkg/errors"
	"github.com/isnlan/coral/pkg/xgin"
	"github.com/isnlan/coral/prometheus"
)

func main() {
	ds, err := consul.New("127.0.0.1:8500")
	errors.Check(err)

	ip, err := discovery.GetLocalIP()
	errors.Check(err)

	prometheus.StartAgent(ip, 9001)
	prometheus.RegisterAgent(ds, "myapp", ip, 9001)

	r := xgin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run("0.0.0.0:8099") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
