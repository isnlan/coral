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
	apis     map[string]*API
	producer Producer
}

func New(appName string, producer Producer) *Gateway {
	return &Gateway{
		appName:  appName,
		apis:     map[string]*API{},
		producer: producer,
	}
}

func (r *Gateway) RegisterHandler(apiName string, apiType string, f gin.HandlerFunc) gin.HandlerFunc {
	api := new(API)
	api.APIName = apiName
	api.APIType = apiType

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
		api.ID = utils.MakeMongoIdFromString(id)
		api.Scheme = "HTTP"
		api.Method = router.Method
		api.Path = router.Path
		api.AppName = r.appName
		api.DocURL = ""

		err := r.producer.APIUpload(api)
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
		err := r.producer.APICallRecord(&APICallEntity{
			APIID:    api.ID,
			Latency:  time.Now().Sub(start).Milliseconds(),
			HttpCode: c.Writer.Status(),
			ClientID: GetClientID(c),
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
		ClientID:  identity.ClientId,
		Address:   identity.Bduid,
		ChainID:   identity.NetworkId,
		ChannelID: channel,
		Contract:  contract,
	}

	return r.producer.ContractCallRecord(entity)
}

func makeFuncName(f gin.HandlerFunc) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
