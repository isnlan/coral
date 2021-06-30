package xgin

import (
	"strconv"
	"time"

	prometheus3 "github.com/isnlan/coral/pkg/prometheus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/gin-gonic/gin"
)

const serverNamespace = "http_server"

var (
	metricServerRequestDurations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   serverNamespace,
		Subsystem:   "requests",
		Name:        "duration_ms",
		Help:        "http server requests duration(ms).",
		ConstLabels: map[string]string{},
		Buckets:     []float64{5, 10, 25, 50, 100, 250, 500, 1000, 5000, 10000, 30000, 60000},
	}, []string{"path"})

	metricServerRequestCodeTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace:   serverNamespace,
		Subsystem:   "requests",
		Name:        "code_total",
		Help:        "http server requests error count.",
		ConstLabels: map[string]string{},
	}, []string{"path", "code"})
)

func RecordMetrics() gin.HandlerFunc {
	if !prometheus3.Enabled() {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		startTime := time.Now()
		defer func() {
			metricServerRequestDurations.WithLabelValues(c.Request.URL.Path).Observe(float64(time.Since(startTime) / time.Millisecond))
			metricServerRequestCodeTotal.WithLabelValues(c.Request.URL.Path, strconv.Itoa(c.Writer.Status())).Inc()
		}()

		c.Next()
	}
}
