package trace

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// sf sampling frequency
var sf = 100

func init() {
	rand.Seed(time.Now().Unix())
}

// SetSamplingFrequency 设置采样频率
// 0 <= n <= 100
func SetSamplingFrequency(n int) {
	sf = n
}

// TracerWrapper tracer 中间件
func TracerWrapper(c *gin.Context) {
	md := make(map[string]string)
	nsf := sf
	var opts []opentracing.StartSpanOption
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	if err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
		nsf = 100
	}
	sp := opentracing.GlobalTracer().StartSpan(c.Request.URL.EscapedPath(), opts...)
	defer sp.Finish()

	if err := sp.Tracer().Inject(sp.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(md)); err != nil {
		logger.Errorf("trace inject error: %v", err)
	}

	ctx := opentracing.ContextWithSpan(c, sp)
	ctx = metadata.NewContext(ctx, md)

	setContext(c, ctx)
	traceId := getTraceIDFromSpan(sp)
	if traceId != "" {
		logger.Infof("trace_id: %s", traceId)
	}

	c.Next()

	statusCode := c.Writer.Status()
	ext.HTTPStatusCode.Set(sp, uint16(statusCode))
	ext.HTTPMethod.Set(sp, c.Request.Method)
	ext.HTTPUrl.Set(sp, c.Request.URL.EscapedPath())
	if statusCode >= http.StatusInternalServerError {
		ext.Error.Set(sp, true)
	} else if rand.Intn(100) > nsf {
		ext.SamplingPriority.Set(sp, 0)
	}
}

func setContext(c *gin.Context, ctx context.Context) {
	_, ok := ctx.Deadline()
	if !ok {
		ctx, _ = context.WithTimeout(ctx, time.Minute)
	}

	ctx = context.WithValue(ctx, _UrlKey, c.Request.URL.EscapedPath())
	ctx = context.WithValue(ctx, _GinContextKey, c)
	c.Set(_ContextTracerKey, ctx)
}

func GetGinContext(ctx context.Context) *gin.Context {
	value := ctx.Value(_GinContextKey)
	if value == nil {
		return nil
	}
	return value.(*gin.Context)
}
