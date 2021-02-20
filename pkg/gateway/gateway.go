package gateway

import (
	"reflect"
	"runtime"
	"time"

	"github.com/isnlan/coral/pkg/protos"

	"github.com/gin-gonic/gin"

	"github.com/isnlan/coral/pkg/logging"
)

var logger = logging.MustGetLogger("gateway")

type Gateway struct {
	appName  string
	fns      map[string]string
	apis     map[string]*Api
	producer Producer
}

func New(appName string, producer Producer) *Gateway {
	return &Gateway{
		appName:  appName,
		fns:      map[string]string{},
		apis:     map[string]*Api{},
		producer: producer,
	}
}

func (r *Gateway) RegisterHandler(apiName string, f gin.HandlerFunc) gin.HandlerFunc {
	r.fns[makeFuncName(f)] = apiName
	return f
}

func (r *Gateway) RecordeApi(rs gin.RoutesInfo) error {
	for _, router := range rs {
		apiName := r.fns[router.Handler]
		if apiName == "" {
			continue
		}

		api := NewApi(r.appName, "http", router.Method, router.Path, apiName, "")

		r.apis[router.Handler] = api
		err := r.producer.ApiUpload(api)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Gateway) Handler(c *gin.Context) {
	// start timer
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	// Process request
	c.Next()

	if raw != "" {
		path = path + "?" + raw
	}

	if api, find := r.apis[c.HandlerName()]; find {
		err := r.producer.ApiCallRecord(&ApiCallEntity{
			ApiId:    api.Id,
			Latency:  time.Now().Sub(start).Milliseconds(),
			HttpCode: c.Writer.Status(),
			ClientId: GetClientId(c),
		})
		if err != nil {
			logger.Errorf("api call record error: %v", err)
		}
	}
}

func (r *Gateway) RecordeContractCall(identity *protos.DigitalIdentity, channel, contract string) error {
	if identity == nil {
		return nil
	}

	entity := &ContractCallEntity{
		ClientId:  identity.ClientId,
		Address:   identity.Bduid,
		ChainId:   identity.NetworkId,
		ChannelId: channel,
		Contract:  contract,
	}

	return r.producer.ContractCallRecord(entity)
}

func makeFuncName(f gin.HandlerFunc) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
