package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parinyapt/prinflix_backend/handler"
)

func configApiRoutes(router *gin.Engine) {
	// No Route 404 Notfound
	router.NoRoute(handler.NoRouteHandler)
	// Health Check
	router.GET("healthz", handler.HealthCheckHandler)
}