package trace

import (
	"context"

	"go.uber.org/zap"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	"github.com/gin-gonic/gin"
	"github.com/isnlan/coral/pkg/logging"
)

const (
	_UrlKey           = "_Url_"
	_ContextTracerKey = "_TracerContext_"
	_GinContextKey    = "_GinContext_"
)

var logger = logging.MustGetLogger("trace")

// ContextWithSpan 返回context
func GetContextFrom(c *gin.Context) (ctx context.Context) {
	v, exist := c.Get(_ContextTracerKey)
	if !exist {
		ctx = context.Background()
		return
	}

	ctx, ok := v.(context.Context)
	if !ok {
		panic("GetContext Error")
	}
	return
}

var StartSpanFromContext = opentracing.StartSpanFromContext

func StartSpan(ctx context.Context, operationName string, f func() error, opts ...opentracing.StartSpanOption) (context.Context, error) {
	sp, _ctx := opentracing.StartSpanFromContext(ctx, operationName, opts...)
	defer sp.Finish()
	return _ctx, f()
}

func GetTraceIDFromContext(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		return getTraceIDFromSpan(span)
	}

	return ""
}

func getTraceIDFromSpan(span opentracing.Span) string {
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		return sc.TraceID().String()
	}
	return ""
}

func GetSpanIDFromContext(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		return getSpanIDFromSpan(span)
	}

	return ""
}

func getSpanIDFromSpan(span opentracing.Span) string {
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		return sc.SpanID().String()
	}
	return ""
}

func GetParentIDFromSpan(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		return getParentIDFromSpan(span)
	}

	return ""
}

func getParentIDFromSpan(span opentracing.Span) string {
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		return sc.ParentID().String()
	}
	return ""
}

func GetUrlFromContext(ctx context.Context) string {
	return getStringFromContext(ctx, _UrlKey)
}

func getStringFromContext(ctx context.Context, key string) string {
	if ctx != nil {
		if v := ctx.Value(key); v != nil {
			if url, ok := v.(string); ok {
				return url
			}
		}
	}
	return ""
}

func GetTraceFieldFrom(ctx context.Context) []interface{} {
	var fields []interface{}
	fields = append(fields, zap.String("service", logging.ServiceName))
	fields = append(fields, zap.String("trace", GetTraceIDFromContext(ctx)))
	fields = append(fields, zap.String("span", GetSpanIDFromContext(ctx)))
	fields = append(fields, zap.String("parent", GetParentIDFromSpan(ctx)))
	return fields
}
