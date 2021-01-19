package trace

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func OpenTracingClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string, req, resp interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) (err error) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "RPC Client "+method)
		defer span.Finish()

		// Save current span context.
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}
		if err = opentracing.GlobalTracer().Inject(
			span.Context(), opentracing.HTTPHeaders, metadataTextMap(md),
		); err != nil {
			logger.Errorf("Failed to inject trace span: %v, error: %v", ctx, err)
		}
		return invoker(metadata.NewOutgoingContext(ctx, md), method, req, resp, cc, opts...)
	}
}

func OpenTracingServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Extract parent trace span.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}
		parentSpanContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders, metadataTextMap(md),
		)

		switch err {
		case nil:
		case opentracing.ErrSpanContextNotFound:
			logger.Debugf("Parent span not found, will start new one. span: %v", ctx)
		default:
			logger.Errorf("Failed to extract trace span: %v, error: %v", ctx, err)
		}

		// Start new trace span.
		span := opentracing.StartSpan(
			"RPC Server "+info.FullMethod,
			ext.RPCServerOption(parentSpanContext),
		)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)

		return handler(ctx, req)
	}
}

const (
	binHeaderSuffix = "_bin"
)

// metadataTextMap extends a metadata.MD to be an opentracing textmap
type metadataTextMap metadata.MD

// Set is a opentracing.TextMapReader interface that extracts values.
func (m metadataTextMap) Set(key, val string) {
	// gRPC allows for complex binary values to be written.
	encodedKey, encodedVal := encodeKeyValue(key, val)
	// The metadata object is a multimap, and previous values may exist, but for opentracing headers, we do not append
	// we just override.
	m[encodedKey] = []string{encodedVal}
}

// ForeachKey is a opentracing.TextMapReader interface that extracts values.
func (m metadataTextMap) ForeachKey(callback func(key, val string) error) error {
	for k, vv := range m {
		for _, v := range vv {
			if decodedKey, decodedVal, err := metadata.DecodeKeyValue(k, v); err == nil {
				if err = callback(decodedKey, decodedVal); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed decoding opentracing from gRPC metadata: %v", err)
			}
		}
	}
	return nil
}

// encodeKeyValue encodes key and value qualified for transmission via gRPC.
// note: copy pasted from private values of grpc.metadata
func encodeKeyValue(k, v string) (string, string) {
	k = strings.ToLower(k)
	if strings.HasSuffix(k, binHeaderSuffix) {
		v = base64.StdEncoding.EncodeToString([]byte(v))
	}
	return k, v
}
