package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/parinyapt/prinflix_backend/handler"
	// "github.com/parinyapt/prinflix_backend/middleware"
)

func InitAuthAPI(router *gin.RouterGroup) {
	r := router.Group("/auth")
	{
		r.POST("/login", handler.LoginHandler)
		r.POST("/register", handler.RegisterHandler)

		// authWithAccessToken := r.Group("/token").Use(middleware.GetHeaderAuthorizationToken, middleware.AuthWithAccessToken)
		// {
		// 	authWithAccessToken.GET("/verify", handler.VerifyTokenHandler)
		// 	authWithAccessToken.POST("/refresh", handler.RefreshTokenHandler)
		// 	authWithAccessToken.POST("/revoke", handler.RevokeTokenHandler)
		// }
	}
}