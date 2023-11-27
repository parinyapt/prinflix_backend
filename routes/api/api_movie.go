package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parinyapt/prinflix_backend/handler"
	"github.com/parinyapt/prinflix_backend/middleware"
)

func InitMovieAPI(router *gin.RouterGroup) {
	r := router.Group("/movie").Use(middleware.GetHeaderAuthorizationToken, middleware.AuthWithAccessToken)
	{
		middlewareUser := middleware.NewMiddleware(middleware.MiddlewareReceiverArgument{})
		user := r.Use(middlewareUser.CheckAccount)
		user.GET("category", handler.GetMovieCategoryListHandler)
		user.GET("", handler.GetMovieListHandler)
		user.GET("/:movie_uuid", handler.GetMovieDetailHandler)
		user.GET("favorite", handler.GetFavoriteMovieListHandler)
		user.POST("/:movie_uuid/favorite", handler.AddFavoriteMovieHandler)
		user.DELETE("/:movie_uuid/favorite", handler.RemoveFavoriteMovieHandler)
		// user.GET("recommend", handler.GetRecommendMovieListHandler)
		// user.GET("continue_watching", handler.GetContinueWatchingMovieListHandler)
		user.POST("/:movie_uuid/watch", handler.RequestWatchMovieHandler)
		watchSessionProtect := user.Use(middleware.WatchSessionCheck)
		{
			watchSessionProtect.POST("/watch/pause", handler.RequestPauseMovieHandler)
			watchSessionProtect.POST("/watch/end", handler.RequestEndMovieHandler)
		}
	}
}