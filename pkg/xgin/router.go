package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/isnlan/coral/pkg/logging"
	"github.com/isnlan/coral/pkg/trace"
)

var logger = logging.MustGetLogger("xgin")

func New(middleware ...gin.HandlerFunc) *gin.Engine {
	middleware = append(middleware, LoggerWriter(), RecoveryWriter(), CorsMiddleware(), trace.TracerWrapper, RecordMetrics())

	router := gin.New()
	router.Use(middleware...)
	router.NoRoute(HandleNotFound)
	router.NoMethod(HandleNotFound)
	return router
}
