package router

import (
	"github.com/gin-gonic/gin"
	"github.com/snlansky/coral/pkg/middleware"
	"github.com/snlansky/coral/pkg/trace"
)

func New() *gin.Engine {
	router := gin.New()
	router.Use(middleware.LoggerWriter(), middleware.RecoveryWriter(), middleware.CorsMiddleware(), trace.TracerWrapper)
	router.NoRoute(middleware.HandleNotFound)
	router.NoMethod(middleware.HandleNotFound)
	return router
}
