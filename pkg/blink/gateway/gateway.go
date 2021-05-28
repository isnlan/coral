package gateway

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/isnlan/coral/pkg/utils"

	"github.com/isnlan/coral/pkg/protos"

	"github.com/gin-gonic/gin"

	"github.com/isnlan/coral/pkg/logging"
)

var logger = logging.MustGetLogger("gateway")

type Gateway struct {
	appName  string
	apis     map[string]*Api
	producer Producer
}

func New(appName string, producer Producer) *Gateway {
	return &Gateway{
		appName:  appName,
		apis:     map[string]*Api{},
		producer: producer,
	}
}

func (r *Gateway) RegisterHandler(apiName string, apiType string, f gin.HandlerFunc) gin.HandlerFunc {
	api := new(Api)
	api.ApiName = apiName
	api.ApiType = apiType

	r.apis[makeFuncName(f)] = api
	return f
}

func (r *Gateway) RecordeApi(rs gin.RoutesInfo) error {
	for _, router := range rs {
		api, find := r.apis[router.Handler]
		if !find {
			continue
		}

		id := fmt.Sprintf("%s:[%s.%s] %s", r.appName, "HTTP", router.Method, router.Path)
		api.Id = utils.MakeMongoIdFromString(id)
		api.Scheme = "HTTP"
		api.Method = router.Method
		api.Path = router.Path
		api.AppName = r.appName
		api.DocUrl = ""

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
