package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/parinyapt/prinflix_backend/handler"
	"github.com/parinyapt/prinflix_backend/middleware"
)

func InitAccountAPI(router *gin.RouterGroup) {
	r := router.Group("/account").Use(middleware.GetHeaderAuthorizationToken, middleware.AuthWithAccessToken)
	{
		middlewareUser := middleware.NewMiddleware(middleware.MiddlewareReceiverArgument{})
		user := r.Use(middlewareUser.CheckAccount)
		user.PUT("/profile", handler.UpdateProfileHandler)
		user.PUT("/password", handler.UpdatePasswordHandler)
	}
}