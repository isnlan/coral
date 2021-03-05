package xgrpc

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"runtime"
	"time"

	"github.com/isnlan/coral/pkg/utils"

	"github.com/isnlan/coral/pkg/logging"

	"google.golang.org/grpc"

	"github.com/isnlan/coral/pkg/errors"
	"google.golang.org/grpc/peer"
)

var (
	logger    = logging.MustGetLogger("grpc")
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
	reset     = string([]byte{27, 91, 48, 109})
)

func LoggerUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		clientIP, service, method := unaryRequestInfo(ctx, info)

		var statusColor, methodColor, resetColor string

		resp, err = handler(ctx, req)
		var statusCode = 200
		if err != nil {
			statusCode = 500
		}

		logger.Infof("[GRPC]%s %3d %s| %v | %s |%s %-2s %s %s %#v\n",
			statusColor, statusCode, resetColor,
			time.Now().Sub(start),
			clientIP,
			methodColor, service, resetColor,
			method,
			req,
		)

		return resp, err
	}
}

func RecoveryUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		clientIP, service, method := unaryRequestInfo(ctx, info)
		reqinfo := fmt.Sprintf("[%s]%s %s", clientIP, service, method)

		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(ctx, r, reqinfo)
			}
		}()
		return handler(ctx, req)
	}
}

func RecoveryStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(stream.Context(), r, "")
			}
		}()

		return handler(srv, stream)
	}
}

func unaryRequestInfo(ctx context.Context, info *grpc.UnaryServerInfo) (string, string, string) {
	var clientIP string
	if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}

	service := utils.MakeTypeName(info.Server)
	method := info.FullMethod
	return clientIP, service, method
}

func recoverFrom(ctx context.Context, r interface{}, reqinfo string) error {
	switch v := r.(type) {
	case runtime.Error:
		stack := stack(3)
		logger.Errorf("runtime error:\n%s\n%s\n%s%s", reqinfo, r, stack, reset)
		return fmt.Errorf("runtime error: %v", v)
	case errors.CodeError:
		logger.Errorf("request %s error, code %d, description %s", reqinfo, v.Code(), v.Error())
		return v
	case error:
		logger.Errorf("request %s error, code %d, description %s", reqinfo, errors.InternalErrorCode, v.Error())
		return v
	default:
		logger.Errorf("[Recovery] panic recovered: %s\n%s\n", reqinfo, r)
		return fmt.Errorf("unknown type error %v", v)
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
