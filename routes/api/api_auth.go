package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/parinyapt/prinflix_backend/handler"
	"github.com/parinyapt/prinflix_backend/middleware"
)

func InitAuthAPI(router *gin.RouterGroup) {
	r := router.Group("/auth")
	{
		r.POST("/login", handler.LoginHandler)
		r.POST("/register", handler.RegisterHandler)
		r.POST("/forgot_password", handler.RequestForgotPasswordHandler)
		r.GET("/forgot_password/:session_id", handler.CheckForgotPasswordSessionHandler)
		r.POST("/reset_password", handler.ResetPasswordHandler)
		r.GET("/email_verify/:code", handler.EmailVerifyHandler)

		google := r.Group("/google")
		{
			public := google.Group("")
			{
				public.GET("/callback", handler.GoogleCallbackHandler)
			}
			external := google.Group("").Use(middleware.GetHeaderAuthorizationToken, middleware.AuthWithAccessToken)
			{
				middlewareUser := middleware.NewMiddleware(middleware.MiddlewareReceiverArgument{})
				user := external.Use(middlewareUser.CheckAccount)
				user.POST("/connect", handler.RequestConnectGoogleOAuthHandler)
				user.POST("/disconnect", handler.RequestDisconnectGoogleOAuthHandler)
			}
			internal := google.Group("").Use(middleware.GetHeaderAPIKey, middleware.AuthWithAPIKey)
			{
				internal.POST("/login", handler.InternalGoogleLoginHandler)
			}
		}

		line := r.Group("/line")
		{
			public := line.Group("")
			{
				public.GET("/callback", handler.LineCallbackHandler)
			}
			external := line.Group("").Use(middleware.GetHeaderAuthorizationToken, middleware.AuthWithAccessToken)
			{
				middlewareUser := middleware.NewMiddleware(middleware.MiddlewareReceiverArgument{})
				user := external.Use(middlewareUser.CheckAccount)
				user.POST("/connect", handler.RequestConnectLineOAuthHandler)
				user.POST("/disconnect", handler.RequestDisconnectLineOAuthHandler)
			}
			internal := line.Group("").Use(middleware.GetHeaderAPIKey, middleware.AuthWithAPIKey)
			{
				internal.POST("/login", handler.InternalLineLoginHandler)
			}
		}

		r.POST("/email_verify", middleware.GetHeaderAuthorizationToken, middleware.AuthWithAccessToken, handler.RequestEmailVerifyHandler)

		authWithAccessToken := r.Group("/token").Use(middleware.GetHeaderAuthorizationToken, middleware.AuthWithAccessToken)
		{
			authWithAccessToken.GET("/verify", handler.VerifyTokenHandler)
			authWithAccessToken.POST("/refresh", handler.RefreshTokenHandler)
			authWithAccessToken.POST("/revoke", handler.RevokeTokenHandler)
		}
	}
}
