package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"runtime"
	"time"

	"github.com/snlansky/coral/pkg/errors"
	"github.com/snlansky/coral/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/snlansky/coral/pkg/logging"
)

var logger = logging.MustGetLogger("http")

// LoggerWithWriter instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWriter() gin.HandlerFunc {
	var skip map[string]struct{}
	return func(c *gin.Context) {
		// start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			var statusColor, methodColor, resetColor string
			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			logger.Infof("[HTTP]%s %3d %s| %v | %s |%s %-2s %s %s\n%s",
				statusColor, statusCode, resetColor,
				latency,
				clientIP,
				methodColor, method, resetColor,
				path,
				comment,
			)
		}
	}
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
	reset     = string([]byte{27, 91, 48, 109})
)

// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func RecoveryWriter() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				switch resp := r.(type) {
				case errors.CodeError:
					c.AbortWithStatusJSON(http.StatusOK, &response.JsonResponse{
						ErrorCode:   resp.Code(),
						Description: resp.Error(),
						Data:        nil,
					})
				case error:
					c.AbortWithStatusJSON(http.StatusOK, &response.JsonResponse{
						ErrorCode:   errors.InternalErrorCode,
						Description: resp.Error(),
						Data:        nil,
					})
				default:
					if logger != nil {
						stack := stack(3)
						httprequest, _ := httputil.DumpRequest(c.Request, false)
						logger.Errorf("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), r, stack, reset)
					}
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
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

func HandleNotFound(c *gin.Context) {
	c.JSON(404, gin.H{
		"message": "404 page not found",
		"request": c.Request.Method + " " + c.Request.URL.String(),
	})
}
