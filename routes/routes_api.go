package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parinyapt/prinflix_backend/handler"
	APIroutes "github.com/parinyapt/prinflix_backend/routes/api"
)

func configApiRoutes(router *gin.Engine) {
	// No Route 404 Notfound
	router.NoRoute(handler.NoRouteHandler)
	// Health Check
	router.GET("healthz", handler.HealthCheckHandler)

	v1 := router.Group("/v1")
	{
		APIroutes.InitAuthAPI(v1)
		APIroutes.InitAccountAPI(v1)

	}
}