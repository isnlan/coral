package gateway

import (
	"context"
	"fmt"

	"github.com/isnlan/coral/pkg/utils"

	"github.com/gin-gonic/gin"

	"github.com/isnlan/coral/pkg/trace"
)

const _ClientIdContextKey = "ClientIdContext"

type Api struct {
	Id      string `json:"id"`       // api id MongoId
	Scheme  string `json:"scheme"`   // http ws grpc
	Method  string `json:"method"`   // http method
	Path    string `json:"path"`     // http route path
	AppName string `json:"app_name"` // application name
	ApiName string `json:"api_name"` // api 中文名称
	ApiType string `json:"api_type"` // api 接口类型
	DocUrl  string `json:"doc_url"`  // 文档地址
}

func NewApi(appName, scheme, method, path, apiName, apiType, docUrl string) *Api {
	api := fmt.Sprintf("%s:[%s.%s] %s", appName, scheme, method, path)
	return &Api{
		Id:      utils.MakeMongoIdFromString(api),
		Scheme:  scheme,
		Method:  method,
		Path:    path,
		AppName: appName,
		ApiName: apiName,
		ApiType: apiType,
		DocUrl:  docUrl,
	}
}

type ApiCallEntity struct {
	ApiId    string `json:"api_id"`    // api id
	Latency  int64  `json:"latency"`   // 耗时
	HttpCode int    `json:"http_code"` // http状态码
	ClientId string `json:"client_id"` // client id
}

type ContractCallEntity struct {
	ClientId  string `json:"client_id" validate:"required"`  // 客户端ID
	Address   string `json:"address" validate:"required"`    // 数字身份
	ChainId   string `json:"chain_id" validate:"required"`   // 网络ID
	ChannelId string `json:"channel_id" validate:"required"` // 链ID
	Contract  string `json:"contract" validate:"omitempty"`  // 合约名称
}

func SetClientId(ctx context.Context, clientId string) {
	c := trace.GetGinContext(ctx)
	if c != nil {
		c.Set(_ClientIdContextKey, clientId)
	}
}

func GetClientId(c *gin.Context) string {
	value := c.Value(_ClientIdContextKey)
	if value == nil {
		return ""
	}
	return value.(string)
}
