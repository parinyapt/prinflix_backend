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

	storage := router.Group("/storage")
	{
		v1 := storage.Group("/v1")
		{
			APIroutes.InitStorageAPI(v1)
		}
	}

	v1 := router.Group("/v1")
	{
		APIroutes.InitAuthAPI(v1)
		APIroutes.InitAccountAPI(v1)
		APIroutes.InitMovieAPI(v1)
		// APIroutes.InitPaymentAPI(v1)

	}

	v2 := router.Group("/v2")
	{
		APIroutes.InitAuthAPIv2(v2)
	}
}