package prometheus

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/isnlan/coral/pkg/discovery"
	"github.com/isnlan/coral/pkg/logging"
	"go.uber.org/atomic"
)

var (
	enabled atomic.Bool
	logger  = logging.MustGetLogger("prometheus")
)

func Enabled() bool {
	return enabled.Load()
}

func StartAgent(port int) {
	enabled.Store(true)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc(discovery.HTTPHealthCheckRouter, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	addr := fmt.Sprintf(":%d", port)
	logger.Infof("starting prometheus agent at %s", addr)

	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			logger.Fatal(err)
		}
	}()
}

func RegisterAgent(sd discovery.ServiceDiscover, svc, host string, port int) {
	if !Enabled() {
		return
	}

	if sd == nil {
		return
	}

	_, err := sd.HTTPServiceRegister("metrics_exporter", host, port, svc)
	if err != nil {
		logger.Fatal(err)
	}
}
